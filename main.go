package main

import (
	"fmt"
	"log"
	"net/http"
	"student_app/handlers"
	"student_app/models"
)

func main() {
	models.InitDB()
	defer models.CloseDB()

	go models.CleanUpExpiredSessions()

	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	http.HandleFunc("/", handlers.AuthMiddleware(handlers.HomeHandler))                    // Rute utama: menampilkan daftar siswa
	http.HandleFunc("/add", handlers.AuthMiddleware(handlers.AddStudentHandler))           // Rute untuk menambah siswa baru (formulir dan proses)
	http.HandleFunc("/edit", handlers.AuthMiddleware(handlers.EditStudentHandler))         // Rute untuk mengedit data siswa (formulir dan proses)
	http.HandleFunc("/delete", handlers.AuthMiddleware(handlers.DeleteStudentHandler))     // Rute untuk menghapus siswa
	http.HandleFunc("/report/pdf", handlers.AuthMiddleware(handlers.PdfReportHandler))     // Rute untuk menghasilkan laporan PDF
	http.HandleFunc("/report/excel", handlers.AuthMiddleware(handlers.ExcelReportHandler)) // Rute untuk menghasilkan laporan Excel

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server berjalan di http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
