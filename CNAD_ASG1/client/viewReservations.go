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

// handles viewing of all reservations user has
func ViewUserReservationsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("user_id")    // retrieve userID from cookie
	if err != nil || cookie.Value == "" { // if user is not logged in
		log.Println("Session error: user_id cookie not found")
		http.Redirect(w, r, "/", http.StatusSeeOther) // redirect user to login page
		return
	}

	decodedUserID, err := base64.StdEncoding.DecodeString(cookie.Value) // decode userID from cookie
	if err != nil {
		log.Println("Error decoding cookie value:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := strconv.Atoi(string(decodedUserID)) // convert userID from string to integer
	if err != nil {
		log.Println("Invalid user ID:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	reservations, err := model.RetrieveReservations(db, userID) // retrieve reservations based on userID
	if err != nil {
		log.Println("Error retrieving user reservations:", err)
		http.Error(w, "Failed to fetch reservations", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("CNAD_ASG1/view/userReservations.html")) // render user reservation html
	if err := tmpl.Execute(w, reservations); err != nil {                              // execute the user reservation template with all reservations
		log.Println("Error rendering reservations template:", err)
		http.Error(w, "Failed to render reservations", http.StatusInternalServerError)
		return
	}
}

// handles updating reservation details
func ModifyReservationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet { // if user is viewing the form
		cookie_user, err := r.Cookie("user_id")    // retrieve userID from cookie
		if err != nil || cookie_user.Value == "" { // if user is not logged in
			log.Println("Session error: user_id cookie not found")
			http.Redirect(w, r, "/", http.StatusSeeOther) // redirect user back to login page
			return
		}

		cookie, err := r.Cookie("reservation_id") // retrieve reservationID from cookie
		if err != nil || cookie.Value == "" {
			log.Println("Error retrieving reservation_id cookie:", err)
			http.Error(w, "Reservation ID cookie not found.", http.StatusBadRequest)
			return
		}

		reservationID, err := strconv.Atoi(cookie.Value) // convert reservationiD from string to integer
		if err != nil {
			log.Println("Invalid reservation ID:", err)
			http.Error(w, "Invalid reservation ID.", http.StatusBadRequest)
			return
		}

		reservation, err := model.RetrieveReservationByID(db, reservationID) // retrieve reservation to update
		if err != nil {
			if err == sql.ErrNoRows { // if not able to retrieve reservations
				log.Println("Reservation not found:", err)
				http.Error(w, "Reservation not found.", http.StatusNotFound)
				return
			}
			log.Println("Error fetching reservation from database:", err)
			http.Error(w, "Failed to retrieve reservation.", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("CNAD_ASG1/view/modifyReservation.html")) // render modify reservation html template
		if err := tmpl.Execute(w, reservation); err != nil {                                // execute the html template with the reservation details
			log.Println("Error rendering template:", err)
			http.Error(w, "Failed to render form.", http.StatusInternalServerError)
			return
		}

	} else if r.Method == http.MethodPost { // if user submits the form
		userCookie, err := r.Cookie("user_id")    // retrieve userID from cookie
		if err != nil || userCookie.Value == "" { // if user not logged in
			log.Println("Error retrieving user_id cookie:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther) // redirect user back to login page
			return
		}

		decodedUserID, err := base64.StdEncoding.DecodeString(userCookie.Value) // decode userID from cookie
		if err != nil {
			log.Println("Error decoding user ID from cookie:", err)
			return
		}

		userID, err := strconv.Atoi(string(decodedUserID)) // convert userID from string to integer
		if err != nil {
			log.Println("Error converting user ID to integer:", err)
			return
		}

		cookie, err := r.Cookie("reservation_id") // retrieve reservationID from cookie
		if err != nil || cookie.Value == "" {
			log.Println("Error retrieving reservation_id cookie:", err)
			http.Error(w, "Reservation ID cookie not found.", http.StatusBadRequest)
			return
		}

		reservationID, err := strconv.Atoi(cookie.Value) // convert reservationID from string to integer
		if err != nil {
			log.Println("Invalid reservation ID in cookie:", err)
			http.Error(w, "Invalid reservation ID.", http.StatusBadRequest)
			return
		}

		currentDateTimeStr := time.Now().Format("2006-01-02T15:04")                // retrieve current time
		currentDateTime, err := time.Parse("2006-01-02T15:04", currentDateTimeStr) // parse current time
		if err != nil {
			log.Println("Error parsing current date time:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		startTimeStr := r.FormValue("start_time") // retrieve start time from URL

		if startTimeStr == "" { // if no start time in URL
			log.Println("Start time is required")
			http.Error(w, "Start time is required.", http.StatusBadRequest)
			return
		}

		startTime, err := time.Parse("2006-01-02T15:04", startTimeStr) // parse start time
		if err != nil {
			log.Println("Invalid start time:", err)
			http.Error(w, "Invalid start time format.", http.StatusBadRequest)
			return
		}

		endTimeStr := r.FormValue("end_time") // retrieve end time from URL

		if endTimeStr == "" { // if no end time in URL
			log.Println("End time is missing")
			http.Error(w, "End time is required.", http.StatusBadRequest)
			return
		}

		endTime, err := time.Parse("2006-01-02T15:04", endTimeStr) // parse end time
		if err != nil {
			log.Println("Invalid end time:", err)
			http.Error(w, "Invalid end time format.", http.StatusBadRequest)
			return
		}

		if currentDateTime.After(startTime) { // if start time is before today's date
			http.Redirect(w, r, "/user/reservations", http.StatusSeeOther)
			return
		}

		if currentDateTime.After(endTime) { // if end time is before today's date
			http.Redirect(w, r, "/user/reservations", http.StatusSeeOther)
			return
		}

		if !endTime.After(startTime) { // if start time is after end time
			http.Redirect(w, r, "/user/reservations", http.StatusSeeOther)
			return
		}

		hours := int(math.Ceil(endTime.Sub(startTime).Hours()))              // calculate number of hours of reservation
		estimatedCost, err := model.RetrieveEstimatedCost(db, userID, hours) // calculate estimated cost based on membership tier and hours

		err = model.ModifyReservation(db, reservationID, startTime, endTime, estimatedCost) // update reservation in the database
		if err != nil {
			log.Println("Error updating reservation in database:", err)
			http.Error(w, "Failed to update reservation.", http.StatusInternalServerError)
			return
		}

		// clear reservation_id cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "reservation_id",
			Value:  "",
			Path:   "/",
			MaxAge: -1, // delete cookie
		})

		http.Redirect(w, r, "/user/reservations", http.StatusSeeOther) // redirect user back to reservations page
	}
}

// handles cancelling reservation
func CancelReservationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost { // if user submits the form
		reservationID, _ := strconv.Atoi(r.FormValue("reservation_id")) // retrieves reservationID from form

		err := model.CancelReservation(db, reservationID) // changes status of reservation to "Cancelled"
		if err != nil {
			log.Println("Error canceling reservation:", err)
			http.Error(w, "Failed to cancel reservation", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/user/reservations", http.StatusSeeOther) // redirect user back to reservation page
	}
}
