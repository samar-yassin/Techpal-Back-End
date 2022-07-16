package routes

import (
	controller "CareerGuidance/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("api/addTrack", controller.AddTrack())
	incomingRoutes.GET("api/gettrack/:track_id", controller.GetTrack())
	incomingRoutes.GET("api/getalltracks", controller.GetAllTracks())
	incomingRoutes.DELETE("api/deletetrack", controller.DeleteTrack())
	incomingRoutes.POST("api/acceptmentor", controller.AcceptMentor())
	incomingRoutes.POST("api/reportmentor/:mentor_email", controller.ReportMentor())
	incomingRoutes.GET("api/getacceptedmentors", controller.GetAcceptedMentors())
	incomingRoutes.GET("api/getnotacceptedmentors", controller.GetNotAcceptedMentors())
	incomingRoutes.GET("api/addskill/:skill_name", controller.AddSkill())
	incomingRoutes.GET("api/getallskills", controller.GetAllSkills())
}
