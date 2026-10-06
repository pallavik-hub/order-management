package utility

import (
	"errors"
	"strings"

	"employee-management/model"
)

type Validator interface {
	ValidateEmployee(employee model.Employee) error
	ValidateEmail(email string) error
}

type employeeValidator struct{}

func NewValidator() Validator {
	return &employeeValidator{}
}

func (v *employeeValidator) ValidateEmail(email string) error {
	if !strings.Contains(email, "@") {
		return errors.New("invalid email: must contain @")
	}
	return nil
}

func (v *employeeValidator) ValidateEmployee(employee model.Employee) error {
	if employee.ID <= 0 {
		return errors.New("invalid employee ID: must be greater than 0")
	}
	if strings.TrimSpace(employee.Name) == "" {
		return errors.New("invalid name: cannot be empty")
	}
	if err := v.ValidateEmail(employee.Email); err != nil {
		return err
	}
	if employee.Age < 18 || employee.Age > 60 {
		return errors.New("invalid age: must be between 18 and 60")
	}
	return nil
}