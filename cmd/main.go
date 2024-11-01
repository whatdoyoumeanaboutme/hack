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

type UserData struct {
	Username string
	Password string
}

// Функция для обработки запросов на регистрацию
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("register.html") // наша регистрация
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Проверяем, существует ли пользователь с таким именем
		// Мб довать бд
		if _, ok := users[username]; ok {
			http.Error(w, "Пользователь с таким именем уже существует", http.StatusBadRequest)
			return
		}

		// Создаем нового пользователя
		users[username] = UserData{Username: username, Password: password}

		// Перенаправление на страницу авторизации после успешной регистрации
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("login.html") // наша авторизация
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Выполняем простую проверку
		var errorString string
		if user, ok := users[username]; ok {
			if user.Password == password {
				http.Redirect(w, r, "/welcome", http.StatusFound)
				return
			} else {

				errorString = "Ошибка авторизации"
				http.Error(w, errorString, http.StatusUnauthorized)
				return
			}
		} else {
			// Пользователь не найден
			errorString = "Пользователь не найден"
			http.Error(w, errorString, http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

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

var users = make(map[string]UserData)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Изменение маршрута: теперь / будет перенаправлять на loginHandler
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/register", registerHandler)

	fmt.Println("Server starting on http://localhost:8080...")
	http.ListenAndServe(":8080", nil)
}
