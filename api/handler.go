package api

import (
	"os"

	"github.com/gin-gonic/gin"
	trackdao "github.com/vyas-git/lti_code_test/TrackDao"
	trackservice "github.com/vyas-git/lti_code_test/TrackService"
	trackhandler "github.com/vyas-git/lti_code_test/Trackhandler"
	"gorm.io/gorm"
)

func Run(router *gin.Engine, db *gorm.DB) {
	// Access environment variables
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	trackDAO := trackdao.NewTrackDAO(db)
	trackService := trackservice.NewtrackService(trackDAO, trackservice.NewSpotifyClient(clientID, clientSecret))
	trackHandler := trackhandler.NewTrackHandler(trackService, trackService.GetSpotifyClient())

	router.POST("/tracks/create", trackHandler.CreateTrackHandler)

	router.GET("/tracks/:isrc", trackHandler.GetTrackByISRCHandler)

	router.GET("/tracks/artist/:artist", trackHandler.GetTracksByArtistHandler)

	router.PUT("/tracks/update/:isrc", trackHandler.UpdateTrackByISRCHandler)
}
