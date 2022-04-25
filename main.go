package main

import (
	"github.com/gin-gonic/gin"

	"os"

	"github.com/kagonlineteam/homework-microservice/controller"
	"github.com/kagonlineteam/homework-microservice/middleware"
	"github.com/kagonlineteam/homework-microservice/models"
)

func main() {
	models.ConnectDatabase()

	r := gin.Default()
	r.Use(middleware.AuthMiddleware())
	r.SetTrustedProxies([]string{os.Getenv("HOMEWORK_PROXY_IP")})

	api := r.Group("/homework/v1")

	api.GET("my", controller.GetOwnHomeworks)
	api.POST("homeworks", controller.CreateHomework)
	api.PUT("homeworks/:id", controller.EditHomework)
	api.POST("report/:id", controller.ReportHomework)

	r.Run()
}
