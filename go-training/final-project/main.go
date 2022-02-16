package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Initializing constant for server path
const serverURL string = "http://localhost:8080"

// Initializing class for a task record
type Task struct {
	ID       int
	UserID   int
	Name     string
	Priority string
	Status   string
}

// Initializing class for user authorization
type User struct {
	ID int
}

func main() {
	var activeUser, taskid int
	var exitApp bool = false
	var endSession bool = false
	var username, password, taskname, priority, status, confirmation string
	fmt.Println("Task Manager v1.0")

	// Start of application. It will only end when user select exit option.
	for !exitApp {
		// Start menu. No active session yet.
		var menuResponse int
		fmt.Println("-------------------")
		fmt.Println("What would you like to to do?")
		fmt.Println("1: Log in to an account")
		fmt.Println("2: Create an account")
		fmt.Println("3: Exit application")
		fmt.Print("Input option: ")
		fmt.Scan(&menuResponse)

		if menuResponse == 1 { // Login option
			fmt.Println("Please input your username and password:")
			fmt.Print("Username: ")
			fmt.Scan(&username)
			fmt.Print("Password: ")
			fmt.Scan(&password)
			activeUser = loginUser(username, password)
			if activeUser != -1 { // Login succesful
				fmt.Println("Log in successful! Welcome back.")
				endSession = false
				for !endSession {
					fmt.Println("-------------------")
					fmt.Println("What would you like to do next?")
					fmt.Println("1: View my tasks")
					fmt.Println("2: Create a task")
					fmt.Println("3: Update a task")
					fmt.Println("4: Delete a task")
					fmt.Println("5: Log out of account")
					fmt.Println("6: Exit application")
					fmt.Print("Input option: ")
					fmt.Scan(&menuResponse)
					if menuResponse == 1 { // Display all tasks of logged in user
						tasks := viewTasks(activeUser)
						for _, task := range tasks {
							fmt.Println("**********************")
							fmt.Printf("Task ID: %v\n", task.ID)
							fmt.Printf("Name: %v\n", task.Name)
							fmt.Printf("Priority: %v\n", task.Priority)
							fmt.Printf("Status: %v\n", task.Status)
						}
					} else if menuResponse == 2 { // Create a task for logged in user
						fmt.Println("-------------------")
						fmt.Println("Please input task details.")
						fmt.Print("Task Name: ")
						fmt.Scan(&taskname)
						fmt.Print("Priority: ")
						fmt.Scan(&priority)
						fmt.Print("Status: ")
						fmt.Scan(&status)
						if createTask(activeUser, taskname, priority, status) {
							fmt.Println("Task successfully created.")
						}
					} else if menuResponse == 3 { // Update a task for logged in user
						fmt.Println("-------------------")
						fmt.Println("Which task do you want to update?")
						fmt.Print("Task ID: ")
						fmt.Scan(&taskid)
						fmt.Println("Please input task details.")
						fmt.Print("Task Name: ")
						fmt.Scan(&taskname)
						fmt.Print("Priority: ")
						fmt.Scan(&priority)
						fmt.Print("Status: ")
						fmt.Scan(&status)
						if updateTask(activeUser, taskid, taskname, priority, status) {
							fmt.Println("Task successfully updated.")
						}
					} else if menuResponse == 4 { // Delete a task for logged in user
						fmt.Println("-------------------")
						fmt.Println("Which task do you want to delete?")
						fmt.Print("Task ID: ")
						fmt.Scan(&taskid)
						fmt.Println("Are you sure?")
						fmt.Print("Y/N: ")
						fmt.Scan(&confirmation)
						if confirmation == "Y" {
							if deleteTask(activeUser, taskid) {
								fmt.Println("Task successfully deleted.")
							}
						}
					} else if menuResponse == 5 { // Logout user
						endSession = true
						if logoutUser() { // Successful log out
							fmt.Println("Log out successful.")
						} else { // Unsuccessful log out
							fmt.Println("Error logging out")
						}

					} else if menuResponse == 6 { // Exit application for logged in user. Should automatically log out first.
						exitApp = true
						endSession = true
						if logoutUser() { // Successful log out
							fmt.Println("Log out successful.")
						} else { // Unsuccessful log out
							fmt.Println("Error logging out")
						}
					} else { // User input not among options
						fmt.Println("Invalid input. Please try again.")
					}
				}
			}
		} else if menuResponse == 2 { // User creation
			fmt.Println("-------------------")
			fmt.Println("Please input your desired username and password:")
			fmt.Print("Username: ")
			fmt.Scan(&username)
			fmt.Print("Password: ")
			fmt.Scan(&password)
			if createUser(username, password) { // User creation success
				fmt.Println("User succesfully created!")
			} else { // User creation failed
				fmt.Println("Error with user creation. Please try again.")
			}
		} else if menuResponse == 3 { // Exit application without logged in user
			exitApp = true
		} else { // User input not among options
			fmt.Println("Invalid input. Please try again.")
		}
	}

}

