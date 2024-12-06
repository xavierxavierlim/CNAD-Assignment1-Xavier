package model

import (
	"database/sql"
	"errors"
	"log"
)

type Billing struct {
	ID            int     `json:"id"`
	ReservationID int     `json:"reservation_id"`
	UserID        int     `json:"user_id"`
	Amount        float64 `json:"amount"`
	PaymentStatus string  `json:"payment_status"`
}

// retrieve discounted amount based on promotion code input
func GetDiscount(db *sql.DB, promocode string) (float64, error) {
	var discount float64

	query := `
	SELECT discount_percentage 
	FROM Promotions 
	WHERE code = ? 
  	AND NOW() BETWEEN valid_from AND valid_until
	`
	err := db.QueryRow(query, promocode).Scan(&discount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("Invalid promotion code")
		}
		return 0, err
	}

	return discount, nil
}

// insert a billing record into database
func CreateBilling(db *sql.DB, userID int, reservationID int, finalPrice float64) error {
	// insert billing record into billing table
	result, err := db.Exec(`
		INSERT INTO Billing (reservation_id, user_id, amount) 
		VALUES (?, ?, ?)`,
		reservationID, userID, finalPrice)
	if err != nil {
		log.Printf("Error inserting billing record into database: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}

	// update reservation status to "Paid"
	updateResult, err := db.Exec(`
		UPDATE Reservations set status = 'Paid'
		WHERE user_id = ? AND id = ?
		`,
		userID, reservationID)
	if err != nil {
		log.Printf("Error inserting billing record into database: %v", err)
		return err
	}

	updateRowsAffected, err := updateResult.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}

	log.Printf("Billing record created successfully: RowsAffected=%d", rowsAffected)
	log.Printf("Billing record updated successfully: RowsAffected=%d", updateRowsAffected)
	return nil
}

// retrieve billing record
func RetrieveBilling(db *sql.DB, reservationID int, userID int) (Billing, error) {
	var billing Billing

	query := `SELECT id, reservation_id, user_id, amount, payment_status
	FROM Billing WHERE reservation_id = ? AND user_id = ?`

	err := db.QueryRow(query, reservationID, userID).Scan(&billing.ID, &billing.ReservationID, &billing.UserID, &billing.Amount, &billing.PaymentStatus)
	if err != nil {
		log.Printf("Error with select query: %v", err)
	}

	return billing, nil
}
