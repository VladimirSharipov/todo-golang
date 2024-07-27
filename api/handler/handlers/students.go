package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"yourmodule/db"
	"yourmodule/models"
)

func HandleStudents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getStudents(w, r)
	case http.MethodPost:
		createStudent(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleStudent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getStudent(w, r)
	case http.MethodPut:
		updateStudent(w, r)
	case http.MethodDelete:
		deleteStudent(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, age, course FROM students")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.Course); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/students/"):]
	var student models.Student
	err := db.DB.QueryRow("SELECT id, name, age, course FROM students WHERE id = $1", id).Scan(&student.ID, &student.Name, &student.Age, &student.Course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.DB.QueryRow("INSERT INTO students(name, age, course) VALUES($1, $2, $3) RETURNING id", student.Name, student.Age, student.Course).Scan(&student.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/students/"):]
	var student models.Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("UPDATE students SET name=$1, age=$2, course=$3 WHERE id=$4", student.Name, student.Age, student.Course, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	student.ID, _ = strconv.Atoi(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/students/"):]
	_, err := db.DB.Exec("DELETE FROM students WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
