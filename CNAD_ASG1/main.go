package main

import (
	controller "CNAD-ASSIGNMENT1-XAVIER/CNAD_ASG1/client"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv" // load environment variables from a .env file
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env")
		return
	}

	// Get database connection details from environment variables
	dbUser, dbPassword, dbHost, dbName := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")

	// Check if all necessary environment variables are set
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbName == "" {
		log.Fatal("Database credentials not fully set in environment variables")
	}

	// Build DSN for database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Database connection error:", err)
		// Log the error and exit, as there's no HTTP request context here
		log.Fatal("Failed to connect to the database")
	}
	defer db.Close()

	// HTTP Handlers
	// login user
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controller.LoginHandler(w, r, db)
	})
	// home page
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		controller.HomeHandler(w, r, db)
	})
	// register
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controller.RegisterHandler(w, r, db)
	})
	// user profile
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		controller.ProfileHandler(w, r, db)
	})
	// update user details
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		controller.UpdateDetailsHandler(w, r, db)
	})
	// user logout
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		controller.Logout(w, r, db)
	})
	// display available vehicles
	http.HandleFunc("/availableVehicles", func(w http.ResponseWriter, r *http.Request) {
		controller.AvailableVehiclesHandler(w, r, db)
	})
	// make reservation
	http.HandleFunc("/makeReservation", func(w http.ResponseWriter, r *http.Request) {
		controller.ReservationHandler(w, r, db)
	})
	// display reservation to user
	http.HandleFunc("/user/reservations", func(w http.ResponseWriter, r *http.Request) {
		controller.ViewUserReservationsHandler(w, r, db)
	})
	// update reservation
	http.HandleFunc("/user/reservations/modify", func(w http.ResponseWriter, r *http.Request) {
		controller.ModifyReservationHandler(w, r, db)
	})
	// cancel reservation
	http.HandleFunc("/user/reservations/cancel", func(w http.ResponseWriter, r *http.Request) {
		controller.CancelReservationHandler(w, r, db)
	})
	// display billing details
	http.HandleFunc("/user/reservations/billing", func(w http.ResponseWriter, r *http.Request) {
		controller.BillingHandler(w, r, db)
	})
	// pay billing
	http.HandleFunc("/user/reservations/billing/pay", func(w http.ResponseWriter, r *http.Request) {
		controller.PayHandler(w, r, db)
	})
	// display invoice
	http.HandleFunc("/user/reservations/invoice", func(w http.ResponseWriter, r *http.Request) {
		controller.InvoiceHandler(w, r, db)
	})

	// Start server
	log.Println("Server running on http://localhost:5001")
	log.Fatal(http.ListenAndServe(":5001", nil))
}
