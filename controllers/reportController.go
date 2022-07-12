package controllers

import (
	"CareerGuidance/models"
	"github.com/gin-gonic/gin"
	gomail "gopkg.in/gomail.v2"
	"net/http"
)

func ReportMentor() gin.HandlerFunc {
	return func(c *gin.Context) {
		mentoremail := c.Param("mentor_email")
		var report models.Report

		err := c.BindJSON(&report)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		report.Mentor_email = mentoremail
		subject := "Mentor Report"
		body := "Mentor" + mentoremail + "<br><br>Message" + *report.Message

		msg := gomail.NewMessage()
		msg.SetHeader("From", "from@gmail.com")
		msg.SetHeader("To", "to@gmail.com")
		msg.SetHeader("Subject", subject)
		msg.SetBody("text/html", body)

		n := gomail.NewDialer("smtp.gmail.com", 587, "from@gmail.com", "password")

		// Send the email
		if err := n.DialAndSend(msg); err != nil {
			panic(err)
		}

	}
}
