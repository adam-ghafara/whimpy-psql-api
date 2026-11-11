package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type hrUser struct {
	Username    string `json:"hrUsername"`
	FullName    string `json:"hrFullName"`
	Email       string `json:"hrEmail"`
	Password    string `json:"hrPassword"`
	PhoneNumber string `json:"hrPhone"`
	AutoRole    int    `json:"userRole"`
}

type internUser struct {
	Username    string `json:"internUsername"`
	FullName    string `json:"internFullName"`
	Email       string `json:"internEmail"`
	Password    string `json:"internPassword"`
	PhoneNumber string `json:"internPhone"`
	AutoRole    int    `json:"userRole"`
}

type mentorUser struct {
	Username    string `json:"mentorUsername"`
	FullName    string `json:"mentorFullName"`
	Email       string `json:"mentorEmail"`
	Password    string `json:"mentorPassword"`
	PhoneNumber string `json:"mentorPhone"`
	AutoRole    int    `json:"userRole"`
}

var db *sql.DB

func main() {
	// NOTE: This is the PUBLIC version of the code.
	// Never commit real database credentials here.
	// Set DATABASE_URL in your environment (or use a .env file that is git-ignored),
	// e.g.:
	//   DATABASE_URL="postgresql://example_user:example_password@ep-example-host-123456.ap-southeast-1.aws.neon.tech/example_db?sslmode=require"
	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		// Dummy/example connection for local testing only.
		DATABASE_URL = "postgresql://example_user:example_password@localhost:5432/example_db?sslmode=disable"
	}
	// Example pq.Config structure (kept for reference):
	// cfg := pq.Config{
	// 	Host:           "ep-example-host-123456.ap-southeast-1.aws.neon.tech",
	// 	Port:           5432,
	// 	User:           "example_user",
	// 	Password:       os.Getenv("DATABASE_PASSWORD"),
	// 	Database:       "example_db",
	// 	SSLMode:        "require",
	// 	SSLSNI:         true,
	// 	ConnectTimeout: 5 * time.Second,
	// }

	cfg, err := pq.NewConfig(DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}

	c, err := pq.NewConnectorConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db = sql.OpenDB(c)
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	fmt.Println("Initializing...")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		fmt.Println("▐██████▌ ▐█▌▐███████▐█▌ ██████ ▐████████ ██▌ ████████▄  ███████▌██             ████████▄  ▄████████▄▐█████████▌  ███████   ")
		fmt.Println("▐██████▌ ▐██▐███████▐██ ▀▀▀▀▀▀ ▐███████▀█▀██▐███████▌█▌▐███████▌██            ▐███████▌█▌▐████▐█████▐███████▌██ ▐███████   ")
		fmt.Println("▐██████▌ ▐██▐███████▐██ ███████▐███████   ██▐███████▌█▌▐███████▌██            ▐███████▌█▌▐████      ▐███████▌██ ▐███████   ")
		fmt.Println("▐██████▌ ▐██▐███████▐██▐███████▐███████   ██▐███████▌█▌▐███████▌██            ▐███████▌█▌▐████████▄ ▐███████▌██ ▐███████▐██")
		fmt.Println("▐██████▌ ▐██▐██████████▐███████▐███████   ██▐████████▀ ▐██████████ ███████████▐████████▀  ▀▀▀▀▀█████▐███████▌██ ▐███████▐██")
		fmt.Println("▐██████▌█▐██▐███████▐██▐███████▐███████   ██▐███████▌   █████████▌            ▐███████▌  ▐████▌█████▐███████▌██ ▐███████▐██")
		fmt.Println("▐███████████▐███████▐██▐███████▐███████   ██▐███████▌   ▐███████              ▐███████▌   ▀███████▀ ▐██████████▌▐██████████")
		fmt.Println("Whimpy API Server (PostgreSQL Database) | Created by: github.com/adam-ghafara                                ▀█       		")
		fmt.Println("See the Documentation guide to know more about the API functions.")
		fmt.Println("Use Postman (Recommended) to test your API.")
	}).Methods("GET")
	r.HandleFunc("/hr/add", AddHRUserHandler).Methods("POST")
	r.HandleFunc("/mentor/add", AddMentorUserHandler).Methods("POST")
	r.HandleFunc("/intern/add", AddInternUserHandler).Methods("POST")
	r.HandleFunc("/hr/get_all", GetAllHRUser).Methods("GET")
	r.HandleFunc("/mentor/get_all", GetAllMentorUser).Methods("GET")
	r.HandleFunc("/intern/get_all", GetAllInternUser).Methods("GET")
	r.HandleFunc("/hr/update/{id}", UpdateHRUserInfo).Methods("PUT")
	r.HandleFunc("/mentor/update/{id}", UpdateMentorUserInfo).Methods("PUT")
	r.HandleFunc("/intern/update/{id}", UpdateInternUserInfo).Methods("PUT")
	r.HandleFunc("/hr/set_new_pass/{id}", SetHRNewPass).Methods("PUT")
	r.HandleFunc("/mentor/set_new_pass/{id}", SetMentorNewPass).Methods("PUT")
	r.HandleFunc("/intern/set_new_pass/{id}", SetInternNewPass).Methods("PUT")
	r.HandleFunc("/hr/delete/{id}", DeleteHRUserHandler).Methods("DELETE")
	r.HandleFunc("/mentor/delete/{id}", DeleteMentorUserHandler).Methods("DELETE")
	r.HandleFunc("/intern/delete/{id}", DeleteInternUserHandler).Methods("DELETE")
	//
	fmt.Println("▐██████▌ ▐█▌▐███████▐█▌ ██████ ▐████████ ██▌ ████████▄  ███████▌██             ████████▄  ▄████████▄▐█████████▌  ███████   ")
	fmt.Println("▐██████▌ ▐██▐███████▐██ ▀▀▀▀▀▀ ▐███████▀█▀██▐███████▌█▌▐███████▌██            ▐███████▌█▌▐████▐█████▐███████▌██ ▐███████   ")
	fmt.Println("▐██████▌ ▐██▐███████▐██ ███████▐███████   ██▐███████▌█▌▐███████▌██            ▐███████▌█▌▐████      ▐███████▌██ ▐███████   ")
	fmt.Println("▐██████▌ ▐██▐███████▐██▐███████▐███████   ██▐███████▌█▌▐███████▌██            ▐███████▌█▌▐████████▄ ▐███████▌██ ▐███████▐██")
	fmt.Println("▐██████▌ ▐██▐██████████▐███████▐███████   ██▐████████▀ ▐██████████ ███████████▐████████▀  ▀▀▀▀▀█████▐███████▌██ ▐███████▐██")
	fmt.Println("▐██████▌█▐██▐███████▐██▐███████▐███████   ██▐███████▌   █████████▌            ▐███████▌  ▐████▌█████▐███████▌██ ▐███████▐██")
	fmt.Println("▐███████████▐███████▐██▐███████▐███████   ██▐███████▌   ▐███████              ▐███████▌   ▀███████▀ ▐██████████▌▐██████████")
	fmt.Println("Whimpy API Server (PostgreSQL Database) | Created by: github.com/adam-ghafara                                ▀█       		")
	fmt.Println("See the Documentation guide to know more about the API functions.")
	fmt.Println("Use Postman (Recommended) to test your API.")
	if err := http.ListenAndServe(":5724", r); err != nil {
		log.Fatal(err)
		return
	}
}

