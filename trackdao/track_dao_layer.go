// trackdao.go

package trackdao

import (
	"github.com/vyas-git/lti_code_test/model"
	"gorm.io/gorm"
)

type TrackDAO struct {
	DB *gorm.DB
}

func NewTrackDAO(db *gorm.DB) *TrackDAO {
	return &TrackDAO{
		DB: db,
	}
}

func (dao *TrackDAO) CreateTrack(track *model.Track) error {

	if err := dao.DB.Create(track).Error; err != nil {
		return err
	}
	return nil
}

func (dao *TrackDAO) GetTrackByISRC(isrc string) (*model.Track, error) {
	var track model.Track
	err := dao.DB.Where("isrc = ?", isrc).First(&track).Error
	if err != nil {
		return nil, err
	}
	return &track, nil
}

func (dao *TrackDAO) GetTracksByArtist(artist string) (*[]model.Track, error) {
	var tracks *[]model.Track
	err := dao.DB.Where("array_to_string(artist_names, '||') ILIKE ?", "%"+artist+"%").Order("popularity DESC").Find(&tracks).Error
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func (dao *TrackDAO) UpdateTrack(track *model.Track) error {
	err := dao.DB.Save(track).Error
	if err != nil {
		return err
	}

	return nil
}
