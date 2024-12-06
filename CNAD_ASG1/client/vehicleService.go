package controller

import (
	model "CNAD-ASSIGNMENT1-XAVIER/CNAD_ASG1/server"
	"database/sql"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"time"
)

// handles retrieval of available vehicles within a selected start date and end date
func AvailableVehiclesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("user_id")    // retrieve user_id cookie
	if err != nil || cookie.Value == "" { // redirect to login page if user is not logged in
		log.Println("User ID not found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	decodedUserID, err := base64.StdEncoding.DecodeString(cookie.Value) // decode userID cookie
	if err != nil {
		log.Println("Error decoding cookie value:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	log.Printf("UserID cookie: %s\n", string(decodedUserID)) // debug

	currentDateTime := time.Now().Format("2006-01-02T15:04") // current time

	// retrieve start date and end date from URL
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var vehicles []model.Vehicle // vehicles list
	var errorMessage string

	if startDate != "" && endDate != "" { // check if startDate and endDate fields are empty

		startDateParsed, err := time.Parse("2006-01-02T15:04", startDate) // format the start date
		if err != nil {
			errorMessage = "Invalid start date format."
			log.Println(errorMessage)
		}

		endDateParsed, err := time.Parse("2006-01-02T15:04", endDate) // format the end date
		if err != nil {
			errorMessage = "Invalid end date format."
			log.Println(errorMessage)
		}

		// check if end_date is earlier than start_date
		if errorMessage == "" && !endDateParsed.After(startDateParsed) {
			errorMessage = "End date must be after start date."
			log.Println(errorMessage)
		}

		// retrieve available vehicles if endDate > startDate (valid)
		if endDateParsed.After(startDateParsed) {
			vehicles, err = model.GetAvailableVehicles(db, startDate, endDate) // retrieve available vehicles
			if err != nil {
				log.Println("Error fetching available vehicles:", err)
				http.Error(w, "Failed to retrieve available vehicles", http.StatusInternalServerError)
				return
			}
		}
	}
	// Parse the HTML template
	tmpl, err := template.ParseFiles("CNAD_ASG1/view/vehicles.html") // render html template
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// pass vehicle data to template
	data := map[string]interface{}{
		"Vehicles":        vehicles,
		"StartDate":       startDate,
		"EndDate":         endDate,
		"ErrorMessage":    errorMessage,
		"CurrentDateTime": currentDateTime,
	}

	log.Println("Rendering template with data:", data)

	err = tmpl.Execute(w, data) // execute html template with vehicle data
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
