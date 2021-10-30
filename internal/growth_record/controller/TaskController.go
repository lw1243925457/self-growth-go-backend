package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	modelV1 "seltGrowth/internal/api/v1"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
	"strconv"
)

type TaskController struct {
	srv srvV1.TaskService
}

func NewTaskController() *TaskController {
	return &TaskController{
		srv: srvV1.NewTaskService(),
	}
}

func (t *TaskController) TaskList(c *gin.Context) {
	isComplete := c.Query("isComplete")
	tasks, err := t.srv.GetTaskList(isComplete, GetLoginUserName(c))
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, tasks)
}

func (t *TaskController) AddTask(c *gin.Context) {
	var taskConfig modelV1.TaskConfig
	err := c.BindJSON(&taskConfig)
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}

	taskConfig.UserName = GetLoginUserName(c)
	err = t.srv.AddTask(taskConfig)
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, "新增任务成功")
}

func (t *TaskController) Complete(c *gin.Context) {
	id := c.Param("id")
	err := t.srv.Complete(id, GetLoginUserName(c))
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, "任务完成")
}

func (t *TaskController) History(c *gin.Context) {
	data, err := t.srv.History(GetLoginUserName(c))
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, data)
}

func (t *TaskController) AddTaskGroup(c *gin.Context) {
	var taskGroup modelV1.TaskGroup
	err := c.BindJSON(&taskGroup)
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}

	taskGroup.UserName = GetLoginUserName(c)
	err = t.srv.AddTaskGroup(taskGroup)
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponseWithoutData(c)
}

func (t *TaskController) TaskListByGroup(c *gin.Context) {
	data, err := t.srv.TaskListByGroup(GetLoginUserName(c))
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, data)
}

func (t *TaskController) Overview(c *gin.Context) {
	startTimeStamp, err := strconv.Atoi(c.Query("startTimeStamp"))
	endTimeStamp, err := strconv.Atoi(c.Query("endTimeStamp"))
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}
	data, err := t.srv.Overview(GetLoginUserName(c), int64(startTimeStamp), int64(endTimeStamp))
	if err != nil {
		log.Error(err)
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponse(c, data)
}

func (t *TaskController) DeleteGroup(c *gin.Context) {
	groupName := c.Param("name")
	err := t.srv.DeleteTaskGroup(groupName, GetLoginUserName(c))
	if err != nil {
		ErrorResponse(c,400, err.Error())
		return
	}
	SuccessResponseWithoutData(c)
}
