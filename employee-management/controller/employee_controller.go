package controller

import (
	"fmt"

	"employee-management/service"
	"employee-management/view"
)

type EmployeeController interface {
	Start()
	Process(choice int)
}

type employeeController struct {
	view    view.EmployeeView
	service service.EmployeeService
}

func NewEmployeeController(view view.EmployeeView, service service.EmployeeService) EmployeeController {
	return &employeeController{
		view:    view,
		service: service,
	}
}

func (c *employeeController) Start() {
	for {
		choice := c.view.ShowMenu()
		if choice == 3 {
			fmt.Println("Thank You!")
			return
		}
		c.Process(choice)
		fmt.Println()
	}
}

func (c *employeeController) Process(choice int) {
	switch choice {
	case 1:
		employee := c.view.ReadEmployee()
		err := c.service.AddEmployee(employee)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Employee added successfully")
		}
	case 2:
		employees, err := c.service.GetEmployees()
		if err != nil {
			fmt.Println("Error retrieving employees:", err)
			return
		}
		if len(employees) == 0 {
			fmt.Println("No employees found.")
			return
		}
		fmt.Printf("%-8s%-12s%-20s%-5s\n", "ID", "Name", "Email", "Age")
		fmt.Println("------------------------------------------------")
		for _, emp := range employees {
			fmt.Printf("%-8d%-12s%-20s%-5d\n", emp.ID, emp.Name, emp.Email, emp.Age)
		}
	default:
		fmt.Println("Invalid choice")
	}
}