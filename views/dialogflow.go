package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Depado/assistant-codelabs/dialogflow"
	"github.com/Depado/assistant-codelabs/models"
)

// PostFulfillment is the handler handling incoming webhook requests from
// dialogflow
func PostFulfillment(c *gin.Context) {
	var err error
	var car *models.Car
	var dfr dialogflow.Request

	if err = c.BindJSON(&dfr); err != nil {
		logrus.WithError(err).Error("Couldn't bind JSON")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if car, err = models.NewCarFromParameters(dfr.Result.Parameters); err != nil {
		logrus.WithError(err).Warn("Couldn't parse parameters")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	out := car.Estimate()
	c.JSON(http.StatusOK, dialogflow.Response{
		DisplayText: out,
		Speech:      out,
	})
}
