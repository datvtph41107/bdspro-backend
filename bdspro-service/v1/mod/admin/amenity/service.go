package a_amenity

import (
	"bdspro/v1/mod/global/amenity"
	"bdspro/v1/models"
	"common/case/crud2"
)

type AmenityService struct {
	crud2.BaseService[models.Amenity, *amenity.AmenityRepo]
}

func NewAmenityService(repo *amenity.AmenityRepo) *AmenityService {
	return &AmenityService{
		BaseService: crud2.BaseService[models.Amenity, *amenity.AmenityRepo]{Repo: repo},
	}
}