// Account Creation

func PasswordHashing(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func AddHRUserHandler(w http.ResponseWriter, r *http.Request) {
	var u hrUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if u.Username == "" || u.Email == "" || len(u.Password) < 8 || len(u.Password) > 72 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username, email required; password must be 8-72 chars"})
		return
	}

	hash, err := PasswordHashing(u.Password)
	if err != nil {
		log.Println("hash error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return // this return was missing in your code
	}
	u.AutoRole = 2
	_, err = db.Exec(
		`INSERT INTO wp_hr_list ("hrUsername", "hrFullname", "hrEmail", "hrPassword", "hrPhone", "roleID")
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		u.Username, u.FullName, u.Email, hash, u.PhoneNumber, u.AutoRole,
	)
	if err != nil {
		log.Println("insert error:", err) // shows the real cause in your terminal

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			log.Println("duplicate on constraint:", pqErr.Constraint)
			writeJSON(w, http.StatusConflict, map[string]string{"error": "username or email already exists"})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message":  "HR user created",
		"username": u.Username,
		"fullname": u.FullName,
		"email":    u.Email,
		"phone":    u.PhoneNumber,
	})
}

func AddMentorUserHandler(w http.ResponseWriter, r *http.Request) {
	var u mentorUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if u.Username == "" || u.Email == "" || len(u.Password) < 8 || len(u.Password) > 72 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username, email required; password must be 8-72 chars"})
		return
	}

	hash, err := PasswordHashing(u.Password)
	if err != nil {
		log.Println("hash error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return // this return was missing in your code
	}
	u.AutoRole = 3
	_, err = db.Exec(
		`INSERT INTO wp_mentor_list ("mtUsername", "mtFullname", "mtEmail", "mtPassword", "mtPhone", "roleID")
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		u.Username, u.FullName, u.Email, hash, u.PhoneNumber, u.AutoRole,
	)
	if err != nil {
		log.Println("insert error:", err) // shows the real cause in your terminal

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			log.Println("duplicate on constraint:", pqErr.Constraint)
			writeJSON(w, http.StatusConflict, map[string]string{"error": "username or email already exists"})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message":  "Mentor user created",
		"username": u.Username,
		"fullname": u.FullName,
		"email":    u.Email,
		"phone":    u.PhoneNumber,
	})
}

