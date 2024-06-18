package types

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type EnergyServiceEvent struct {
	ID             string               `json:"id" validate:"required"`
	ChildProcesses []string             `json:"childProcesses" validate:"required"`
	Context        EnergyServiceContext `json:"context"`
	Completed_At   string               `json:"completed_at"`
}

type EnergyServiceContext struct {
	RequestID string `json:"requestId" validate:"required"`
	Date      string `json:"date" validate:"required,customDate"`
	Status    string `json:"status"`
	Service   string `json:"service"`
}

func customErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required", err.Field())
	case "customDate":
		return fmt.Sprintf("The %s field must be in YYYY-MM-DD format", err.Field())
	default:
		return fmt.Sprintf("The %s field is invalid", err.Field())
	}
}

func validateYYYYMMDD(fl validator.FieldLevel) bool {
	// Regular expression to match YYYY-MM-DD format
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return re.MatchString(fl.Field().String())
}

func formatValidationErrors(errors validator.ValidationErrors) error {
	var errorMessages []string
	for _, err := range errors {
		errorMessages = append(errorMessages, customErrorMessage(err))
	}
	return fmt.Errorf(strings.Join(errorMessages, "\n"))
}

func Validate(val interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("customDate", validateYYYYMMDD)
	err := validate.Struct(val)

	if err != nil {
		return formatValidationErrors(err.(validator.ValidationErrors))
	}

	return nil
}
