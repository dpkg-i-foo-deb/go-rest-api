package services

import (
	"backend/auth"
	"backend/database"
	"backend/models"
	"backend/models/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
)

func CreateTaskService(connection *fiber.Ctx) error {

	var task *models.Task
	var tokenString string
	decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(connection.Body())))
	var response utils.GenericResponse
	err := decoder.Decode(&task)

	if err != nil {

		response.Response = "The task is malformed"
		connection.Status(fiber.StatusBadRequest).JSON(response)
		return nil

	}

	tokenString = connection.Cookies("access-token")

	//We retrieve the token string from the request cookie
	task.UserEmail = new(string)
	*task.UserEmail, err = auth.EmailFromToken(tokenString)

	if err != nil {

		return errors.New("Could not create the task, try again later")

	}

	err = database.CreateTaskStatement.QueryRow(
		task.Title, task.Description, task.UserEmail, task.StartDate, task.DueDate, task.Status,
		task.MainTask,
	).Scan(
		&task.Title, &task.Description, &task.UserEmail,
		&task.StartDate, &task.DueDate, &task.Status, &task.MainTask, &task.Code,
	)

	if err != nil {

		return errors.New("Could not create the task, try again later")

	}

	connection.Status(fiber.StatusCreated).JSON(task)
	return nil
}

func GetTaskService(connection *fiber.Ctx) error {

	var taskCode int
	var tokenString string

	var userEmail string

	var task models.Task

	var response utils.GenericResponse

	taskCode, err := connection.ParamsInt("code")

	if taskCode == 0 || err != nil {

		response.Response = "The received parameter is not valid"
		connection.Status(fiber.StatusBadRequest).JSON(response)
		return nil
	}

	tokenString = connection.Cookies("access-token")

	userEmail, err = auth.EmailFromToken(tokenString)

	if err != nil {
		return errors.New("Could not get the user's email from the token")
	}

	err = database.GetTaskStatement.QueryRow(
		taskCode, userEmail,
	).Scan(
		&task.Title, &task.Description, &task.Code,
		&task.MainTask, &task.UserEmail, &task.StartDate, &task.DueDate, &task.Status,
	)

	if err != nil {

		response.Response = "The task does not exist or you have no access to it"
		connection.Status(fiber.StatusBadRequest).JSON(response)
		return nil
	}

	connection.Status(fiber.StatusOK).JSON(task)
	return nil
}

func GetAllTasksService(connection *fiber.Ctx) error {
	var tokenString string

	var userEmail string

	var task models.Task

	var allTasks []models.Task

	var response utils.GenericResponse

	var rows *sql.Rows

	tokenString = connection.Cookies("access-token")

	userEmail, err := auth.EmailFromToken(tokenString)

	if err != nil {
		return errors.New("Could not get the user's email from the token")
	}

	rows, err = database.GetAllTasksStatement.Query(userEmail)

	if err != nil {
		response.Response = "You have no tasks created"
		connection.Status(fiber.StatusOK).JSON(allTasks)
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&task.Title, &task.Description, &task.Code,
			&task.MainTask, &task.UserEmail, &task.StartDate, &task.DueDate, &task.Status)

		if err != nil {
			return errors.New("Could not retrieve your tasks")
		}
		allTasks = append(allTasks, task)

	}

	err = rows.Err()

	if err != nil {
		return errors.New("Could not retrieve your tasks")
	}

	connection.Status(fiber.StatusOK).JSON(allTasks)

	return nil
}

func EditTaskService(connection *fiber.Ctx) error {

	var tokenString string

	var userEmail string

	var task models.Task

	taskCode, err := connection.ParamsInt("code")

	var response utils.GenericResponse

	var reader = ioutil.NopCloser(bytes.NewBuffer(connection.Body()))

	if taskCode == 0 || err != nil {

		response.Response = "The received parameter is not valid"
		connection.Status(fiber.StatusBadRequest).JSON(response)
	}

	tokenString = connection.Cookies("access-token")

	userEmail, err = auth.EmailFromToken(tokenString)

	if err != nil {

		return errors.New("Could not get the user's email from the token")

	}

	decoder := json.NewDecoder(reader)

	err = decoder.Decode(&task)

	if err != nil {
		response.Response = "The task is malformed"
		connection.Status(fiber.StatusBadRequest).JSON(response)
	}

	task.UserEmail = &userEmail

	code := int(taskCode)

	task.Code = &code

	err = database.EditTaskStatement.QueryRow(task.Title,
		task.Description, task.MainTask, task.StartDate,
		task.DueDate, task.Status, task.Code, task.UserEmail).Scan(
		&task.Title, &task.Description, &task.MainTask, &task.StartDate,
		&task.DueDate, &task.Status, &task.Code,
	)

	if err != nil {
		return errors.New("Could not edit the task, try again later")
	}

	connection.Status(fiber.StatusAccepted).JSON(task)

	return nil

}

func DeleteTaskService(connection *fiber.Ctx) error {
	var tokenString string

	var userEmail string

	var task models.Task

	taskCode, err := connection.ParamsInt("code")

	var response utils.GenericResponse

	if taskCode == 0 || err != nil {
		response.Response = "The received parameter is not valid"
		connection.Status(fiber.StatusBadRequest).JSON(response)
		return nil
	}

	tokenString = connection.Cookies("access-token")

	userEmail, err = auth.EmailFromToken(tokenString)

	if err != nil {
		return errors.New("Could not delete the task, try again later")
	}

	task.UserEmail = &userEmail

	code := int(taskCode)

	task.Code = &code

	result, err := database.DeleteTaskStatement.Exec(task.Code, task.UserEmail)

	if err != nil {
		return errors.New("Could not delete the task, try again later")
	}

	affectedRows, err := result.RowsAffected()

	if err != nil || affectedRows == 0 {
		return errors.New("Could not delete the task, try again later")
	}

	connection.Status(fiber.StatusOK)
	return nil
}
