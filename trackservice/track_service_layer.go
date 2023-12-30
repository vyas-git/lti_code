// service.go

package trackservice

import (
	"context"
	"errors"
	"fmt"

	trackdao "github.com/vyas-git/lti_code_test/TrackDao"
	"github.com/vyas-git/lti_code_test/model"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyClient struct {
	ClientID     string
	ClientSecret string
}

type TrackService struct {
	trackDAO      *trackdao.TrackDAO
	spotifyClient *SpotifyClient
}

func NewtrackService(trackDAO *trackdao.TrackDAO, spotifyClient *SpotifyClient) *TrackService {
	return &TrackService{trackDAO: trackDAO, spotifyClient: spotifyClient}
}

func NewSpotifyClient(clientID, clientSecret string) *SpotifyClient {
	return &SpotifyClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

func (s *TrackService) GetSpotifyClient() *SpotifyClient {
	return s.spotifyClient
}

func (s *TrackService) CreateTrack(isrc string) (*model.Track, error) {
	// checks for avilable tracks match with isrc
	existingTrack, _ := s.trackDAO.GetTrackByISRC(isrc)
	// if no existing data was found then it searches form Spotify api
	if existingTrack == nil || existingTrack.ISRC == "" {
		trackMetadata, err := s.spotifyClient.GetTrackMetadata(isrc)
		if err != nil {
			return nil, err
		}

		// Create a new track object
		track := &model.Track{
			ISRC:         isrc,
			SpotifyImage: trackMetadata.SpotifyImage,
			Title:        trackMetadata.Title,
			ArtistNames:  trackMetadata.ArtistNames,
			Popularity:   trackMetadata.Popularity,
		}

		// Store the track in the database
		err = s.trackDAO.CreateTrack(track)
		if err != nil {
			return nil, err
		}

		return track, nil
	}

	// If existingTrack is not nil and ISRC is not empty, return an error
	return nil, fmt.Errorf("track already found")
}

func (sc *SpotifyClient) GetTrackMetadata(isrc string) (*model.Track, error) {
	// for retriveing datas from spotify api
	config := &clientcredentials.Config{
		ClientID:     sc.ClientID,
		ClientSecret: sc.ClientSecret,
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		fmt.Println("Error getting Spotify token:", err)
		return nil, err
	}

	client := spotify.Authenticator{}.NewClient(token)

	result, err := client.Search(fmt.Sprintf("isrc:%s", isrc), spotify.SearchTypeTrack)
	if err != nil {
		fmt.Println("Error searching for track on Spotify:", err)
		return nil, err
	}
	// checks to for results to found on db
	if len(result.Tracks.Tracks) == 0 {
		fmt.Println("Track not found on Spotify")
		return nil, errors.New("track not found on Spotify")
	}

	// to return the highest popular track
	highestPopularityTrack := result.Tracks.Tracks[0]
	for _, track := range result.Tracks.Tracks {
		if track.Popularity > highestPopularityTrack.Popularity {
			highestPopularityTrack = track
		}
	}

	trackMetadata := &model.Track{
		SpotifyImage: highestPopularityTrack.Album.Images[0].URL,
		Title:        highestPopularityTrack.Name,
		ArtistNames:  getArtistNames(highestPopularityTrack.Artists),
		Popularity:   highestPopularityTrack.Popularity,
	}

	return trackMetadata, nil
}

func getArtistNames(artists []spotify.SimpleArtist) []string {

	var names []string
	for _, artist := range artists {
		names = append(names, artist.Name)
	}
	return names
}

func (s *TrackService) GetTrackByISRC(isrc string) (*model.Track, error) {
	// for retriveing datas from DB
	track, err := s.trackDAO.GetTrackByISRC(isrc)
	if err == nil {
		return track, nil
	}

	trackMetadata, err := s.spotifyClient.GetTrackMetadata(isrc)
	if err != nil {
		return nil, err
	}

	track = &model.Track{
		ISRC:         isrc,
		SpotifyImage: trackMetadata.SpotifyImage,
		Title:        trackMetadata.Title,
		ArtistNames:  trackMetadata.ArtistNames,
		Popularity:   trackMetadata.Popularity,
	}
	// creates track on DB
	err = s.trackDAO.CreateTrack(track)
	if err != nil {
		return nil, err
	}

	return track, nil
}

func (s *TrackService) GetTracksByArtist(artist string) (*[]model.Track, error) {
	// checks for artists on db
	tracks, err := s.trackDAO.GetTracksByArtist(artist)

	if err != nil {
		return nil, err
	}
	if len(*tracks) == 0 {
		return nil, err
	}
	return tracks, nil
}

func (s *TrackService) UpdateTrackByISRC(isrc string, updatedTrack *model.Track) (*model.Track, error) {
	// Check if the track exists in the database
	existingTrack, err := s.trackDAO.GetTrackByISRC(isrc)
	if err != nil {
		return nil, errors.New("track not found")
	}

	existingTrack.SpotifyImage = updatedTrack.SpotifyImage
	existingTrack.Title = updatedTrack.Title
	existingTrack.ArtistNames = updatedTrack.ArtistNames
	existingTrack.Popularity = updatedTrack.Popularity

	err = s.trackDAO.UpdateTrack(existingTrack)
	if err != nil {
		return nil, errors.New("failed to update track")
	}

	return existingTrack, nil
}
