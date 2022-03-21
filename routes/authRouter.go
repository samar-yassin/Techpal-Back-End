package routes

import (
	controller "CareerGuidance/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	//incomingRoutes.POST("api/signup", controller.Signup())
	incomingRoutes.POST("api/login", controller.Login())

}
