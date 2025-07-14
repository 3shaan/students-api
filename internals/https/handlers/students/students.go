package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/3shaan/students-api/internals/storage"
	"github.com/3shaan/students-api/internals/types"
	"github.com/3shaan/students-api/internals/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		decodeErr := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(decodeErr, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if decodeErr != nil {
			response.WriteJson(w, http.StatusBadRequest, decodeErr)
			return
		}

		// validations
		validateError := validator.New().Struct(student)
		if validateError != nil {

			vErr := validateError.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidatorError(vErr))
			return
		}

		// insert into database
		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}
		slog.Info("User created with id %s", slog.String("user id", fmt.Sprint(lastId)))

		response.WriteJson(w, http.StatusCreated, map[string]int64{"userId": (lastId)})
	}

}

func GetAll(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, result)

	}
}

func GetStudentById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusOK, response.GeneralError(err))
			return
		}
		result, err := storage.GetStudentById(intId)
		if err != nil {
			response.WriteJson(w, http.StatusOK, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, result)

	}
}
