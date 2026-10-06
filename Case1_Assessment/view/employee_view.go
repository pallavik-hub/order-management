package view

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"employee-management/model"
)

type EmployeeView interface {
	ShowMenu() int
	ReadEmployee() model.Employee
}

type consoleEmployeeView struct {
	reader *bufio.Reader
}

func NewEmployeeView() EmployeeView {
	return &consoleEmployeeView{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (v *consoleEmployeeView) ShowMenu() int {
	fmt.Println("================================")
	fmt.Println("      EMPLOYEE MANAGEMENT")
	fmt.Println("================================")
	fmt.Println("1. Add Employee")
	fmt.Println("2. View Employees")
	fmt.Println("3. Exit")
	fmt.Print("Enter your choice: ")

	input, _ := v.reader.ReadString('\n')
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return 0
	}
	return choice
}

func (v *consoleEmployeeView) ReadEmployee() model.Employee {
	var emp model.Employee

	fmt.Print("Enter Employee ID: ")
	idStr, _ := v.reader.ReadString('\n')
	emp.ID, _ = strconv.Atoi(strings.TrimSpace(idStr))

	fmt.Print("Enter Employee Name: ")
	emp.Name, _ = v.reader.ReadString('\n')
	emp.Name = strings.TrimSpace(emp.Name)

	fmt.Print("Enter Employee Email: ")
	emp.Email, _ = v.reader.ReadString('\n')
	emp.Email = strings.TrimSpace(emp.Email)

	fmt.Print("Enter Employee Age: ")
	ageStr, _ := v.reader.ReadString('\n')
	emp.Age, _ = strconv.Atoi(strings.TrimSpace(ageStr))

	return emp
}