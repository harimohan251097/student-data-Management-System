package main

// Imports
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Initiate DB and Error variable
var db *gorm.DB
var err error

// Model ...

// Food Model
type Login struct {
	ID       uint   `gorm:"primary_key"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Student struct {
	ID     uint   `gorm:"primary_key"`
	Name   string `json:"name"`
	Class  string `json:"class"`
	Branch string `json:"branch"`
}

// Initialize Database with Food Table
func initialMigration() {
	dsn := "host=localhost user=postgres password=root dbname=hari sslmode=disable"
	// Check For Connection
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err.Error())
		panic("Connection to Database Failed")
	}

	defer db.Close()

	db.AutoMigrate(&Student{})
	db.AutoMigrate(&Login{})
}

// Get All Students
func getStudents(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: getStudents")

	dsn := "host=localhost user=postgres password=root dbname=hari  sslmode=disable"
	// Check For Connection
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed To Retrieve Foods, Database Connection Failed")
	}
	defer db.Close()

	// Find All students
	var students []Student
	db.Find(&students)

	// Return All students
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func getloginStudents(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: getStudents")

	dsn := "host=localhost user=postgres password=root dbname=hari  sslmode=disable"
	// Check For Connection
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed To Retrieve Foods, Database Connection Failed")
	}
	defer db.Close()

	// Find All students
	var login []Login
	db.Find(&login)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(login)
}

func createloginStudents(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: getStudents")

	dsn := "host=localhost user=postgres password=root dbname=hari  sslmode=disable"
	// Check For Connection
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed To Retrieve Foods, Database Connection Failed")
	}
	defer db.Close()

	var login Login
	_ = json.NewDecoder(r.Body).Decode(&login)

	// Create student
	db.Create(&Login{Username: login.Username, Password: login.Password})
	json.NewEncoder(w).Encode(login)

}

// Get specific student
func getStudent(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: getStudent")

	dsn := "host=localhost user=postgres password=root dbname=hari  sslmode=disable"
	// Check For Connection
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed To Retrieve Food, Database Connection Failed")
	}
	defer db.Close()

	// Retrieve ID from url variables
	vars := mux.Vars(r)
	id := vars["id"]

	// Find student with the given ID
	var student Student
	db.Where("id = ?", id).Find(&student)

	// Return The student
	json.NewEncoder(w).Encode(student)
}

// Create student
func createStudent(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: createStudent")

	dsn := "host=localhost user=postgres password=root dbname=hari  sslmode=disable"
	// Check For Connection
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		panic("Failed To Create Student, Database Connection Failed")
	}
	defer db.Close()

	// Decode the raw Data from response Body into student object
	var student Student
	_ = json.NewDecoder(r.Body).Decode(&student)

	// Create student
	db.Create(&Student{Name: student.Name, Class: student.Class, Branch: student.Branch})

	fmt.Fprintf(w, "Sucessfuly Created Food Item")
}

// Delete student
func deleteStudent(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: deleteStudent")
	dsn := "host=localhost user=postgres password=root dbname=hari  sslmode=disable"
	// Check For Connection
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		panic("Failed To Delete, Database Connection Failed")
	}
	defer db.Close()

	// Retrieve ID from url variables
	vars := mux.Vars(r)
	id := vars["id"]

	// Find Student with the Given ID
	var student Student
	db.Where("id = ?", id).Find(&student)

	// Delete the Student
	db.Delete(&student)

	fmt.Fprintf(w, "Sucessfuly Deleted student with ID :"+id)
}

// Update Student
func updateStudent(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: updateStudent")

	dsn := "host=localhost user=postgres password=root dbname=hari  sslmode=disable"
	// Check For Connection
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		panic("Failed to Update Student, Database Connection Failed")
	}
	defer db.Close()

	// Retrieve ID from url variables
	vars := mux.Vars(r)
	id := vars["id"]

	// Decode new Updated Item into Temporary food Object
	var Tempstudent Student
	_ = json.NewDecoder(r.Body).Decode(&Tempstudent)

	// Find the student to be Updated Using ID
	var student Student
	db.Where("id = ?", id).Find(&student)

	// Apply changes to the Student
	student.Name = Tempstudent.Name
	student.Class = Tempstudent.Class
	student.Branch = Tempstudent.Branch

	// Save it into Database
	db.Save(&student)

	fmt.Fprintf(w, "Sucessfuly Updated Student with ID: "+id)
}

// Main Function
func main() {

	// Initialize Database
	initialMigration()

	// Router
	router := mux.NewRouter()

	// Router Handlers
	router.HandleFunc("/api/student", getStudents).Methods("GET")
	router.HandleFunc("/api/student/login", getloginStudents).Methods("GET")
	router.HandleFunc("/api/student/login", createloginStudents).Methods("POST")
	router.HandleFunc("/api/student/{id}", getStudent).Methods("GET")
	router.HandleFunc("/api/student/", createStudent).Methods("POST")
	router.HandleFunc("/api/student/{id}/", updateStudent).Methods("POST")
	router.HandleFunc("/api/student/delete/{id}", deleteStudent).Methods("POST")

	// Listen
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
