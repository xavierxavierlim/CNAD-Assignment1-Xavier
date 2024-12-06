package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// User struct represents a user in the system
type User struct {
	ID                 int
	Name               string
	Email              string
	Phone              string
	PasswordHash       string
	MembershipID       int
	RegistrationStatus string
}

type Membership struct {
	ID              int
	Tier            string
	HourlyRate      float32
	PriorityAccess  bool
	MaxBookingLimit int
}

// initialise the database connection
func InitDB(user, password, host, dbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// retrieve user by userID
func RetrieveUserByID(db *sql.DB, userID int) (User, error) {
	var user User
	err := db.QueryRow(`
        SELECT id, name, email, phone, membership_id, registration_status 
        FROM users WHERE id = ?`, userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.MembershipID, &user.RegistrationStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

// retrieve user by email or phone
func RetrieveUser(db *sql.DB, email, phone string) (User, string, error) {
	var user User
	var passwordHash string

	query := `
        SELECT id, name, email, phone, password_hash, membership_id, registration_status 
        FROM users WHERE email = ? OR phone = ?`
	err := db.QueryRow(query, email, phone).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &passwordHash, &user.MembershipID, &user.RegistrationStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, "", errors.New("user not found")
		}
		return user, "", err
	}
	return user, passwordHash, nil
}

// insert a new user into database
func CreateUser(db *sql.DB, newUser User) error {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	result, err := db.Exec(`
        INSERT INTO users (name, email, phone, password_hash, membership_id, registration_status) 
        VALUES (?, ?, ?, ?, ?, ?)`,
		newUser.Name, newUser.Email, newUser.Phone, string(hashedPassword), newUser.MembershipID, newUser.RegistrationStatus,
	)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}
	log.Printf("User created successfully: Email=%s, RowsAffected=%d", newUser.Email, rowsAffected)

	return nil
}

// update user's profile details
func UpdateUserProfile(db *sql.DB, userID int, updatedUser User, password string) error {
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = db.Exec(`
            UPDATE users SET name = ?, email = ?, phone = ?, password_hash = ? 
            WHERE id = ?`,
			updatedUser.Name, updatedUser.Email, updatedUser.Phone, string(hashedPassword), userID,
		)
		return err
	}

	_, err := db.Exec(`
        UPDATE users SET name = ?, email = ?, phone = ? 
        WHERE id = ?`,
		updatedUser.Name, updatedUser.Email, updatedUser.Phone, userID,
	)
	return err
}

/*
// retrieve all users
func GetAllUsers(db *sql.DB) ([]User, error) {
	var users []User
	rows, err := db.Query(`
        SELECT id, name, email, phone, membership_id, registration_status
        FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.MembershipID, &user.RegistrationStatus)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
*/

// retrieve membership details by id
func RetrieveMembershipByID(db *sql.DB, membershipID int) (Membership, error) {
	var membership Membership
	err := db.QueryRow(`
        SELECT tier, hourly_rate, priority_access, max_booking_limit
        FROM membership 
        WHERE id = ?`, membershipID).
		Scan(&membership.Tier, &membership.HourlyRate, &membership.PriorityAccess, &membership.MaxBookingLimit)
	if err != nil {
		if err == sql.ErrNoRows {
			return membership, errors.New("membership not found")
		}
		return membership, err
	}
	return membership, nil
}
