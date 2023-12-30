package trackhandler

import (
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	trackservice "github.com/vyas-git/lti_code_test/TrackService"
	"github.com/vyas-git/lti_code_test/model"
)

type TrackHandler struct {
	trackService  *trackservice.TrackService
	spotifyClient *trackservice.SpotifyClient
}

func NewTrackHandler(trackService *trackservice.TrackService, spotifyClient *trackservice.SpotifyClient) *TrackHandler {
	return &TrackHandler{trackService: trackService, spotifyClient: spotifyClient}
}

// @Summary Create a new track
// @Description Create a new track record in the database
// @Accept json
// @Produce json
// @tags Create or Update Tracks
// @Param track body model.RequestBody true "Track details to create"
// @Success 200 {object} model.TrackDetails	"Track details"
// @Failure 400 {object} model.ErrorResponse "Invalid request body"
// @Failure 409 {object} model.ErrorResponse "Track with ISRC code already exists"
// @Router /tracks/create [post]
func (h *TrackHandler) CreateTrackHandler(c *gin.Context) {
	// to store user isrc value
	var requestBody model.RequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// To verify the client have entered isrc value
	if requestBody.ISRC == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "ISRC is required"})
		return
	}

	track, err := h.trackService.CreateTrack(requestBody.ISRC)

	if err != nil {
		if strings.Contains(err.Error(), "track already found") {
			c.JSON(http.StatusConflict, gin.H{"error": "Track already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Track created successfully", "track": track})
}

// @Summary Search tracks by ISRC code
// @Description Search tracks from the database or Spotify by ISRC code
// @Produce json
// @tags GetTracks
// @Param isrc path string true "ISRC code of the track"
// @Success 200 {object} model.TrackDetails "Track details"
// @Failure 400 {object} model.ErrorResponse "Invalid request body"
// @Failure 404 {object} model.ErrorResponse "Track not found for the given ISRC"
// @Router /tracks/{isrc} [get]
func (h *TrackHandler) GetTrackByISRCHandler(c *gin.Context) {

	isrc := c.Param("isrc")

	track, err := h.trackService.GetTrackByISRC(isrc)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get track by ISRC"})
		return
	}

	c.JSON(200, gin.H{"track": track})
}

// @Summary Get tracks by artist
// @Description Retrieves tracks by the specified artist from the database or Spotify.
// @Produce json
// @tags GetTracks
// @Param artist path string true "Artist name"
// @Success 200 {object} model.TrackDetails "Successfully retrieved tracks"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /tracks/artist/{artist} [get]
func (h *TrackHandler) GetTracksByArtistHandler(c *gin.Context) {

	artist := c.Param("artist")
	// checks the db for avilable tracks with this artist
	tracks, err := h.trackService.GetTracksByArtist(artist)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get tracks by artist"})
		return
	}
	if tracks == nil {
		c.JSON(500, gin.H{"error": "No songs related to this artist"})
		return
	}
	c.JSON(200, gin.H{"tracks": tracks})
}

// @Summary Update a track by ISRC
// @Description Updates a track with the specified ISRC.
// @Accept json
// @Produce json
// @tags Create or Update Tracks
// @Param isrc path string true "ISRC of the track"
// @Param existingTrack body model.TrackDetails true "Updated track details"
// @Success 200 {object} model.TrackDetails "Track updated successfully"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 404 {object} model.ErrorResponse "Track not found"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /tracks/update/{isrc} [put]
func (h *TrackHandler) UpdateTrackByISRCHandler(c *gin.Context) {
	isrc := c.Param("isrc")

	var updatedTrack *model.Track

	if err := c.ShouldBindJSON(&updatedTrack); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	track, err := h.trackService.UpdateTrackByISRC(isrc, updatedTrack)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Track updated successfully", "track": track})
}