func AddInternUserHandler(w http.ResponseWriter, r *http.Request) {
	var u internUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if u.Username == "" || u.Email == "" || len(u.Password) < 8 || len(u.Password) > 72 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username, email required; password must be 8-72 chars"})
		return
	}

	hash, err := PasswordHashing(u.Password)
	if err != nil {
		log.Println("hash error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return // this return was missing in your code
	}
	u.AutoRole = 4
	_, err = db.Exec(
		`INSERT INTO wp_intern_list ("internUsername", "internFullname", "internEmail", "internPassword", "internPhone", "roleID")
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		u.Username, u.FullName, u.Email, hash, u.PhoneNumber, u.AutoRole,
	)
	if err != nil {
		log.Println("insert error:", err) // shows the real cause in your terminal

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			log.Println("duplicate on constraint:", pqErr.Constraint)
			writeJSON(w, http.StatusConflict, map[string]string{"error": "username or email already exists"})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message":  "Internship Participant user created",
		"username": u.Username,
		"fullname": u.FullName,
		"email":    u.Email,
		"phone":    u.PhoneNumber,
	})
}

// GET ALL USERS

type hrInfoList struct {
	ID          int    `json:"hrID"`
	Username    string `json:"hrUsername"`
	FullName    string `json:"hrFullName"`
	Email       string `json:"hrEmail"`
	Password    string `json:"hrPassword"`
	PhoneNumber string `json:"hrPhone"`
	UserRole    string `json:"userRole"`
}

type mentorInfoList struct {
	ID          int    `json:"mentorID"`
	Username    string `json:"mentorUsername"`
	FullName    string `json:"mentorFullName"`
	Email       string `json:"mentorEmail"`
	Password    string `json:"mentorPassword"`
	PhoneNumber string `json:"mentorPhone"`
	UserRole    string `json:"userRole"`
}

type InternInfoList struct {
	ID          int    `json:"internID"`
	Username    string `json:"internUsername"`
	FullName    string `json:"internFullName"`
	Email       string `json:"internEmail"`
	Password    string `json:"internPassword"`
	PhoneNumber string `json:"internPhone"`
	UserRole    string `json:"userRole"`
}

var roleName = map[int]string{
	1: "Admin",       // Access : Administration Content, Role Assignment
	2: "HR",          // Access : User Management, Mentor Assignment
	3: "Mentor",      // Access : Participant Control, Task Assignment
	4: "Participant", // Access : Task Management, Propose Assignment
}

func getRoleName(roleID int) string {
	if name, ok := roleName[roleID]; ok {
		return name
	}
	return "Unknown"
}

func GetAllHRUser(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * from wp_hr_list")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	userlist := []hrInfoList{}
	for rows.Next() {
		var u hrInfoList
		var roleID int // scan the number into a temporary variable
		if err := rows.Scan(&u.ID, &u.Username, &u.FullName, &u.Email, &u.Password, &u.PhoneNumber, &roleID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.UserRole = getRoleName(roleID) // convert to name
		userlist = append(userlist, u)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(userlist)
}

func GetAllMentorUser(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * from wp_mentor_list")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	userlist := []mentorInfoList{}
	for rows.Next() {
		var u mentorInfoList
		var roleID int // scan the number into a temporary variable
		if err := rows.Scan(&u.ID, &u.Username, &u.FullName, &u.Email, &u.Password, &u.PhoneNumber, &roleID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.UserRole = getRoleName(roleID) // convert to name
		userlist = append(userlist, u)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(userlist)
}

func GetAllInternUser(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * from wp_intern_list")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	userlist := []InternInfoList{}
	for rows.Next() {
		var u InternInfoList
		var roleID int // scan the number into a temporary variable
		if err := rows.Scan(&u.ID, &u.Username, &u.FullName, &u.Email, &u.Password, &u.PhoneNumber, &roleID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.UserRole = getRoleName(roleID) // convert to name
		userlist = append(userlist, u)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(userlist)
}

// SignIn Operator, Password Hashing is being used for SignIn Credential

// This struct can be used for any role signin.
type signinRequest struct {
	CredentialString string `json:"userCredential"`
	Password         string `json:"userPassword"`
}

// func passwordCheck(hash, password string) bool {
// 	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
// }

func SessionIDGenerator() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", b)
}

const sessionCookieName = "sessionID"

func HRsignInOperator(w http.ResponseWriter, r *http.Request) {
	var c signinRequest
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if c.CredentialString == "" || c.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Username/Email and Password are required"})
		return
	}

	// Hash and match the password
	var hash string
	var userID int
	err := db.QueryRow(`SELECT "hrID", "hrPassword" FROM wp_hr_list WHERE "hrUsername" = $1:;text or LOWER("hrEmail") = LOWER($1:;text) LIMIT 1`, c.CredentialString).Scan(&userID, &hash)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("login query error:", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
			return
		}
		// no such user: same response as a wrong password
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(c.Password)) != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "User Credential Accepted, continue to the signin process."})
	// Create SessionID and store it to Cookie
	sessionID := SessionIDGenerator()

	if err != nil {
		http.Error(w, "Serve Error..", http.StatusInternalServerError)
	}

	expires := time.Now().Add(720 * time.Hour)

	_, err = db.Exec("INSERT INTO wp_user_session (sessionid, userID, senssionExpire) (VALUES $1, $2, $3)",
		sessionID, userID, expires)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,                 // JS can't read it (blocks XSS theft)
		Secure:   false,                // set to true in production (HTTPS)
		SameSite: http.SameSiteLaxMode, // basic CSRF protection
	})
	writeJSON(w, http.StatusOK, map[string]any{
		"message":   "SessionID created. Info : ",
		"userID":    userID,
		"sessionID": sessionID})
}

func MentorsignInOperator(w http.ResponseWriter, r *http.Request) {
	var c signinRequest
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if c.CredentialString == "" || c.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Username/Email and Password are required"})
		return
	}

	// Hash and match the password
	var hash string
	err := db.QueryRow(`SELECT "mtPassword" FROM wp_mentor_list WHERE "mtUsername" = $1:;text or LOWER("mtEmail") = LOWER($1:;text) LIMIT 1`, c.CredentialString).Scan(&hash)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("login query error:", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
			return
		}
		// no such user: same response as a wrong password
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(c.Password)) != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "User Credential Accepted, continue to the signin process."})
	// Create SessionID and store it in Cookie

}

func InternsignInOperator(w http.ResponseWriter, r *http.Request) {
	var c signinRequest
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if c.CredentialString == "" || c.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Username/Email and Password are required"})
		return
	}

	// Hash and match the password
	var hash string
	err := db.QueryRow(`SELECT "internPassword" FROM wp_intern_list WHERE "internUsername" = $1:;text or LOWER("internEmail") = LOWER($1:;text) LIMIT 1`, c.CredentialString).Scan(&hash)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("login query error:", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
			return
		}
		// no such user: same response as a wrong password
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(c.Password)) != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "User Credential Accepted, continue to the signin process."})
}

// User info Update

func UpdateHRUserInfo(w http.ResponseWriter, r *http.Request) {
	getID := mux.Vars(r)["id"]
	var u hrUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	upt, err := db.Exec("UPDATE wp_hr_list SET hrUsername = $1, hrFullname = $2, hrEmail = $3, hrPhone = $4 WHERE hrID = $5", u.Username, u.FullName, u.Email, u.PhoneNumber, getID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := upt.RowsAffected(); n == 0 {
		http.Error(w, "User Not Found (Or no changes affected)", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "User Info has been Updated."})
}

func SetHRNewPass(w http.ResponseWriter, r *http.Request) {
	getID := mux.Vars(r)["id"]
	var u hrUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	hash, err := PasswordHashing(u.Password)
	if err != nil {
		log.Println("hash error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return // this return was missing in your code
	}
	upt, err := db.Exec("UPDATE wp_hr_list SET hrPassword = $1 WHERE hrID = $2", hash, getID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := upt.RowsAffected(); n == 0 {
		http.Error(w, "User Not Found (Or no changes affected)", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Password has been Changed."})
}

func UpdateMentorUserInfo(w http.ResponseWriter, r *http.Request) {
	getID := mux.Vars(r)["id"]
	var u mentorUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	upt, err := db.Exec("UPDATE wp_mentor_list SET mentorUsername = $1, mentorFullname = $2, mentorEmail = $3, mentorPhone = $4 WHERE mentorID = $5", u.Username, u.FullName, u.Email, u.PhoneNumber, getID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := upt.RowsAffected(); n == 0 {
		http.Error(w, "User Not Found (Or no changes affected)", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "User Info has been Updated."})
}

func SetMentorNewPass(w http.ResponseWriter, r *http.Request) {
	getID := mux.Vars(r)["id"]
	var u mentorUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	hash, err := PasswordHashing(u.Password)
	if err != nil {
		log.Println("hash error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return // this return was missing in your code
	}
	upt, err := db.Exec("UPDATE wp_mentor_list SET mentorPassword = $1 WHERE mentorID = $2", hash, getID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := upt.RowsAffected(); n == 0 {
		http.Error(w, "User Not Found (Or no changes affected)", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Password has been Changed."})
}

func UpdateInternUserInfo(w http.ResponseWriter, r *http.Request) {
	getID := mux.Vars(r)["id"]
	var u internUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	upt, err := db.Exec("UPDATE wp_intern_list SET internUsername = $1, internFullname = $2, internEmail = $3, internPhone = $4 WHERE internID = $5", u.Username, u.FullName, u.Email, u.PhoneNumber, getID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := upt.RowsAffected(); n == 0 {
		http.Error(w, "User Not Found (Or no changes affected)", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "User Info has been Updated."})
}

func SetInternNewPass(w http.ResponseWriter, r *http.Request) {
	getID := mux.Vars(r)["id"]
	var u internUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	hash, err := PasswordHashing(u.Password)
	if err != nil {
		log.Println("hash error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return // this return was missing in your code
	}
	upt, err := db.Exec("UPDATE wp_intern_list SET internPassword = $1 WHERE internID = $2", hash, getID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := upt.RowsAffected(); n == 0 {
		http.Error(w, "User Not Found (Or no changes affected)", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Password has been Changed."})
}

// User Delete (by ID)

func DeleteHRUserHandler(w http.ResponseWriter, r *http.Request) {
	getID := mux.Vars(r)["id"]
	res, err := db.Exec(`DELETE FROM wp_hr_list WHERE "hrID" = $1`, getID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "HR user has been Deleted.", "userID": getID})
}

func DeleteMentorUserHandler(w http.ResponseWriter, r *http.Request) {
	getID := mux.Vars(r)["id"]
	res, err := db.Exec(`DELETE FROM wp_mentor_list WHERE "mtID" = $1`, getID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Mentor user has been Deleted.", "userID": getID})
}

func DeleteInternUserHandler(w http.ResponseWriter, r *http.Request) {
	getID := mux.Vars(r)["id"]
	res, err := db.Exec(`DELETE FROM wp_intern_list WHERE "internID" = $1`, getID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Intern user has been Deleted.", "userID": getID})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
