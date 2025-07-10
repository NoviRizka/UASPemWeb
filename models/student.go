package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Student struct {
	ID    int
	Name  string
	Class string
	Age   int
}

var db *sql.DB

func InitDB() {
	var err error
	dataSourceName := "root:@tcp(127.0.0.1:3306)/student_db"
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Koneksi database berhasil!")
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func GetAllStudents() ([]Student, error) {
	rows, err := db.Query("SELECT id, name, class, age FROM students ORDER BY id DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to query students: %v", err)
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Class, &s.Age); err != nil {
			return nil, fmt.Errorf("failed to scan student row: %v", err)
		}
		students = append(students, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	return students, nil
}

func GetStudentByID(id int) (*Student, error) {
	var s Student
	row := db.QueryRow("SELECT id, name, class, age FROM students WHERE id = ?", id)
	err := row.Scan(&s.ID, &s.Name, &s.Class, &s.Age)
	if err == sql.ErrNoRows {
		return nil, nil // Siswa tidak ditemukan
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get student by ID: %v", err)
	}
	return &s, nil
}

func AddStudent(s Student) (int64, error) {
	result, err := db.Exec("INSERT INTO students (name, class, age) VALUES (?, ?, ?)", s.Name, s.Class, s.Age)
	if err != nil {
		return 0, fmt.Errorf("failed to insert student: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}
	return id, nil
}

func UpdateStudent(s Student) error {
	_, err := db.Exec("UPDATE students SET name = ?, class = ?, age = ? WHERE id = ?", s.Name, s.Class, s.Age, s.ID)
	if err != nil {
		return fmt.Errorf("failed to update student: %v", err)
	}
	return nil
}

func DeleteStudent(id int) error {
	_, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %v", err)
	}
	return nil
}
