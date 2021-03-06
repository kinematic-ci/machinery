package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kinematic-ci/machinery/server/executors"
	"github.com/kinematic-ci/machinery/server/mappers"
	"github.com/kinematic-ci/machinery/server/responses"
	"github.com/kinematic-ci/machinery/server/tasks"
	"net/http"
)

type TaskController struct {
	tasks    []tasks.Task
	executor executors.Executor
}

func NewTaskController(tasks []tasks.Task, executor executors.Executor) *TaskController {
	return &TaskController{tasks: tasks, executor: executor}
}

func (t TaskController) List(c *gin.Context) {
	c.JSON(http.StatusOK, mappers.TasksToResponse(t.tasks))
}

func (t TaskController) Get(c *gin.Context) {
	id := c.Param("id")
	if task, found := t.findTask(id); found {
		c.JSON(http.StatusOK, mappers.TaskToTaskDetail(task))
		return
	}

	c.JSON(http.StatusNotFound, responses.NotFound)
}

func (t TaskController) findTask(id string) (tasks.Task, bool) {
	for _, task := range t.tasks {
		if task.Id == id {
			return task, true
		}
	}
	return tasks.Task{}, false
}

func (t TaskController) Execute(c *gin.Context) {
	id := c.Param("id")
	task, found := t.findTask(id)

	if !found {
		c.JSON(http.StatusNotFound, responses.NotFound)
		return
	}

	var executeArgs gin.H

	err := c.BindJSON(&executeArgs)

	if err != nil {
		c.JSON(http.StatusBadRequest, responses.BadRequest)
	}

	result, err := t.executor.Execute(c, task, makeParams(executeArgs))

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.ResultToResponse(result))
	return
}

func makeParams(args gin.H) []string {
	params := make([]string,0, len(args))

	for k, v := range args {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}

	return params
}
