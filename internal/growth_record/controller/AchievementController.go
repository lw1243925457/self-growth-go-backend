package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	srvV1 "seltGrowth/internal/growth_record/service/v1"
	"strconv"
	"time"
)

type AchievementController struct {
	srv srvV1.AchievementService
}

func NewAchievementController() *AchievementController {
	return &AchievementController{
		srv: srvV1.NewAchievementService(),
	}
}

func (a *AchievementController) Get(c *gin.Context) {
	timestamp, err := strconv.Atoi(c.Query("timestamp"))
	if err != nil {
		ErrorResponse(c, 400, errors.New("请传入有效时间").Error())
		return
	}
	data, err := a.srv.Get(time.Unix(int64(timestamp), 0), GetLoginUserName(c))
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}

	SuccessResponse(c, data)
}

// Sync 每日成就同步（生成成就列表）
func (a *AchievementController) Sync(c *gin.Context) {
	timestamp, err := strconv.Atoi(c.Query("timestamp"))
	if err != nil {
		ErrorResponse(c, 400, errors.New("请传入有效时间").Error())
		return
	}
	err = a.srv.Sync(time.Unix(int64(timestamp), 0), GetLoginUserName(c))
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}

	SuccessResponseWithoutData(c)
}

func (a *AchievementController) Import(c *gin.Context) {
	err := a.srv.Import(c.Query("id"))
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}
	SuccessResponseWithoutData(c)
}