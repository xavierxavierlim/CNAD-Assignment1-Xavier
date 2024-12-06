package model

import (
	"database/sql"
	"errors"
	"time"
)

type Vehicle struct {
	ID                int    `json:"id"`
	Model             string `json:"model"`
	LicensePlate      string `json:"license_plate"`
	Location          string `json:"location"`
	ChargeLevel       int    `json:"charge_level"`
	CleanlinessStatus string `json:"cleanliness_status"`
}

type Reservation struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	VehicleID     int       `json:"vehicle_id"`
	Model         string    `json:"vehicle_name"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Status        string    `json:"status"`
	EstimatedCost float64   `json:"estimated_cost"`
}

// retrieve all available vehicles
func GetAvailableVehicles(db *sql.DB, startTime, endTime string) ([]Vehicle, error) {
	var vehicles []Vehicle
	query := `
    SELECT id, model, license_plate, location, charge_level, cleanliness_status
    FROM vehicles
    WHERE id NOT IN (
        SELECT vehicle_id
        FROM reservations
        WHERE (? < end_time AND ? > start_time)
		AND status NOT IN ("Cancelled")
    )`
	rows, err := db.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// looped through retrieved vehicles
	for rows.Next() {
		var v Vehicle
		err := rows.Scan(&v.ID, &v.Model, &v.LicensePlate, &v.Location, &v.ChargeLevel, &v.CleanlinessStatus)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v) // add vehicle into vehicles list
	}

	return vehicles, nil
}

/*
// retrieves vehicle information
func GetVehicleStatus(db *sql.DB, vehicleID int) (Vehicle, error) {
	var vehicle Vehicle

	query := `
        SELECT model, license_plate, location, charge_level, cleanliness_status
        FROM vehicles WHERE id = ?`
	err := db.QueryRow(query, vehicleID).
		Scan(&vehicle.Model, &vehicle.LicensePlate, &vehicle.Location, &vehicle.ChargeLevel, &vehicle.CleanlinessStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return vehicle, errors.New("vehicle not found")
		}
		return vehicle, err
	}

	return vehicle, nil
}
*/
/*
// retrieves all vehicles that are currently available
func GetVehicles(db *sql.DB) ([]Vehicle, error) {
	var vehicles []Vehicle

	query := `
	SELECT id, model, license_plate, location, charge_level, cleanliness_status`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v Vehicle
		err := rows.Scan(&v.ID, &v.Model, &v.LicensePlate, &v.Location, &v.ChargeLevel, &v.CleanlinessStatus)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return vehicles, nil
}
*/

/*
// retrieves all available vehicles for a specific date
func GetVehiclesByDate(db *sql.DB, date string) ([]Vehicle, error) {
	var vehicles []Vehicle

	query := `
        SELECT id, model, license_plate, location, charge_level, cleanliness_status
        FROM vehicles
          AND id NOT IN (
              SELECT vehicle_id
              FROM reservations
              WHERE DATE(start_time) <= ? AND DATE(end_time) >= ?
          )`
	rows, err := db.Query(query, date, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v Vehicle
		err := rows.Scan(&v.ID, &v.Model, &v.LicensePlate, &v.Location, &v.ChargeLevel, &v.CleanlinessStatus)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}
*/

// retrieve vehicle by vehicleID
func GetVehicleByID(db *sql.DB, vehicleID int) (Vehicle, error) {
	var vehicle Vehicle

	query := `
        SELECT id, model, license_plate, location, charge_level, cleanliness_status
        FROM vehicles
        WHERE id = ?`

	err := db.QueryRow(query, vehicleID).Scan(
		&vehicle.ID,
		&vehicle.Model,
		&vehicle.LicensePlate,
		&vehicle.Location,
		&vehicle.ChargeLevel,
		&vehicle.CleanlinessStatus,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return vehicle, errors.New("vehicle not found")
		}
		return vehicle, err
	}

	return vehicle, nil
}
