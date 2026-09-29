package repo

import (
	"gorm.io/gorm"
)

// Interface kế thừa Base
// type DeveloperRepoInterface interface {
// 	// crud2.BaseRepoInterface[domain.Developer]
// 	UpdateDevelopers(c *gin.Context, Developers []domain.Developer, attributeID uint64) error
// 	UpdateStatusDevelopers(c *gin.Context, Developers []domain.Developer, status int, archived int) error
// 	UpdateDeveloper2(c *gin.Context, Developers []domain.Developer, attributeID uint64) error
// }

type DeveloperRepo struct {
	// crud2.BaseRepo[domain.Developer]
}

func NewDeveloperRepo(db *gorm.DB) *DeveloperRepo {
	return &DeveloperRepo{
		// BaseRepo: crud2.BaseRepo[domain.Developer]{DB: db},
	}
}
