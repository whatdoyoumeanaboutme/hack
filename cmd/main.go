package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

type Credentials struct {
	Username string
	Password string
}

type WelcomeData struct {
	Message string
}

type UserData struct {
	Username string
	Password string
}

func connectToDB() (*sql.DB, error) {
    db, err := sql.Open("postgres", "user=postgres password=123 host=localhost port=5432 dbname=postgres sslmode=disable")
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(`
        DO $$ 
        BEGIN
            IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'userauth') THEN
                CREATE DATABASE userauth;
            END IF;
        END $$;
    `)
    if err != nil {
        return nil, err
    }
    
    db.Close()

    db, err = sql.Open("postgres", "user=postgres password=123 host=localhost port=5432 dbname=userauth sslmode=disable")
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            username VARCHAR(50) PRIMARY KEY,
            password_hash VARCHAR(100) NOT NULL
        )
    `)
    if err != nil {
        return nil, err
    }

    return db, nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        db, err := connectToDB()
        if err != nil {
            http.Error(w, "Database connection error", http.StatusInternalServerError)
            return
        }
        defer db.Close()

        username := r.FormValue("username")
        password := r.FormValue("password")

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(w, "Password hashing error", http.StatusInternalServerError)
            return
        }

        var exists bool
        err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
        if err != nil {
            http.Error(w, "User check error", http.StatusInternalServerError)
            return
        }
        if exists {
            http.Error(w, "Username already exists", http.StatusBadRequest)
            return
        }

        _, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, string(hashedPassword))
        if err != nil {
            http.Error(w, "Registration error", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/", http.StatusFound)
    } else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        db, err := connectToDB()
        if err != nil {
            http.Error(w, "Database connection error", http.StatusInternalServerError)
            return
        }
        defer db.Close()

        username := r.FormValue("username")
        password := r.FormValue("password")

        var hashedPassword string
        err = db.QueryRow("SELECT password_hash FROM users WHERE username = $1", username).Scan(&hashedPassword)
        if err != nil {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
        if err == nil {
            http.Redirect(w, r, "/welcome", http.StatusFound)
            return
        } else {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }
    } else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	data := WelcomeData{Message: "Добро пожаловать"}
	tmpl, err := template.ParseFiles("welcome.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func main() {
	cssFiles := http.FileServer(http.Dir("css"))
	jsFiles := http.FileServer(http.Dir("js"))
	
	http.Handle("/css/", http.StripPrefix("/css/", cssFiles))
	http.Handle("/js/", http.StripPrefix("/js/", jsFiles))

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)

	fmt.Println("Server starting on http://localhost:8080...")
	http.ListenAndServe(":8080", nil)
}
