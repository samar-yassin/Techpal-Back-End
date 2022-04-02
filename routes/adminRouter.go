package routes

import (
	controller "CareerGuidance/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("api/addTrack", controller.AddTrack())
	incomingRoutes.GET("api/getalltracks", controller.GetAllTracks())
}
