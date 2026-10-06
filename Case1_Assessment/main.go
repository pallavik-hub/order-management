package main

import (
	"employee-management/controller"
	"employee-management/dao"
	"employee-management/service"
	"employee-management/utility"
	"employee-management/view"
)

func main() {
	// Initialize dependencies and inject them bottom-up
	empDAO := dao.NewEmployeeDAO()
	validator := utility.NewValidator()
	empService := service.NewEmployeeService(empDAO, validator)
	empView := view.NewEmployeeView()

	empController := controller.NewEmployeeController(empView, empService)

	// Start application flow
	empController.Start()
}