package dao

import (
	"employee-management/model"
)

type EmployeeDAO interface {
	Save(employee model.Employee) error
	FindAll() ([]model.Employee, error)
}

type memoryEmployeeDAO struct {
	employees []model.Employee
}

func NewEmployeeDAO() EmployeeDAO {
	return &memoryEmployeeDAO{
		employees: make([]model.Employee, 0),
	}
}

func (d *memoryEmployeeDAO) Save(employee model.Employee) error {
	d.employees = append(d.employees, employee)
	return nil
}

func (d *memoryEmployeeDAO) FindAll() ([]model.Employee, error) {
	return d.employees, nil
}