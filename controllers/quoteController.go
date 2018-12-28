package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sdstolworthy/go-fly/environment"
	"github.com/sdstolworthy/go-fly/models"
	skyscanner "github.com/sdstolworthy/go-skyscanner"
)

// TODO: Remove
var params = skyscanner.Parameters{
	Adults:           1,
	Country:          "US",
	Currency:         "USD",
	Locale:           "en-US",
	OriginPlace:      "BNA-sky",
	DestinationPlace: "BNA-sky",
	OutbandDate:      "anytime",
	InboundDate:      "anytime",
}

type quoteParameters struct {
	skyscanner.Parameters
}

// QuoteController handles all quote routes
type QuoteController struct {
	BaseController
}

// SetRoutes initializes the routes for the Quote Controller
func (c *QuoteController) SetRoutes(router *gin.RouterGroup) {
	c.setRouter(router)
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong")
	})
	router.POST("/getQuote", getQuote)
	router.GET("/allQuotes", allQuotes)
}

func getQuote(context *gin.Context) {
	var requestBody quoteParameters
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newQuotes []*skyscanner.QuoteSummary
	quoteChannels := make(chan *skyscanner.QuoteSummary)
	for _, v := range DestinationAirports {
		go skyscanner.ProcessDestination(v, &requestBody.Parameters, quoteChannels)
	}
	for range DestinationAirports {
		q := <-quoteChannels
		if q == nil {
			continue
		}
		newQuotes = append(newQuotes, q)
		environment.Env.Db.AddQuote(&models.Quote{
			Price:              q.Price,
			DestinationAirport: q.DestinationCity,
			OriginAirport:      q.OriginCity,
		})
	}
	context.JSON(http.StatusOK, newQuotes)
}

func allQuotes(context *gin.Context) {
	quotes, _ := environment.Env.Db.AllQuotes()
	context.JSON(http.StatusOK, quotes)
}
