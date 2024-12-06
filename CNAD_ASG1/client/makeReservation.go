package controller

import (
	model "CNAD-ASSIGNMENT1-XAVIER/CNAD_ASG1/server"
	"database/sql"
	"encoding/base64"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

// handles displaying available vehicles and making reservations
func ReservationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("user_id")    // retrieve userID from cookie
	if err != nil || cookie.Value == "" { // redirect to login page if user is not logged in
		log.Println("User ID not found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	decodedUserID, err := base64.StdEncoding.DecodeString(cookie.Value) // decode userID from cookie
	if err != nil {
		log.Println("Error decoding user ID from cookie:", err)
		return
	}

	userID, err := strconv.Atoi(string(decodedUserID)) // convert userID from string to integer
	if err != nil {
		log.Println("Error converting user ID to integer:", err)
		return
	}

	tmpl, err := template.ParseFiles("CNAD_ASG1/view/makeReservation.html") // render the make reservation html
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet { // retrieve the form
		vehicleIDStr := r.URL.Query().Get("vehicle_id") // retrieve vehicleID from URL

		if vehicleIDStr == "" { // if no vehicleID selected
			tmpl.Execute(w, map[string]interface{}{
				"Error": "No vehicle selected for reservation",
			})
			return
		}

		vehicleID, err := strconv.Atoi(vehicleIDStr) // convert vehicleID from string to integer
		if err != nil {
			log.Printf("Error converting vehicle ID to integer: %s\n", err)
			tmpl.Execute(w, map[string]interface{}{
				"Error": "Invalid vehicle ID.",
			})
			return
		}

		vehicle, err := model.GetVehicleByID(db, vehicleID) // retrieve vehicle by vehicleID
		if err != nil {
			log.Printf("Error fetching vehicle details: %s\n", err)
			tmpl.Execute(w, map[string]interface{}{
				"Error": "Failed to fetch vehicle details: " + err.Error(),
			})
			return
		}

		startDateStr := r.URL.Query().Get("start_date") // retrieve start date from URL
		if startDateStr == "" {
			tmpl.Execute(w, map[string]interface{}{
				"Error": "No start date selected",
			})
			return
		}

		endDateStr := r.URL.Query().Get("end_date") // retrieve end date from URL
		if endDateStr == "" {
			tmpl.Execute(w, map[string]interface{}{
				"Error": "No end date selected",
			})
			return
		}

		// parse start date and end date strings into time.Time objects
		startDate, err := time.Parse("2006-01-02T15:04", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date and time format", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02T15:04", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date and time format", http.StatusBadRequest)
			return
		}

		hours := int(math.Ceil(endDate.Sub(startDate).Hours()))              // calculate the number of hours between the dates
		estimatedCost, err := model.RetrieveEstimatedCost(db, userID, hours) // calculate estimated cost by retrieving user membership details

		// format dates for the HTML form
		formattedStartDate := startDate.Format("2006-01-02T15:04")
		formattedEndDate := endDate.Format("2006-01-02T15:04")

		// pass vehicle details and formatted dates to reservation html page
		tmpl.Execute(w, map[string]interface{}{
			"Vehicle":       vehicle,
			"StartDate":     formattedStartDate,
			"EndDate":       formattedEndDate,
			"EstimatedCost": estimatedCost,
		})
		return
	}

	if r.Method == http.MethodPost { // if user submitted the form
		cookie, err := r.Cookie("user_id")    // retrieve userID from cookie
		if err != nil || cookie.Value == "" { // redirect user to login page if not logged in
			log.Println("User ID not found")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		decodedUserID, err := base64.StdEncoding.DecodeString(cookie.Value) // decode userID from cookie
		if err != nil {
			log.Println("Error decoding user ID from cookie:", err)
			return
		}

		userID, err := strconv.Atoi(string(decodedUserID)) // convert userID from string to integer
		if err != nil {
			log.Println("Error converting user ID to integer:", err)
			return
		}

		vehicleIDStr := r.FormValue("vehicle_id")         // retrieve vehicleID from form
		startDateStr := r.FormValue("start_date")         // retrieve start date from form
		endDateStr := r.FormValue("end_date")             // retrieve end date from form
		estimatedCostStr := r.FormValue("estimated_cost") // retrieve estimated cost from form

		// check if there are any empty fields
		if vehicleIDStr == "" || startDateStr == "" || endDateStr == "" || estimatedCostStr == "" {
			tmpl.Execute(w, map[string]interface{}{
				"Error": "All fields are required to make a reservation",
			})
			return
		}

		vehicleID, err := strconv.Atoi(vehicleIDStr) // convert vehicleID from string to integer
		if err != nil {
			log.Printf("Error converting vehicle ID to integer: %s\n", err)
			return
		}

		startDate, err := time.Parse("2006-01-02T15:04", startDateStr) // parse start date with time
		if err != nil {
			log.Println("Error parsing start date:", err)
			return
		}

		endDate, err := time.Parse("2006-01-02T15:04", endDateStr) // parse end date with time
		if err != nil {
			log.Println("Error parsing end date:", err)
			tmpl.Execute(w, map[string]interface{}{
				"Error": "Invalid end date format. Please use the correct format (YYYY-MM-DDTHH:MM).",
			})
			return
		}

		estimatedCost, err := strconv.ParseFloat(estimatedCostStr, 64) // convert estimated cost from string to float
		if err != nil {
			log.Printf("Error converting estimated cost to float: %s\n", err)
			return
		}

		reservationID, err := model.CreateReservation(db, userID, vehicleID, startDate, endDate, estimatedCost) // create reservation
		if err != nil {
			log.Printf("Error creating reservation: %s\n", err)
			return
		}

		log.Printf("Successfully created reservation ID: %d\n", reservationID) // debug

		http.Redirect(w, r, "/home", http.StatusSeeOther) // redirect user to home page
	}
}
