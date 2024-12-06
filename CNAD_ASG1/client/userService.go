package controller

import (
	model "CNAD-ASSIGNMENT1-XAVIER/CNAD_ASG1/server"
	"database/sql"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// encodeCookie takes a string value, encodes it using base64 URL encoding and returns the encoded string
func encodeCookie(value string) string {
	// Convert the input string to a byte slice and encode it using base64 URL encoding.
	// base64.URLEncoding ensures the result is URL-safe (e.g., replacing '+' and '/' with URL-safe characters).
	return base64.URLEncoding.EncodeToString([]byte(value))
}

// decodeCookie takes a base64 URL-encoded string as input, decodes it and returns the decoded string along with any error encountered.
func decodeCookie(value string) (string, error) {
	// Decode the input base64 URL-encoded string into a byte slice.
	// base64.URLEncoding.DecodeString performs the decoding, which will return the decoded bytes and an error (if any).
	decoded, err := base64.URLEncoding.DecodeString(value)

	// Convert the decoded byte slice into a string and return it along with the error (if any).
	return string(decoded), err
}

// handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("CNAD_ASG1/view/login.html") // render the html file
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet { // check if the http method is get
		tmpl.Execute(w, nil)
		return
	}

	var loginUser model.User                         // create a user object
	loginMethod := r.FormValue("loginMethod")        // retrieve the login method: email or phone
	loginUser.PasswordHash = r.FormValue("password") // retrieve the password

	if loginMethod == "email" {
		loginUser.Email = strings.TrimSpace(r.FormValue("email")) // trim all whitespaces
	} else if loginMethod == "phone" {
		loginUser.Phone = strings.TrimSpace(r.FormValue("phone")) // trim all whitespaces
	}

	user, passwordHash, err := model.RetrieveUser(db, loginUser.Email, loginUser.Phone) // retrieve user object
	if err != nil {
		tmpl.Execute(w, map[string]interface{}{"Error": "Invalid credentials"}) // print our error message
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginUser.PasswordHash)) != nil { // hash the password
		tmpl.Execute(w, map[string]interface{}{"Error": "Invalid credentials"})
		return
	}

	// Set user_id session cookie
	cookie := &http.Cookie{
		Name:     "user_id",
		Value:    encodeCookie(strconv.Itoa(user.ID)),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	// redirect to home page upon successful login
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// handles home page
func HomeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("CNAD_ASG1/view/home.html") // render the home page
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("user_id") // retrieve user ID cookie
	if err != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	decodedUserID, err := decodeCookie(cookie.Value) // decode userID cookie
	if err != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	userID, err := strconv.Atoi(string(decodedUserID)) // convert userID from string to integer
	if err != nil {
		log.Println("Error converting user ID to integer:", err)
		return
	}

	user, err := model.RetrieveUserByID(db, userID) // retrieve the user by the user ID
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, map[string]interface{}{"User": user}) // print out the user details

}

// handles user profile
func ProfileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("CNAD_ASG1/view/profile.html") // render the profile page
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("user_id") // retrieve userID cookie
	if err != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	decodedUserID, err := decodeCookie(cookie.Value) // decode userID cookie
	if err != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	userID, err := strconv.Atoi(string(decodedUserID)) // convert userID from string to integer
	if err != nil {
		log.Println("Error with converting user ID to integer:", err)
		return
	}

	user, err := model.RetrieveUserByID(db, userID) // retrieve the user
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	membership, err := model.RetrieveMembershipByID(db, user.MembershipID) // retrieve the membership
	if err != nil {
		http.Error(w, "Failed to fetch membership details", http.StatusInternalServerError)
		return
	}

	// pass user and membership details to html template
	data := map[string]interface{}{
		"User":       user,
		"Membership": membership,
	}

	tmpl.Execute(w, data) // execute the html template with user and membership details
}

// handles update details
func UpdateDetailsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("CNAD_ASG1/view/updateDetails.html") // render html to update details
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("user_id") // retrieve user_id cookie
	if err != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	decodedUserID, err := decodeCookie(cookie.Value) // decode retrieved user_id cookie
	if err != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	userID, err := strconv.Atoi(string(decodedUserID)) // convert user_id from string to integer
	if err != nil {
		log.Println("Error converting user ID to integer:", err)
		return
	}

	if r.Method == http.MethodGet { // check if the http method is get
		user, err := model.RetrieveUserByID(db, userID) // retrieve user by the user_id
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		tmpl.Execute(w, map[string]interface{}{"User": user}) // pass user to html template
		return
	}

	// Handle POST update
	var updatedUser model.User                                  // new user object
	updatedUser.Name = r.FormValue("name")                      // retrieve name from form
	updatedUser.Email = strings.TrimSpace(r.FormValue("email")) // retrieve email from form
	updatedUser.Phone = r.FormValue("phone")                    // retrieve phone number from form
	password := r.FormValue("password")                         // retrieve password from form

	err = model.UpdateUserProfile(db, userID, updatedUser, password) // update the user details
	if err != nil {
		tmpl.Execute(w, map[string]interface{}{"Error": "Failed to update profile"}) // display error message
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther) // redirect to profile page upon successful update
}

// handles user register
func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("CNAD_ASG1/view/register.html") // render register html template
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet { // check if http method is get
		tmpl.Execute(w, nil) // execute the html template
		return
	}

	var newUser model.User                                  // create a User object
	newUser.Name = r.FormValue("name")                      // retrieve user name from form
	newUser.Email = strings.TrimSpace(r.FormValue("email")) // retrieve user email from form
	newUser.Phone = r.FormValue("phone")                    // retrieve phone from form
	newUser.PasswordHash = r.FormValue("password")          // retrieve password from form
	newUser.MembershipID = 1                                // membership tier basic
	newUser.RegistrationStatus = "Verified"                 // set registration status to registered

	// check if there are any blank fields
	if newUser.Name == "" || newUser.Email == "" || newUser.Phone == "" || newUser.PasswordHash == "" {
		tmpl.Execute(w, map[string]interface{}{"Error": "All fields are required"}) // display error message on html
		return
	}

	err = model.CreateUser(db, newUser) // create the new user
	if err != nil {
		tmpl.Execute(w, map[string]interface{}{"Error": "Failed to register user"}) // display error message
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther) // redirect to login page
}

func Logout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Clear the user_id cookie (or any other relevant cookies)
	cookies := r.Cookies()
	for _, cookie := range cookies {
		// Set the cookie expiration date to a past date to clear it
		http.SetCookie(w, &http.Cookie{
			Name:    cookie.Name,
			Value:   "", // Clear the value
			Path:    "/",
			Expires: time.Now().Add(-time.Hour), // Set expiration to 1 hour ago
			MaxAge:  -1,                         // Ensure the cookie is immediately expired
		})
	}

	// Optionally, redirect to the login page or home page after logout
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
