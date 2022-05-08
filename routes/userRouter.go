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
	incomingRoutes.GET("api/deleteprofile/:user_id", controller.DeleteProfile())
	incomingRoutes.GET("api/getcurrentprofile/:user_id", controller.GetCurrentProfile())
	incomingRoutes.GET("api/getallprofiles", controller.GetAllProfiles())
	incomingRoutes.GET("api/user", controller.User())
	incomingRoutes.GET("api/getuser/:user_id", controller.GetUser())
	incomingRoutes.POST("api/updatestudent/:user_id", controller.UpdateStudent())
	incomingRoutes.POST("api/updatementor/:user_id", controller.UpdateMentor())
	incomingRoutes.POST("api/changepassword/:user_id", controller.ChangePassword())
}
