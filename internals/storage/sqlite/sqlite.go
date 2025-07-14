package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/3shaan/students-api/internals/config"
	"github.com/3shaan/students-api/internals/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, createTableErr := db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if createTableErr != nil {
		return nil, createTableErr
	}
	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil

}

func (s *Sqlite) GetStudents() ([]types.Student, error) {

	rows, err := s.Db.Query(`Select id, name, email, age FROM students`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)

	}
	return students, nil

}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students where id=? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()
	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with this %s", fmt.Sprint(id))

		}
		return types.Student{}, fmt.Errorf("query error %w", err)
	}

	return student, nil

}

// delete functions
// if deleted it will "OK"
func (s Sqlite) DeleteStudentById(id int64) (string, error) {
	result, err := s.Db.Exec("DELETE FROM students WHERE id=?", id)
	if err != nil {
		return "", err
	}

	num, err := result.RowsAffected()

	if err != nil {
		return "", err
	}

	if num == 0 {
		return "", fmt.Errorf("delete failed with id %s", fmt.Sprint(id))
	}
	return "OK", nil

}
