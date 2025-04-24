package a_developer

import (
	"bdspro/v1/models"
	"common/case/crud2"

	"gorm.io/gorm"
)

// Interface kế thừa Base
// type DeveloperRepoInterface interface {
// 	// crud2.BaseRepoInterface[models.Developer]
// 	UpdateDevelopers(c *gin.Context, Developers []models.Developer, attributeID uint64) error
// 	UpdateStatusDevelopers(c *gin.Context, Developers []models.Developer, status int, archived int) error
// 	UpdateDeveloper2(c *gin.Context, Developers []models.Developer, attributeID uint64) error
// }

type DeveloperRepo struct {
	crud2.BaseRepo[models.Developer]
}

func NewDeveloperRepo(db *gorm.DB) *DeveloperRepo {
	return &DeveloperRepo{
		BaseRepo: crud2.BaseRepo[models.Developer]{DB: db},
	}
}
