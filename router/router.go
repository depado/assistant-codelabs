package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Depado/assistant-codelabs/conf"
	"github.com/Depado/assistant-codelabs/views"
)

// Run is a function to run the router
func Run() {
	var err error

	r := gin.Default()
	r.Use(AuthMiddleware())

	r.POST("/fulfillment", views.PostFulfillment)

	if err = r.Run(conf.C.ListenAddress()); err != nil {
		logrus.WithError(err).Fatal("Couldn't start server")
	}
}
