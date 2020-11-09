package controller

import (
	"encoding/json"
	"net/http"

	"caringAPI/config"
)

// User Model
//Girish comment-  TBD : Need to move all model in seperate folder
type User struct {
	Username string `form:"username" json:"username"`
	Dob      string `form:"dob" json:"dob"`
	Age      string `form:"age" json:"age"`
	Email    string `form:"email" json:"email"`
	Phone    string `form:"phone" json:"phone"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	db := config.Connect()
	defer db.Close()
	username := r.FormValue("username")
	selecteduser, err := db.Query("SELECT * FROM users WHERE username = ? ", username)
	if err != nil {
		panic(err.Error())
	}
	usr := User{}
	res := []User{}
	for selecteduser.Next() {

		var username, dob, age, email, phone string
		err = selecteduser.Scan(&username, &dob, &age, &email, &phone)
		if err != nil {
			panic(err.Error())
		}
		usr.Username = username
		usr.Dob = dob
		usr.Age = age
		usr.Email = email
		usr.Phone = phone
		res = append(res, usr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("status", "200")
	w.Header().Set("message", "User data has been fetched")
	json.NewEncoder(w).Encode(res)

}

func Proxy(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")

	if username != "" {
		// Database check for user data!
		result, err := http.Get("http://localhost:1234/auth?username=" + username)
		//result, err := http.Get("/auth?username=" + username)
		if err != nil {
			panic(err)
		}
		defer result.Body.Close()
		//json.NewEncoder(w).Encode(result.Status)
		if result.Status == "200 OK" {
			http.Redirect(w, r, "/user/profile?username="+username, 302) //StatusFound
		} else if result.Status == "401 Unauthorized" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("User request is unauthorized")
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Bad request")
	}

}

func Authenticator(w http.ResponseWriter, r *http.Request) {

	db := config.Connect()
	defer db.Close()
	username := r.FormValue("username")
	rows, err := db.Query("SELECT count(*) FROM users WHERE username = ? ", username)
	if err != nil {
		panic(err.Error())
	}
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			panic(err.Error())
		}
	}
	if count > 0 {
		w.WriteHeader(http.StatusOK)

	} else {
		w.WriteHeader(http.StatusUnauthorized)

	}

}
