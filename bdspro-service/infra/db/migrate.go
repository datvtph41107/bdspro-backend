package db

import (
	_db "common/db"

	"bdspro/internal/domain"
)

func MigrateDomain() {
	_db.DB.AutoMigrate(
		// Core entities
		&domain.Asset{},
		&domain.Product{},
		&domain.Post{},
		&domain.Intent{},

		// Asset related
		&domain.AssetCost{},
		&domain.AssetCostType{},
		&domain.AssetIncomeType{},
		&domain.AssetExploitation{},
		&domain.AssetLegal{},
		&domain.AssetOrganization{},
		&domain.AssetUser{},

		// Product related
		&domain.ProductUser{},
		&domain.ProductOrganization{},
		&domain.ProductHistory{},
		&domain.ProductPrice{},
		&domain.ProductMedia{},
		&domain.ProductPrivate{},

		// Post related
		&domain.PostMediaEntity{},
		&domain.PostOrganization{},
		&domain.PostUser{},

		// Document related
		&domain.CostDocument{},
		&domain.IncomeDocument{},

		// Project related
		&domain.Project{},
		&domain.ProjectBuild{},
		&domain.ProjectInventory{},
		&domain.BuildRoom{},

		// Sharing and Access
		&domain.SharingAccess{},
		&domain.RecordHistory{},

		// Other entities
		&domain.Amenity{},
		&domain.Apartment{},
		&domain.Developer{},
		&domain.PropertyType{},
		// &domain.Region{},
		&domain.DocType{},
		&domain.PaymentMethod{},
		// &domain.Deposite{},
		&domain.Transaction{},
		// &domain.Plan{},
		&domain.Task{},
		&domain.Customer{},
	)
}
