package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// @title       My Awesome API
// @version     1.0
// @description An API to demonstrate OpenAPI integration.
// @termsOfService  http://swagger.io/terms/
// @contact.name  API Support
// @contact.email support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
// @schemes   http

type Credentials struct {
	Username string
	Password string
}

type WelcomeData struct {
	Message string
}

// @Summary Login to the system
// @Description Logs in a user with given username and password.
// @Accept  json
// @Produce json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {object} Credentials
// @Failure 401 {object} error "Invalid credentials"
// @Router /login [post]

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

// @Summary Get welcome page
// @Description Returns the welcome page with a message.
// @Accept  json
// @Produce  json
// @Success 200 {object} WelcomeData
// @Router /welcome [get]

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
