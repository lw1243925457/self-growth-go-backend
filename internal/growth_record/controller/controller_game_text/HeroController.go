package controller_game_text

import (
	"github.com/gin-gonic/gin"
	"seltGrowth/internal/growth_record/controller"
	"seltGrowth/internal/growth_record/service/service_game_text"
)

type HeroController struct {
	srv service_game_text.HeroService
}

func NewHeroController() *HeroController {
	return &HeroController{
		srv: service_game_text.NewHeroService(),
	}
}

// List 获取所有的角色列表
func (c *HeroController) List(ctx *gin.Context) {
	data, err := c.srv.List()
	if err != nil {
		controller.ErrorResponse(ctx, 400, err.Error())
		return
	}
	controller.SuccessResponse(ctx, data)
}

func (c *HeroController) UserInfo(context *gin.Context) {
	data, err := c.srv.UserInfo(controller.GetLoginUserName(context))
	if err != nil {
		controller.ErrorResponse(context, 400, err.Error())
		return
	}
	controller.SuccessResponse(context, data)
}

func (c *HeroController) HeroRound(ctx *gin.Context) {
	data, err := c.srv.HeroRound(controller.GetLoginUserName(ctx))
	if err != nil {
		controller.ErrorResponse(ctx, 400, err.Error())
		return
	}
	controller.SuccessResponse(ctx, data)
}

func (c *HeroController) OwnHeroes(ctx *gin.Context) {
	data, err := c.srv.OwnHeroes(controller.GetLoginUserName(ctx))
	if err != nil {
		controller.ErrorResponse(ctx, 400, err.Error())
		return
	}
	controller.SuccessResponse(ctx, data)
}