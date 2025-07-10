package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"student_app/models"
	"student_app/utils"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/xuri/excelize/v2"
)

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request: Could not parse form", http.StatusBadRequest)
			log.Printf("Error parsing form: %v", err)
			return
		}

		name := r.FormValue("name")
		class := r.FormValue("class")
		ageStr := r.FormValue("age")

		age, err := strconv.Atoi(ageStr)
		if err != nil {
			http.Error(w, "Bad Request: Invalid age format", http.StatusBadRequest)
			log.Printf("Error converting age: %v", err)
			return
		}

		newStudent := models.Student{
			Name:  name,
			Class: class,
			Age:   age,
		}

		_, err = models.AddStudent(newStudent)
		if err != nil {
			http.Error(w, "Internal Server Error: Could not add student", http.StatusInternalServerError)
			log.Printf("Error adding student to DB: %v", err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	utils.RenderTemplate(w, "add_student.html", nil)
}

func EditStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request: Invalid student ID", http.StatusBadRequest)
		return
	}

	studentToEdit, err := models.GetStudentByID(id)
	if err != nil {
		http.Error(w, "Internal Server Error: Could not retrieve student", http.StatusInternalServerError)
		log.Printf("Error retrieving student from DB: %v", err)
		return
	}
	if studentToEdit == nil {
		http.Error(w, "Not Found: Student not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request: Could not parse form", http.StatusBadRequest)
			log.Printf("Error parsing form: %v", err)
			return
		}

		studentToEdit.Name = r.FormValue("name")
		studentToEdit.Class = r.FormValue("class")
		age, err := strconv.Atoi(r.FormValue("age"))
		if err != nil {
			http.Error(w, "Bad Request: Invalid age format", http.StatusBadRequest)
			log.Printf("Error converting age: %v", err)
			return
		}
		studentToEdit.Age = age

		err = models.UpdateStudent(*studentToEdit)
		if err != nil {
			http.Error(w, "Internal Server Error: Could not update student", http.StatusInternalServerError)
			log.Printf("Error updating student in DB: %v", err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	utils.RenderTemplate(w, "edit_student.html", studentToEdit)
}

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request: Invalid student ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteStudent(id)
	if err != nil {
		http.Error(w, "Internal Server Error: Could not delete student", http.StatusInternalServerError)
		log.Printf("Error deleting student from DB: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func PdfReportHandler(w http.ResponseWriter, r *http.Request) {
	students, err := models.GetAllStudents()
	if err != nil {
		http.Error(w, "Internal Server Error: Could not retrieve students for PDF report", http.StatusInternalServerError)
		log.Printf("Error retrieving students for PDF: %v", err)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Laporan Data Siswa")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(20, 7, "ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 7, "Nama", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 7, "Kelas", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 7, "Umur", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	for _, s := range students {
		pdf.CellFormat(20, 7, strconv.Itoa(s.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 7, s.Name, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 7, s.Class, "1", 0, "", false, 0, "")
		pdf.CellFormat(20, 7, strconv.Itoa(s.Age), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=laporan_siswa.pdf")

	err = pdf.Output(w)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		http.Error(w, "Internal Server Error: Could not generate PDF", http.StatusInternalServerError)
	}
}

func ExcelReportHandler(w http.ResponseWriter, r *http.Request) {
	students, err := models.GetAllStudents()
	if err != nil {
		http.Error(w, "Internal Server Error: Could not retrieve students for Excel report", http.StatusInternalServerError)
		log.Printf("Error retrieving students for Excel: %v", err)
		return
	}

	f := excelize.NewFile()
	sheetName := "Data Siswa"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Printf("Error creating new Excel sheet: %v", err)
		http.Error(w, "Internal Server Error: Could not create Excel sheet", http.StatusInternalServerError)
		return
	}
	f.SetActiveSheet(index)

	// Set header
	f.SetCellValue(sheetName, "A1", "ID")
	f.SetCellValue(sheetName, "B1", "Nama")
	f.SetCellValue(sheetName, "C1", "Kelas")
	f.SetCellValue(sheetName, "D1", "Umur")

	// Populate data
	for i, s := range students {
		rowNum := i + 2 // Data dimulai dari baris kedua (setelah header)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), s.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), s.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), s.Class)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), s.Age)
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=laporan_siswa.xlsx")

	if err := f.Write(w); err != nil {
		log.Printf("Error writing Excel file: %v", err)
		http.Error(w, "Internal Server Error: Could not write Excel file", http.StatusInternalServerError)
	}
}
