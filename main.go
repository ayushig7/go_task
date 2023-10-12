/*
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Tweet struct {
	Content string
}

func main() {
	usersTweets := map[int][]Tweet{
		1: {},
		2: {},
		3: {},
		4: {},
		5: {},
	}
	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		// Path /api/users/%d/tweets
		idAndTweets := r.URL.Path[len("/api/users/"):]
		segments := strings.Split(idAndTweets, "/")
		if len(segments) != 2 || segments[1] != "tweets" {
			w.WriteHeader(400)
			w.Write([]byte("Invalid path should be /api/users/%d/tweets"))
			return
		}
		userId, err := strconv.Atoi(segments[0])
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Invalid path should be /api/users/%d/tweets"))
			return
		}
		tweets, ok := usersTweets[userId]
		if !ok {
			w.WriteHeader(400)
			w.Write([]byte("User id is not valid"))
			return
		}
		switch r.Method {
		case "GET":
			jsonData, err := json.Marshal(tweets)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Unknown Server error"))
				return
			}
			w.WriteHeader(200)
			w.Write(jsonData)
			return
		case "POST":
			jsonInput, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Unknown Server error"))
				return
			}
			var tweet Tweet
			if err := json.Unmarshal(jsonInput, &tweet); err != nil {
				w.WriteHeader(400)
				w.Write([]byte("Invalid input"))
				return
			}
			if tweet.Content == "" {
				w.WriteHeader(400)
				w.Write([]byte("Content is required"))
				return
			}
			if len(tweet.Content) > 250 {
				w.WriteHeader(400)
				w.Write([]byte("Content must be less than 250"))
				return
			}
			usersTweets[userId] = append(tweets, tweet)
		default:
			w.WriteHeader(400)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
*/

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tweet struct {
	Content string
}

func addUserToDatabase(userID int, userName string, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO users (id, name) VALUES ($1, $2)", userID, userName)
	if err != nil {
		return err
	}
	return nil
}
func main() {
	var err error // declare err outside of the if statement
	db, err := sql.Open("pgx", "user=postgres password=ayushigupta7@ host=localhost port=5432 database=codesmith sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	//go get github.com/jackc/pgx/v4

	http.HandleFunc("/api/users/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		/*also we can use
		        idandname := r.URL.Path[len("/api/users/add"):]
				segments := strings.Split(idandname, "/")
		        userId, err := strconv.Atoi(segments[0])
		        if err != nil {
					w.WriteHeader(400)
					w.Write([]byte("Invalid path should be /api/users/add/%d/name"))
					return
				}

		*/

		var user User

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		if user.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		err := addUserToDatabase(user.ID, user.Name, db)
		if err != nil {
			http.Error(w, "Failed to add user to database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

	//http.HandleFunc("/api/users/get", getUserByIDHandler)

	http.HandleFunc("/api/users/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Extract user ID from the URL query parameters
		userIDStr := r.URL.Query().Get("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Query the database for the user with the specified ID
		row := db.QueryRow("SELECT id, name FROM users WHERE id = $1", userID)

		var id int
		var name string
		err = row.Scan(&id, &name)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, "Error querying database", http.StatusInternalServerError)
			}
			return
		}

		user := User{
			ID:   id,
			Name: name,
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	// http.HandleFunc("/api/users/getAll", getAllUsersHandler)

	http.HandleFunc("/api/users/getAll", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Extract pagination parameters
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		// Convert page and limit to integers
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			http.Error(w, "Invalid limit value", http.StatusBadRequest)
			return
		}

		// Calculate offset for pagination
		offset := (page - 1) * limit

		// Query the database for users with pagination
		rows, err := db.Query("SELECT id, name FROM users LIMIT $1 OFFSET $2", limit, offset)
		if err != nil {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Iterate through the rows and construct a slice of User objects
		var users []User
		for rows.Next() {
			var id int
			var name string
			err := rows.Scan(&id, &name)
			if err != nil {
				http.Error(w, "Error reading data from database", http.StatusInternalServerError)
				return
			}
			user := User{
				ID:   id,
				Name: name,
			}
			users = append(users, user)
		}

		// Check for errors from iterating over rows
		err = rows.Err()
		if err != nil {
			http.Error(w, "Error iterating through database results", http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	// http.HandleFunc("/api/users/tweets", getUserTweetsHandler)

	http.HandleFunc("/api/users/tweets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Extract user ID from the URL query parameters
		userIDStr := r.URL.Query().Get("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID < 1 {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Extract pagination parameters
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		// Convert page and limit to integers
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			http.Error(w, "Invalid limit value", http.StatusBadRequest)
			return
		}

		// Calculate offset for pagination
		offset := (page - 1) * limit

		// Query the database for tweets of the specified user with pagination
		rows, err := db.Query("SELECT content FROM tweets WHERE user_id = $1 LIMIT $2 OFFSET $3", userID, limit, offset)
		if err != nil {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Iterate through the rows and construct a slice of Tweet objects
		var tweets []Tweet
		for rows.Next() {
			var content string
			err := rows.Scan(&content)
			if err != nil {
				http.Error(w, "Error reading data from database", http.StatusInternalServerError)
				return
			}
			tweet := Tweet{
				Content: content,
			}
			tweets = append(tweets, tweet)
		}

		// Check for errors from iterating over rows
		err = rows.Err()
		if err != nil {
			http.Error(w, "Error iterating through database results", http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(tweets)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

}
