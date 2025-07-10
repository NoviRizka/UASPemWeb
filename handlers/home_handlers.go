package handlers

import (
	"net/http"
	"student_app/models"
	"student_app/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	students, err := models.GetAllStudents()
	if err != nil {
		http.Error(w, "Internal Server Error: Could not retrieve students", http.StatusInternalServerError)
		return
	}
	utils.RenderTemplate(w, "index.html", students)
}
