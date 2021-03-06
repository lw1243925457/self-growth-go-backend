package controller_game_text

import (
	"github.com/gin-gonic/gin"
	"seltGrowth/internal/growth_record/controller"
	"seltGrowth/internal/growth_record/service/service_game_text"
	"strconv"
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

func (c *HeroController) ModifyOwnHeroProperty(ctx *gin.Context) {
	heroName := ctx.Query("hero")
	property := ctx.Query("property")
	modifyType := ctx.Query("type")
	err := c.srv.ModifyOwnHeroProperty(heroName, property, modifyType, controller.GetLoginUserName(ctx))
	if err != nil {
		controller.ErrorResponse(ctx, 400, err.Error())
		return
	}
	controller.SuccessResponseWithoutData(ctx)
}

func (c *HeroController) BattleHero(ctx *gin.Context) {
	heroName := ctx.Query("hero")
	err := c.srv.BattleHero(heroName, controller.GetLoginUserName(ctx))
	if err != nil {
		controller.ErrorResponse(ctx, 400, err.Error())
		return
	}
	controller.SuccessResponseWithoutData(ctx)
}

func (c *HeroController) BattleLog(ctx *gin.Context) {
	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil {
		controller.ErrorResponse(ctx, 400, err.Error())
		return
	}
	pageIndex, err := strconv.Atoi(ctx.Query("pageIndex"))
	if err != nil {
		controller.ErrorResponse(ctx, 400, err.Error())
		return
	}
	data, total, err := c.srv.BattleLog(controller.GetLoginUserName(ctx), pageIndex, pageSize)
	if err != nil {
		controller.ErrorResponse(ctx, 400, err.Error())
		return
	}
	controller.SuccessResponseWithPage(ctx, data, int64(pageSize), int64(pageIndex), total)
}