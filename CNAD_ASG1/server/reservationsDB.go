package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// retrieve all reservations by user ID
func RetrieveReservations(db *sql.DB, userID int) ([]Reservation, error) {
	updateQuery := `
        UPDATE reservations
        SET status = 'Completed'
        WHERE end_time < NOW() AND status != 'Completed' AND status != 'Cancelled'
    `
	// Execute the update query
	_, err := db.Exec(updateQuery)
	if err != nil {
		// Handle query execution error
		log.Println("Error executing update query:", err)
		return nil, err
	} else {
		log.Println("Ran update statement")
	}

	rows, err := db.Query(`SELECT r.id, r.user_id, r.vehicle_id, v.model, r.start_time, r.end_time, r.status, r.estimated_cost 
	FROM Reservations r INNER JOIN Vehicles v ON r.vehicle_id = v.id WHERE user_id = ? ORDER BY r.id DESC`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var r Reservation
		err := rows.Scan(&r.ID, &r.UserID, &r.VehicleID, &r.Model, &r.StartTime, &r.EndTime, &r.Status, &r.EstimatedCost)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}

	return reservations, nil
}

// retrieve reservation by reservation ID
func RetrieveReservationByID(db *sql.DB, reservationID int) (Reservation, error) {
	var reservation Reservation

	query := `
        SELECT id, user_id, vehicle_id, start_time, end_time, status, estimated_cost
        FROM reservations
        WHERE id = ?
    `
	err := db.QueryRow(query, reservationID).Scan(
		&reservation.ID,
		&reservation.UserID,
		&reservation.VehicleID,
		&reservation.StartTime,
		&reservation.EndTime,
		&reservation.Status,
		&reservation.EstimatedCost,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return reservation, sql.ErrNoRows
		}
		return reservation, err
	}

	return reservation, nil
}

// insert reservation record into database
func CreateReservation(db *sql.DB, userID, vehicleID int, startTime time.Time, endTime time.Time, estimatedCost float64) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// check if the vehicle is booked during selected time with status 'Confirmed'
	query := `
        SELECT COUNT(*) 
        FROM reservations 
        WHERE vehicle_id = ? 
          AND (start_time BETWEEN ? AND ? 
               OR end_time BETWEEN ? AND ? 
               OR ? BETWEEN start_time AND end_time 
               OR ? BETWEEN start_time AND end_time)
          AND status IN ('Confirmed')`
	var count int
	err = tx.QueryRow(query, vehicleID, startTime, endTime, startTime, endTime, startTime, endTime).Scan(&count)
	if err != nil {
		return 0, err
	}

	// vehicle booked
	if count > 0 {
		return 0, errors.New("vehicle is already reserved during the selected time frame")
	}

	// insert reservation into database
	insertQuery := `
        INSERT INTO reservations (user_id, vehicle_id, start_time, end_time, status, estimated_cost) 
        VALUES (?, ?, ?, ?, 'Confirmed', ?)`
	res, err := tx.Exec(insertQuery, userID, vehicleID, startTime, endTime, estimatedCost)
	if err != nil {
		return 0, err
	}

	// Get the newly created reservation ID
	reservationID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(reservationID), nil
}

// updates the start_time and end_time of a reservation
func ModifyReservation(db *sql.DB, reservationID int, startTime, endTime time.Time, estimatedCost float64) error {
	query := `
        UPDATE reservations
        SET start_time = ?, end_time = ?, estimated_cost = ?
        WHERE id = ?`
	_, err := db.Exec(query, startTime, endTime, estimatedCost, reservationID)
	return err
}

// cancels a reservation - changes the reservation status to "Cancelled"
func CancelReservation(db *sql.DB, reservationID int) error {
	// Fetch vehicle ID linked to the reservation
	var vehicleID int
	query := `SELECT vehicle_id FROM reservations WHERE id = ?`
	err := db.QueryRow(query, reservationID).Scan(&vehicleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("reservation not found")
		}
		return err
	}

	// Update reservation status
	updateReservationQuery := `
        UPDATE reservations 
        SET status = 'Cancelled' 
        WHERE id = ?`
	_, err = db.Exec(updateReservationQuery, reservationID)
	if err != nil {
		return err
	}

	return nil
}

// retrieve esimtated cost based on user's membership tier
func RetrieveEstimatedCost(db *sql.DB, userID int, days int) (float64, error) {
	query := `
	SELECT m.hourly_rate 
	FROM Membership m INNER JOIN Users u ON m.id = u.membership_id
	WHERE u.id = ?
	`
	var hourlyRate float64
	// Execute the query and scan the result into hourlyRate
	err := db.QueryRow(query, userID).Scan(&hourlyRate)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve hourly rate: %v", err)
	}

	estimatedCost := hourlyRate * float64(days)
	return estimatedCost, nil
}
