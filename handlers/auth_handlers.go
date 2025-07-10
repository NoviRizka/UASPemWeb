package handlers

import (
	"log"
	"net/http"
	"student_app/models"
	"student_app/utils"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request: Could not parse form", http.StatusBadRequest)
			log.Printf("Error parsing login form: %v", err)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if models.IsAuthenticated(username, password) {
			sessionToken, err := models.CreateSession(username)
			if err != nil {
				http.Error(w, "Internal Server Error: Could not create session", http.StatusInternalServerError)
				log.Printf("Error creating session: %v", err)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    sessionToken,
				Expires:  time.Now().Add(1 * time.Hour),
				HttpOnly: true,
				Path:     "/",
			})

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		data := map[string]string{"Error": "Username atau password salah!"}
		utils.RenderTemplate(w, "login.html", data)
		return
	}
	utils.RenderTemplate(w, "login.html", nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Dapatkan cookie sesi.
	c, err := r.Cookie("session_token")
	if err == nil {
		models.DeleteSession(c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, ok := models.GetSession(c.Value)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
