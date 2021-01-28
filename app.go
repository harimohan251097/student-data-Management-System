package main

// Imports
import (
	"bytes"
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// FoodItem ...
type StudentInfo struct {
	ID     int
	Name   string
	Class  string
	Branch string
}

type user struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

var userList []user

// GetAllStudents ...
func GetAllStudents(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: GetAllStudents")

	if r.Method == http.MethodGet {
		a := false
		username := r.FormValue("username")
		password := r.FormValue("password")

		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/hello")

		if err != nil {
			panic(err.Error())
		}
		result, err := db.Query("SELECT * from user5")
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()

		for result.Next() {
			var a1 user
			err := result.Scan(&a1.Username, &a1.Password)
			if err != nil {
				panic(err.Error())
			}
			userList = append(userList, a1)
		}
		for _, u := range userList {
			if u.Username == username && u.Password == password {
				a = true
			}
		}
		if a == false {
			http.Redirect(w, r, "/login", 301)
		}
	}
	// Post Request
	if r.Method == http.MethodPost {

		// Get the values from Form Data
		requestBody, err := json.Marshal(map[string]string{
			"name":   r.FormValue("name"),
			"class":  r.FormValue("class"),
			"branch": r.FormValue("branch"),
		})

		// Error Check
		if err != nil {
			fmt.Println("Error Occured")
			return
		}

		// Send the post request with the required Values
		resp, postErr := http.Post("http://localhost:8000/api/student/", "application/json", bytes.NewBuffer(requestBody))

		if postErr != nil {
			fmt.Println("Error Occured")
			return
		}

		fmt.Println(resp.Body)
	}

	// Fetch All the Food Items
	response, err := http.Get("http://localhost:8000/api/student")
	if err != nil {
		fmt.Printf("Could Not Fetch Foods, Error: %s", err)
		return
	}

	defer response.Body.Close()

	// Store the Fetched Food Items in an Array
	var items []StudentInfo
	_ = json.NewDecoder(response.Body).Decode(&items)

	// Template
	templ := template.Must(template.ParseFiles("templates/home.html"))

	templ.Execute(w, items)
	return
}

// updateStudentInfo ...
func updateStudentInfo(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: Updating Food Item")

	// Retrieve Food Item ID
	id := r.FormValue("id")

	// New Values
	requestBody, err := json.Marshal(map[string]string{
		"name":   r.FormValue("name"),
		"class":  r.FormValue("class"),
		"branch": r.FormValue("branch"),
	})

	// Error Check
	if err != nil {
		fmt.Println("Error Occured")
		return
	}

	// Specific URL for updating the required Item
	url := "http://localhost:8000/api/student/" + id + "/"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		fmt.Println("error occured")
		return
	}
	fmt.Println(resp.Body)

	// Redirect to Home Page
	http.Redirect(w, r, "/", 301)

}

func deleteStudentInfo(w http.ResponseWriter, r *http.Request) {

	// Send Post request to the URL with id to delete the Item
	id := r.FormValue("id")

	url := "http://localhost:8000/api/student/delete/" + id

	resp, err := http.Post(url, "application/json", nil)

	if err != nil {
		fmt.Println("error occured")
		return
	}

	fmt.Println(resp.Body)

	// Redirect to Home Page
	http.Redirect(w, r, "/", 301)
}

func loginStudentInfo(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("http://localhost:8000/api/student")
	if err != nil {
		fmt.Printf("Could Not Fetch Foods, Error: %s", err)
		return
	}

	defer response.Body.Close()

	// Store the Fetched Food Items in an Array
	var items []StudentInfo
	_ = json.NewDecoder(response.Body).Decode(&items)

	templ := template.Must(template.ParseFiles("templates/login.html"))

	templ.Execute(w, items)
	return
}

func main() {

	// Handler functions
	http.HandleFunc("/", GetAllStudents)
	http.HandleFunc("/update", updateStudentInfo)
	http.HandleFunc("/delete", deleteStudentInfo)
	http.HandleFunc("/login", loginStudentInfo)
	http.ListenAndServe(":8080", nil)
}
