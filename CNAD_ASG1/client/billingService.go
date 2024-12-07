package controller

import (
	model "CNAD-ASSIGNMENT1-XAVIER/CNAD_ASG1/server"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// handles displaying reseravtion details and processes promotion code
func BillingHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl := template.Must(template.ParseFiles("CNAD_ASG1/view/bilingdetails.html")) // parse biling html

	if r.Method == http.MethodGet { // if user is viewing the form
		cookie, err := r.Cookie("user_id")    // retrieve user ID from cookie
		if err != nil || cookie.Value == "" { // if user is not logged in
			log.Println("Session error: user_id cookie not found")
			http.Redirect(w, r, "/", http.StatusSeeOther) // redirect user back to home page
			return
		}

		reservationCookie, err := r.Cookie("reservation_id") // retrieve reservationID from cookie
		if err != nil || reservationCookie.Value == "" {
			log.Println("Error retrieving reservation_id cookie:", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Parse reservation ID directly (no decoding needed)
		reservationID, err := strconv.Atoi(reservationCookie.Value) // convert reservationID from string to integer
		if err != nil {
			log.Println("Invalid reservation ID:", err)
			http.Error(w, "Invalid reservation ID.", http.StatusBadRequest)
			return
		}

		reservation, err := model.RetrieveReservationByID(db, reservationID) // retrieve reservation based on reservationID
		if err != nil {
			if err == sql.ErrNoRows { // if there is no reservation
				log.Println("Reservation not found:", err)
				http.Error(w, "Reservation not found.", http.StatusNotFound)
				return
			}
			log.Println("Error fetching reservation from database:", err)
			http.Error(w, "Failed to fetch reservation.", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, map[string]interface{}{ // execute the html template with reservation details
			"ID":            reservation.ID,
			"VehicleID":     reservation.VehicleID,
			"StartTime":     reservation.StartTime,
			"EndTime":       reservation.EndTime,
			"EstimatedCost": reservation.EstimatedCost,
		}); err != nil {
			log.Println("Error rendering template:", err)
			http.Error(w, "Failed to render form.", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost { // when user submits form - handle promotion code
		// Retrieve reservation details from form
		reservationIDStr := r.FormValue("reservation_id") // retrieve reservationID from form
		estimatedCostStr := r.FormValue("estimated_cost") // retrieve estimated cost from form
		promocode := r.FormValue("promocode")             // retrieve promotion code from form

		// handle promocode submission
		var discount float64
		var successMessage string
		var err error
		if promocode != "" {
			discount, err = model.GetDiscount(db, promocode) // retrieve discount percentage based on the promotion code input
			if err != nil {
				log.Println("Error fetching discount for promocode:", err)
				successMessage = "Failed to apply promocode. Please try again."
			} else { // success message to display in HTML form
				successMessage = fmt.Sprintf("Your discount percentage is %.2f%%.", discount)

			}
		} else {
			successMessage = "No promocode provided." // if promotion code submitted blank
		}

		reservationID, err := strconv.Atoi(string(reservationIDStr))   // convert reservationID from string to integer
		estimatedCost, err := strconv.ParseFloat(estimatedCostStr, 64) // convert estimated cost from string to float

		finalCost := ((100 - discount) / 100) * estimatedCost // calculate final cost based on discount percentage

		reservation, err := model.RetrieveReservationByID(db, reservationID) // retrieve rservation details
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("Reservation not found:", err)
				http.Error(w, "Reservation not found.", http.StatusNotFound)
				return
			}
			log.Println("Error fetching reservation from database:", err)
			http.Error(w, "Failed to fetch reservation.", http.StatusInternalServerError)
			return
		}

		// re-render the billing details template with the success message
		if err := tmpl.Execute(w, map[string]interface{}{
			"ID":            reservationID,
			"VehicleID":     reservation.VehicleID,
			"StartTime":     reservation.StartTime,
			"EndTime":       reservation.EndTime,
			"EstimatedCost": estimatedCost,
			"Success":       successMessage,
			"FinalCost":     finalCost,
		}); err != nil {
			log.Println("Error rendering template:", err)
			http.Error(w, "Failed to render form.", http.StatusInternalServerError)
			return
		}
	} else {
		// method not allowed
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}

// handles the payment of reservation bill
func PayHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("user_id")    // retrieve userID from cookie
	if err != nil || cookie.Value == "" { // if user not logged in
		log.Println("Session error: user_id cookie not found")
		http.Redirect(w, r, "/", http.StatusSeeOther) // redirect user back to login page
		return
	}

	decodedUserID, err := base64.StdEncoding.DecodeString(cookie.Value) // decode userID cookie
	if err != nil {
		log.Println("Error decoding user ID from cookie:", err)
		return
	}

	userID, err := strconv.Atoi(string(decodedUserID)) // convert userID from string to integer
	if err != nil {
		log.Println("Error converting user ID to integer:", err)
		return
	}

	reservationIDStr := r.URL.Query().Get("reservation_id") // retrieve reservationID from URL
	reservationID, err := strconv.Atoi(reservationIDStr)    // convert resreverationID to integer
	if err != nil {
		log.Println("Error converting reservation_id to int:", err)
		http.Error(w, "Invalid reservation ID.", http.StatusBadRequest)
		return
	}

	estimatedPriceStr := r.URL.Query().Get("estimated_cost") // retrieve estimated price from URL
	finalPriceStr := r.URL.Query().Get("final_cost")         // retrieve final price from URL
	var finalPrice float64
	var estimatedPrice float64
	if finalPriceStr != "" { // if final price has value
		finalPrice, err = strconv.ParseFloat(finalPriceStr, 64) // convert final price from string to float
		if err != nil {
			log.Println("Error converting finalPrice to float:", err)
			http.Error(w, "Invalid final price.", http.StatusBadRequest)
			return
		}
		model.CreateBilling(db, userID, reservationID, finalPrice) // create billing with final price
	} else { // if no final price as no or invalid promotion code
		estimatedPrice, err = strconv.ParseFloat(estimatedPriceStr, 64) // convert estimated price from string to float
		model.CreateBilling(db, userID, reservationID, estimatedPrice)  // create billing with estimated price
	}

	http.Redirect(w, r, "/user/reservations", http.StatusSeeOther) // redirect user back to reservation page
}

// handles invoice details
func InvoiceHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("user_id")    // retrieve userID from cookie
	if err != nil || cookie.Value == "" { // if user not logged in
		log.Println("Session error: user_id cookie not found")
		http.Redirect(w, r, "/", http.StatusSeeOther) // redirect user back to home page
		return
	}

	decodedUserID, err := base64.StdEncoding.DecodeString(cookie.Value) // decode userID cookie
	if err != nil {
		log.Println("Error decoding user ID from cookie:", err)
		return
	}

	userID, err := strconv.Atoi(string(decodedUserID)) // convert userID from string to integer
	if err != nil {
		log.Println("Error converting user ID to integer:", err)
		return
	}

	reservationIDStr := r.URL.Query().Get("reservation_id") // retrieve reservationID from URL
	reservationID, err := strconv.Atoi(reservationIDStr)    // convert reservationID from string to integer
	if err != nil {
		log.Printf("Invalid reservation_id: %v", err)
		http.Error(w, "Invalid reservation ID.", http.StatusBadRequest)
		return
	}

	vehicleIDStr := r.URL.Query().Get("vehicle_id") // retrieve vehicleID from URL
	vehicleID, err := strconv.Atoi(vehicleIDStr)    // convert vehicleID from string to integer
	if err != nil {
		log.Printf("Invalid vehicle_id: %v", err)
		http.Error(w, "Invalid vehicle ID.", http.StatusBadRequest)
		return
	}

	user, err := model.RetrieveUserByID(db, userID)                      // retrieve user record
	reservation, err := model.RetrieveReservationByID(db, reservationID) // retrieve reservation record
	vehicle, err := model.GetVehicleByID(db, vehicleID)                  // retrieve vehicle record
	billing, err := model.RetrieveBilling(db, reservationID, userID)     // retrieve billing record

	tmpl := template.Must(template.ParseFiles("CNAD_ASG1/view/invoice.html")) // render invoice html template
	if err := tmpl.Execute(w, map[string]interface{}{                         // execute html template with user, vehicle, reservation and billing details
		"BillingID":     billing.ID,
		"Amount":        billing.Amount,
		"PaymentStatus": billing.PaymentStatus,
		"VehicleModel":  vehicle.Model,
		"LicensePlate":  vehicle.LicensePlate,
		"Location":      vehicle.Location,
		"ChargeLevel":   vehicle.ChargeLevel,
		"UserName":      user.Name,
		"UserEmail":     user.Email,
		"UserPhone":     user.Phone,
		"MembershipID":  user.MembershipID,
		"StartTime":     reservation.StartTime,
		"EndTime":       reservation.EndTime,
	}); err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Failed to render form.", http.StatusInternalServerError)
		return
	}
}
