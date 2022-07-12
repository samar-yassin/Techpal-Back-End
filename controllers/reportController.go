package controllers

import (
	"CareerGuidance/models"
	"github.com/gin-gonic/gin"
	gomail "gopkg.in/gomail.v2"
	"net/http"
)

func ReportMentor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var report models.Report

		err := c.BindJSON(&report)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		subject := "Report!"
		body := "Report message " + *report.Message

		msg := gomail.NewMessage()
		msg.SetHeader("From", "techpal.guidance@gmail.com")
		msg.SetHeader("To", "farabi.marwa@gmail.com")
		msg.SetHeader("Subject", subject)
		msg.SetBody("text/html", body)

		n := gomail.NewDialer("smtp.gmail.com", 587, "techpal.guidance@gmail.com", "osijygroequycuww")

		// Send the email
		if err := n.DialAndSend(msg); err != nil {
			panic(err)
		}

	}
}
