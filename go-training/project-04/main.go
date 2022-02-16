package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Contact struct {
	ID       int
	Last     string
	First    string
	Company  string
	Address  string
	Country  string
	Position string
}

type Database struct {
	nextID int
	mu     sync.Mutex
	recs   []Contact
}

func main() {
	db := &Database{recs: []Contact{}}
	http.ListenAndServe(":8081", db.handler())
}

func (db *Database) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int
		if r.URL.Path == "/contacts" {
			db.process(w, r)
		} else if n, _ := fmt.Sscanf(r.URL.Path, "/contacts/%d", &id); n == 1 {
			db.processID(id, w, r)
		}
	}
}

func (db *Database) process(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var rec Contact
		if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		exists := false
		// Access database
		db.mu.Lock()

		// Check if record exists
		for _, item := range db.recs {
			if (rec.First == item.First) && (rec.Last == item.Last) && (rec.Company == item.Company) && (rec.Address == item.Address) && (rec.Country == item.Country) && (rec.Position == item.Position) {
				exists = true
			}
		}
		if exists { // Record exists, show error
			fmt.Fprintln(w, "Error: record already exists")
			fmt.Fprintln(w, "{\"success\": false}")
			http.Error(w, "Record already exists", http.StatusNotFound)
		} else { // Record does not exist
			rec.ID = db.nextID
			db.nextID++
			db.recs = append(db.recs, rec)
			fmt.Fprintln(w, "{\"success\": true}")
		}

		db.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")

	case "GET":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(db.recs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "PUT":
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return

	case "DELETE":
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
func (db *Database) processID(id int, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE": // Delete contact record
		exists := false
		db.mu.Lock()
		for j, item := range db.recs {
			if id == item.ID {
				db.recs = append(db.recs[:j], db.recs[j+1:]...)
				exists = true
				break
			}
		}
		db.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if exists {
			fmt.Fprintln(w, "{\"success\": true}")
		} else {
			fmt.Fprintln(w, "{\"success\": false}")
			http.Error(w, "Record not found", http.StatusNotFound)
		}
	case "POST": // Not allowed
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	case "GET": // Get contact record
		exists := false
		rec := Contact{}
		db.mu.Lock()
		for _, item := range db.recs {
			if id == item.ID {
				exists = true
				rec = item
				break
			}
		}
		db.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if exists {
			if err := json.NewEncoder(w).Encode(rec); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			fmt.Fprintln(w, "{\"success\": false}")
			http.Error(w, "Record not found", http.StatusNotFound)
		}
		return
	case "PUT": // Update contact record
		// Check if record exists
		exists := false
		var recId int
		for j, item := range db.recs {
			if id == item.ID {
				exists = true
				recId = j
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if exists {
			// Read json input
			var rec Contact
			if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			db.mu.Lock()
			db.recs[recId].Last = rec.Last
			db.recs[recId].First = rec.First
			db.recs[recId].Company = rec.Company
			db.recs[recId].Address = rec.Address
			db.recs[recId].Country = rec.Country
			db.recs[recId].Position = rec.Position
			db.mu.Unlock()
			fmt.Fprintln(w, "{\"success\": true}")
		} else {
			fmt.Fprintln(w, "{\"success\": false}")
			http.Error(w, "Record not found", http.StatusNotFound)
		}

		return
	}
}
