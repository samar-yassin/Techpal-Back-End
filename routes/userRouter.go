package routes

import (
	controller "CareerGuidance/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	//incomingRoutes.Use(middleware.Authenticate())
	//incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("api/createprofile/:user_id", controller.CreateProfile())
	incomingRoutes.POST("api/switchprofile/:user_id", controller.SwitchProfile())
	incomingRoutes.GET("api/user", controller.User())

}
