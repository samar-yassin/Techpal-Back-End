package routes

import (
	controller "CareerGuidance/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("api/addTrack", controller.AddTrack())
	incomingRoutes.GET("api/gettrack/:track_id", controller.GetTrack())
	incomingRoutes.GET("api/getalltracks", controller.GetAllTracks())
	incomingRoutes.POST("api/acceptmentor", controller.AcceptMentor())
	incomingRoutes.POST("api/reportmentor", controller.ReportMentor())

}
