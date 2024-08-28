package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type todoService struct {
	db *database
}

func (s *todoService) AddTodo(c *gin.Context) {
	var t userTask
	if err := c.BindJSON(&t); err != nil {
		c.String(400, "Неправильный формат данных.")
		return
	}
	err := s.db.AddTask(t)
	if err != nil {
		c.String(500, "Проблема на сервере.")
		return
	}

	c.String(201, "Created")
}

func (s *todoService) GetTodos(c *gin.Context) {
	tasks, err := s.db.GetTasks()
	if err != nil {
		c.String(500, "Проблема на сервере.")
		return
	}

	c.JSON(200, tasks)
}

func (s *todoService) GetTodo(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.String(400, "Неправильный формат данных.")
		return
	}
	task, err := s.db.GetTask(id)
	if err == errNotFound {
		c.String(404, "Задача не найдена.")
		return
	} else if err != nil {
		c.String(500, "Проблема на сервере.")
		return
	}
	c.JSON(200, task)
}

func (s *todoService) UpdateTodo(c *gin.Context) {
	var t userTask
	if err := c.BindJSON(&t); err != nil {
		c.String(400, "Неправильный формат данных.")
		return
	}
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.String(400, "Неправильный формат данных.")
		return
	}
	task, err := s.db.UpdateTask(id, t)
	if err == errNotFound {
		c.String(404, "Задача не найдена.")
		return
	} else if err != nil {
		c.String(500, "Проблема на сервере.")
		return
	}

	c.JSON(200, task)
}

func (s *todoService) DeleteTodo(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.String(400, "Неправильный формат данных.")
		return
	}

	err = s.db.DeleteTask(id)
	if err == errNotFound {
		c.String(404, "Задача не найдена.")
		return
	} else if err != nil {
		c.String(500, "Проблема на сервере.")
		return
	}
	c.String(204, "Задача удалена.")
}
