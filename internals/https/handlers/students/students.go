package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/3shaan/students-api/internals/types"
	"github.com/3shaan/students-api/internals/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
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

		response.WriteJson(w, http.StatusCreated, map[string]string{"Success": "OK"})
	}

}
