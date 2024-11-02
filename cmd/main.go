package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

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

// Функция для подключения к базе данных PostgreSQL
func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres password=123 host=localhost port=5432 dbname=auto sslmode=disable")
	if err != nil {
		return nil, err
	}
	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Функция для регистрации нового пользователя
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		db, err := connectToDB()
		if err != nil {
			http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Проверяем, существует ли пользователь с таким именем
		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
		if err != nil {
			http.Error(w, "Ошибка проверки существования пользователя", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Пользователь с таким именем уже существует", http.StatusBadRequest)
			return
		}

		// Вставляем нового пользователя в базу данных
		_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
		if err != nil {
			http.Error(w, "Ошибка регистрации пользователя", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Функция для авторизации пользователя
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		db, err := connectToDB()
		if err != nil {
			http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Проверяем, существует ли пользователь с таким именем
		var storedPassword string
		err = db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&storedPassword)
		if err != nil {
			// Отображаем сообщение об ошибке в том же шаблоне
			tmpl, err := template.ParseFiles("login.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data := map[string]interface{}{
				"ErrorMessage": "Ошибка проверки пользователя",
			}
			tmpl.Execute(w, data)
			return
		}

		// Сравниваем введённый пароль с паролем из базы данных
		if storedPassword == password {
			http.Redirect(w, r, "/welcome", http.StatusFound)
			return
		} else {
			// Отображаем сообщение об ошибке в том же шаблоне
			tmpl, err := template.ParseFiles("login.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data := map[string]interface{}{
				"ErrorMessage": "Ошибка авторизации",
			}
			tmpl.Execute(w, data)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Функция для отображения страницы приветствия
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	// Отображаем страницу приветствия с использованием шаблона
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
