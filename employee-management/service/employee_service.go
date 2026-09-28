package service

import (
	"employee-management/dao"
	"employee-management/model"
	"employee-management/utility"
)

type EmployeeService interface {
	AddEmployee(employee model.Employee) error
	GetEmployees() ([]model.Employee, error)
}

type employeeService struct {
	dao       dao.EmployeeDAO
	validator utility.Validator
}

func NewEmployeeService(dao dao.EmployeeDAO, validator utility.Validator) EmployeeService {
	return &employeeService{
		dao:       dao,
		validator: validator,
	}
}

func (s *employeeService) AddEmployee(employee model.Employee) error {
	if err := s.validator.ValidateEmployee(employee); err != nil {
		return err
	}
	return s.dao.Save(employee)
}

func (s *employeeService) GetEmployees() ([]model.Employee, error) {
	return s.dao.FindAll()
}