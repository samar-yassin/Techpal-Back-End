package routes

import (
	controller "CareerGuidance/controllers"

	"github.com/gin-gonic/gin"
)

func MentorRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("api/addsession/:user_id", controller.AddSession())
	incomingRoutes.GET("api/getallsessions/:user_id", controller.GetAllSessions())
	incomingRoutes.GET("api/removesession/:session_id", controller.RemoveSession())
}
