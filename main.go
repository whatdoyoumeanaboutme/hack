package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Credentials struct {
	Username string
	Password string
}

type WelcomeData struct {
	Message string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Тут мы к файлу идём
		tmpl, err := template.ParseFiles("login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// это наш логин и пароль
		if username == "user" && password == "pass" {
			// типо если верно велком , а если нет в попу
			http.Redirect(w, r, "/welcome", http.StatusFound)
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	// Отображаем страницу приветствия с использованием шаблона
	data := WelcomeData{Message: "Welcome to the website!"}
	tmpl, err := template.ParseFiles("welcome.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/welcome", welcomeHandler)

	fmt.Println("Server starting on http://localhost:8080...")
	http.ListenAndServe(":8080", nil)
}
