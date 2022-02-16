package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// Initialize variable to store active session state
var activeSession bool = false

// Initializing class for a user account
type User struct {
	ID       int    `json:"userID"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Initializing class for a user's task
type Task struct {
	ID       int    `json:"taskID"`
	UserID   int    `json:"userID"`
	Name     string `json:"name"`
	Priority string `json:"priority"`
	Status   string `json:"status"`
}

// Initializing class for database records
type Database struct {
	NextUserID int    `json:"nextUserID"`
	NextTaskID int    `json:"nextTaskID"`
	Recs_user  []User `json:"users"`
	Recs_task  []Task `json:"tasks"`
	Mu         sync.Mutex
}

// Initialize logging functions
// Reference: https://www.honeybadger.io/blog/golang-logging/
var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

// Initializing file for logging
func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Function to initialize database from data.json
func initializeDB(filepath string) *Database {
	file, err := os.Open(filepath)
	if err != nil {
		ErrorLogger.Println(err.Error())
	}
	defer file.Close()
	byteData, err := ioutil.ReadAll(file)
	if err != nil {
		ErrorLogger.Println(err.Error())
	}
	var result *Database
	if err := json.Unmarshal([]byte(byteData), &result); err != nil {
		ErrorLogger.Println(err.Error())
	}
	result.Mu = sync.Mutex{}
	return result
}

// Function to save database to data.json
func updateDatabase(db *Database) {
	db.Mu.Lock()
	jsonString, err := json.Marshal(db)
	if err != nil {
		ErrorLogger.Println(err.Error())
		return
	}
	if err := ioutil.WriteFile("data/data.json", jsonString, os.ModePerm); err != nil {
		ErrorLogger.Println(err.Error())
		return
	}
	db.Mu.Unlock()

}

func main() {
	db := initializeDB("data/data.json")
	// db := &Database{Recs_user: []User{}, Recs_task: []Task{}}
	if err := http.ListenAndServe(":8080", db.handler()); err != nil {
		log.Fatalf("Error ListenAndServe(): %s", err.Error())
	}
}

func (db *Database) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int
		if r.URL.Path == "/tasks" {
			userID := r.URL.Query().Get("userid")
			if u, err := strconv.Atoi(userID); err == nil {
				db.processTasks(u, w, r)
			}
		} else if n, _ := fmt.Sscanf(r.URL.Path, "/tasks/%d", &id); n == 1 {
			userID := r.URL.Query().Get("userid")
			if u, err := strconv.Atoi(userID); err == nil {
				db.processTaskID(u, id, w, r)
			}
		} else if r.URL.Path == "/login" {
			db.createSession(w, r)
		} else if r.URL.Path == "/logout" {
			db.endSession(w, r)
		} else if r.URL.Path == "/signup" {
			db.createUser(w, r)

		}
		updateDatabase(db)
	}
}

func (db *Database) createSession(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.Header().Set("Content-Type", "application/json")
		var userid int
		var username string
		// Check if there is a user logged in
		if activeSession { // Ongoing active session, log in fail
			fmt.Fprintln(w, "Already logged in. Sign out to switch to another account.")
		} else { // No active session, proceed to attempt login
			var rec User
			if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				ErrorLogger.Println(err.Error())
				return
			}

			exists := false
			// Access database
			db.Mu.Lock()

			// Check if User record exists
			for _, item := range db.Recs_user {
				if (rec.Username == item.Username) && (rec.Password == item.Password) {
					exists = true
					userid = item.ID
					username = item.Username
				}
			}

			db.Mu.Unlock()

			if exists { // Username and password correct, login successful
				res := make(map[string]int)
				res["id"] = userid
				jsonResp, err := json.Marshal(res)
				if err != nil {
					ErrorLogger.Println(err.Error())
				}
				w.Write(jsonResp)
				activeSession = true
				InfoLogger.Printf("Login sucessful for user ID: %d (%s)\n", userid, username)
				return
			} else { // Username and password incorrect, login unsuccessful
				http.Error(w, "Username and/or password is incorrect.", http.StatusNotFound)
				ErrorLogger.Println("Username and/or password is incorrect.")
				return
			}

		}

	case "GET": // Get method not allowed for login
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Get method not allowed")
		return
	case "PUT": // Put method not allowed for login
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Put method not allowed")
		return

	case "DELETE": // Delete method not allowed for login
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Get method not allowed")
		return
	}

}

func (db *Database) endSession(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// Check if there is a user logged in
		if activeSession { // Active session present; log out success
			activeSession = false
			fmt.Fprintln(w, "Logout sucesssful.")
			InfoLogger.Println("Logout sucessful.")
		} else { // No active session; log out fail
			fmt.Fprintln(w, "Already logged out.")
			WarningLogger.Println("Attempted to log out without an active session")
		}
	case "GET": // Get method not allowed for logout
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Get method not allowed")
		return
	case "PUT": // Put method not allowed for logout
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Put method not allowed")
		return
	case "DELETE": // Delete method not allowed for logout
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Delete method not allowed")
		return
	}
}

func (db *Database) createUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.Header().Set("Content-Type", "application/json")
		var rec User
		if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			ErrorLogger.Println(err.Error())
			return
		}

		exists := false
		// Access database
		db.Mu.Lock()

		// Check if Username exists
		for _, item := range db.Recs_user {
			if rec.Username == item.Username {
				exists = true
			}
		}

		if exists { // Username already exists, user creation unusuccesful successful
			http.Error(w, "Username already exists", http.StatusNotFound)
			ErrorLogger.Println("Username already exists")
			db.Mu.Unlock()
			return
		} else { // Username does not exist, user creation successful
			rec.ID = db.NextUserID
			db.NextUserID++
			db.Recs_user = append(db.Recs_user, rec)
			db.Mu.Unlock()
			InfoLogger.Printf("User creation sucessful for user ID: %d (%s)\n", rec.ID, rec.Username)
		}

	case "GET": // Get method not allowed for user creation
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Get method not allowed")
		return
	case "PUT": // Put method not allowed fo  user creation
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Put method not allowed")
		return
	case "DELETE": // Delete method not allowed for user creation
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Delete method not allowed")
		return
	}

}

func (db *Database) processTasks(userID int, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.Header().Set("Content-Type", "application/json")
		var rec Task
		if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			ErrorLogger.Println(err.Error())
			return
		}

		exists := false
		// Access database
		db.Mu.Lock()

		// Check if record exists
		for _, item := range db.Recs_task {
			if (rec.Name == item.Name) && (rec.Priority == item.Priority) && (rec.Status == item.Status) && (userID == item.UserID) {
				exists = true
			}
		}
		if exists { // Task record exists, show error
			http.Error(w, "Task already exists", http.StatusForbidden)
			ErrorLogger.Println("Task already exists")
		} else { // Task record does not exist, create task record
			rec.ID = db.NextTaskID
			rec.UserID = userID
			db.NextTaskID++
			db.Recs_task = append(db.Recs_task, rec)
			fmt.Fprintln(w, "{\"success\": true}")
			InfoLogger.Printf("Task creation sucessful (task ID %d) for user ID %d\n", rec.ID, rec.UserID)
		}

		db.Mu.Unlock()

	case "GET": // Get all tasks belonging to the user
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(filterTasksByUser(db.Recs_task, userID)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			ErrorLogger.Println(err.Error())
			return
		}
	case "PUT": // Put method not allowed for /tasks
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Put method not allowed")
		return

	case "DELETE": // Put method not allowed for /tasks
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Delete method not allowed")
		return
	}

}

// Function to get all tasks by logged in user
func filterTasksByUser(task_list []Task, userID int) []Task {
	var filtered_list []Task
	for _, task := range task_list {
		if task.UserID == userID {
			filtered_list = append(filtered_list, task)
		}
	}
	return filtered_list
}

func (db *Database) processTaskID(userID int, id int, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE": // Delete contact record
		w.Header().Set("Content-Type", "application/json")
		exists := false
		db.Mu.Lock()
		for j, item := range db.Recs_task { // Loop through all records to check if task ID exists
			if id == item.ID { // Task ID exists, check if task belongs to logged in user
				if userID == item.UserID {
					db.Recs_task = append(db.Recs_task[:j], db.Recs_task[j+1:]...)
					exists = true
					InfoLogger.Printf("Task deletion sucessful (task ID %d) for user ID %d\n", id, userID)
					break
				} else {
					http.Error(w, "You are not authorized to delete this record.", http.StatusInternalServerError)
					ErrorLogger.Println("Unauthorized access to record")
				}
			}
		}
		db.Mu.Unlock()
		if exists {
			fmt.Fprintln(w, "{\"success\": true}")
		} else {
			fmt.Fprintln(w, "{\"success\": false}")
			http.Error(w, "Record not found", http.StatusForbidden)
			ErrorLogger.Println("Record not found")
		}
	case "POST": // Post method not allowed for /tasks/
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		ErrorLogger.Println("Post method not allowed")
		return
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		exists := false
		allowed := false
		rec := Task{}
		db.Mu.Lock()
		// Loop through all records to check if task ID exists
		for _, item := range db.Recs_task {
			if id == item.ID {
				exists = true
				// Task ID exists, check if task belongs to logged in user
				if userID == item.UserID {
					rec = item
					allowed = true
					break
				}
			}
		}
		db.Mu.Unlock()

		if exists {
			if allowed {
				if err := json.NewEncoder(w).Encode(rec); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					InfoLogger.Printf("Task retrieval sucessful (task ID %d) for user ID %d\n", id, userID)
				}
			} else {
				http.Error(w, "You are not authorized to view this record.", http.StatusInternalServerError)
				ErrorLogger.Println("Unauthorized access to record")
			}
		} else {
			fmt.Fprintln(w, "{\"success\": false}")
			http.Error(w, "Record not found", http.StatusNotFound)
			ErrorLogger.Println("Record not found")
		}
		return
	case "PUT":
		w.Header().Set("Content-Type", "application/json")
		// Check if record exists
		exists := false
		allowed := false
		var recId int
		// Loop through all records to check if task ID exists
		for j, item := range db.Recs_task {
			if id == item.ID {
				exists = true
				// Task ID exists, check if task belongs to logged in user
				if userID == item.UserID {
					recId = j
					allowed = true
					break
				}
			}
		}

		if exists {
			if allowed {
				var rec Task
				if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				db.Mu.Lock()
				db.Recs_task[recId].Name = rec.Name
				db.Recs_task[recId].Priority = rec.Priority
				db.Recs_task[recId].Status = rec.Status
				db.Mu.Unlock()
				fmt.Fprintln(w, "{\"success\": true}")
				InfoLogger.Printf("Task update sucessful (task ID %d) for user ID %d\n", id, userID)
			} else {
				http.Error(w, "You are not authorized to view this record.", http.StatusInternalServerError)
				ErrorLogger.Println("Unauthorized access to record")
			}

		} else {
			fmt.Fprintln(w, "{\"success\": false}")
			http.Error(w, "Record not found", http.StatusNotFound)
			ErrorLogger.Println("Record not found")
		}

		return
	}
}