// Function for user login. Returns true if successful, false if not
func loginUser(username string, password string) int {
	var userid int
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	jsonRec := fmt.Sprintf("{\"Username\": \"%s\", \"Password\": \"%s\"}", username, password)
	outData := bytes.NewBuffer([]byte(jsonRec))
	res, err := c.Post(serverURL+"/login", "application/json", outData)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		fmt.Println("Invalid username and/or password. Please try again.")
		userid = -1
	} else {
		var u User
		err3 := json.Unmarshal(body, &u)
		if err3 != nil {
			fmt.Println("Invalid username and/or password. Please try again.")
			userid = -1
		} else {
			userid = u.ID
		}
	}

	return userid

}

// Function for user logout. Returns true if successful, false if not
func logoutUser() bool {
	var success bool
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	res, err := c.Post(serverURL+"/logout", "application/json", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode < 400 {
		success = true
	} else {
		success = false
		fmt.Println(res.StatusCode)
	}

	return success

}

// Function for user creation. Returns true if successful, false if not
func createUser(username string, password string) bool {
	var success bool
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	jsonRec := fmt.Sprintf("{\"Username\": \"%s\", \"Password\": \"%s\"}", username, password)
	outData := bytes.NewBuffer([]byte(jsonRec))
	res, err := c.Post(serverURL+"/signup", "application/json", outData)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode < 400 {
		success = true
	} else {
		success = false

	}

	return success
}

// Function for retrieving a user's tasks. Returns all tasks by user.
func viewTasks(userid int) []Task {
	var task_rec []Task
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	res, err := c.Get(serverURL + "/tasks?userid=" + strconv.Itoa(userid))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&task_rec); err != nil {
		log.Fatal(err)
	}
	return task_rec
}

// Function for task creation. Returns true if successful, false if not
func createTask(userid int, taskname string, priority string, status string) bool {
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	body := fmt.Sprintf("{\"Name\": \"%s\", \"Priority\": \"%s\", \"Status\": \"%s\"}", taskname, priority, status)
	payload := bytes.NewBuffer([]byte(body))

	res, err := c.Post(serverURL+"/tasks?userid="+strconv.Itoa(userid), "application/json", payload)
	if err != nil {
		fmt.Println("Task creation unsuccessful: ", err)
	}
	defer res.Body.Close()

	return res.StatusCode < 400
}

// Function for task update. Returns true if successful, false if not
func updateTask(userid int, taskid int, taskname string, priority string, status string) bool {
	var success bool
	client := http.Client{Timeout: time.Duration(1) * time.Second}
	body := fmt.Sprintf("{\"Name\": \"%s\", \"Priority\": \"%s\", \"Status\": \"%s\"}", taskname, priority, status)
	payload := bytes.NewBuffer([]byte(body))

	request, err := http.NewRequest("PUT", serverURL+"/tasks/"+fmt.Sprint(taskid)+"?userid="+strconv.Itoa(userid), payload)
	if err != nil {
		fmt.Println("Task update unsuccessful: ", err)
		success = false
	}
	res, err2 := client.Do(request)
	if err2 != nil {
		fmt.Println("Task update unsuccessful: ", err)
		success = false
	} else {
		success = true
	}
	defer res.Body.Close()

	return success
}

// Function for task deletion. Returns true if successful, false if not
func deleteTask(userid int, taskid int) bool {
	var success bool
	client := http.Client{Timeout: time.Duration(1) * time.Second}
	request, err := http.NewRequest("DELETE", serverURL+"/tasks/"+fmt.Sprint(taskid)+"?userid="+strconv.Itoa(userid), nil)
	if err != nil {
		fmt.Println("Task delete unsuccessful: ", err)
		success = false
	}

	res, err2 := client.Do(request)
	if err2 != nil {
		fmt.Println("Task delete unsuccessful: ", err)
		success = false
	} else {
		success = true
	}
	defer res.Body.Close()

	return success
}
