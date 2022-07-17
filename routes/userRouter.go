package routes

import (
	controller "CareerGuidance/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	//incomingRoutes.Use(middleware.Authenticate())
	//incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.POST("api/createprofile/:user_id", controller.CreateProfile())
	incomingRoutes.POST("api/switchprofile/:user_id", controller.SwitchProfile())
	incomingRoutes.DELETE("api/deleteprofile/:user_id", controller.DeleteProfile())
	incomingRoutes.POST("api/markcompleted", controller.MarkCompleted())
	incomingRoutes.POST("api/enroll", controller.EnrollCourse())
	incomingRoutes.GET("api/getcurrentprofile/:user_id", controller.GetCurrentProfile())
	incomingRoutes.GET("api/user", controller.User())
	incomingRoutes.GET("api/getuser/:user_id", controller.GetUser())
	incomingRoutes.POST("api/updatestudent/:user_id", controller.UpdateStudent())
	incomingRoutes.POST("api/updatementor/:user_id", controller.UpdateMentor())
	incomingRoutes.POST("api/changementorpassword/:user_id", controller.ChangeMentorPassword())
	incomingRoutes.GET("api/getallprofiles/:user_id", controller.GetAllProfiles())
	incomingRoutes.DELETE("api/removementor/:user_id", controller.RemoveMentor())
	incomingRoutes.POST("api/addresume", controller.AddResume())
	incomingRoutes.GET("api/getresume/:profile_id", controller.GetResume())
	incomingRoutes.POST("api/updateresume/:profile_id", controller.UpdateResume())
	incomingRoutes.GET("api/getallsessions", controller.GetAllSessions())
	incomingRoutes.GET("api/getenrolledcourses/:profile_id", controller.GetEnrolledCourses())
	incomingRoutes.GET("api/getcompletedcourses/:profile_id", controller.GetCompletedCourses())
	incomingRoutes.GET("api/getallusers", controller.GetAllUsers())
	incomingRoutes.GET("api/leadershipboard/:track_id", controller.LeadershipBoard())
	incomingRoutes.POST("api/deletecourse", controller.DeleteCourse())
	incomingRoutes.POST("api/rate", controller.RateCourse())
	incomingRoutes.POST("api/contact", controller.ContactUs())

}
