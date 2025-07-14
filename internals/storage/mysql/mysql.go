package mySql

import (
	"database/sql"
	"fmt"

	"github.com/3shaan/students-api/internals/config"
	"github.com/3shaan/students-api/internals/types"
	_ "github.com/go-sql-driver/mysql"
)

type MySql struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*MySql, error) {

	dsn := cfg.DbUser + ":" + cfg.DbPassword + "@tcp(localhost:3306)/" + cfg.DbName
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	_, createTableErr := db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
	name varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	age int
	)`)

	if createTableErr != nil {
		return nil, createTableErr
	}
	return &MySql{
		Db: db,
	}, nil

}

func (s *MySql) CreateStudent(name string, email string, age int) (int64, error) {
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

func (s *MySql) GetStudents() ([]types.Student, error) {

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

func (s *MySql) GetStudentById(id int64) (types.Student, error) {
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
func (s MySql) DeleteStudentById(id int64) (string, error) {
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
