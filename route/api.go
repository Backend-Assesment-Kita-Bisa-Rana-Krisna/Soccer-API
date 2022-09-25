package route

import (
	controller "soccer-api/app/controller"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func ApiRoutes(route *gin.RouterGroup) {
	route.POST("v1/team", controller.CreateTeam())
	route.GET("v1/team/:id", controller.GetTeam())
	route.GET("v1/team/:id/players", controller.GetTeamWithPlayer())
	route.PUT("v1/team/:id", controller.UpdateTeam())
	route.DELETE("v1/team/:id", controller.DeleteTeam())
	route.GET("v1/teams", controller.GetAllTeam())

	route.POST("v1/player", controller.CreatePlayer())
	route.GET("v1/player/:id", controller.GetPlayer())
	route.PUT("v1/player/:id", controller.UpdatePlayer())
	route.DELETE("v1/player/:id", controller.DeletePlayer())
	route.GET("v1/players", controller.GetAllPlayer())
}
