package db

import (
	"bdspro/internal/domain"
	_db "common/db"
)

func MigrateDomain() {
	_db.DB.AutoMigrate(
		// &domain.AreaRegion{},
		// &domain.PropertyAreaRegion{},
		// Core entities
		// &domain.Observation{},
		// &domain.Asset{},
		// &domain.Product{},
		// &domain.Post{},
		// &domain.Intent{},
		// &domain.Observation{},
		// &domain.Property{},
		// &domain.Block{},
		// &domain.Building{},
		// &domain.DistrictV2{},
		&domain.PropertyImpactVersion{},
		&domain.PropertyAuditLog{},
		// &domain.ProductStats{},

		// // Asset related
		// &domain.AssetCost{},
		// &domain.AssetCostType{},
		// &domain.AssetIncomeType{},
		// &domain.AssetExploitation{},
		// &domain.AssetLegal{},
		// &domain.AssetOrganization{},
		// &domain.AssetUser{},
		// &domain.AssetShare{},

		// // Product related
		// &domain.ProductUser{},
		// &domain.ProductOrganization{},
		// &domain.ProductHistory{},
		// &domain.ProductPrice{},
		// &domain.ProductMedia{},
		// &domain.ProductPrivate{},
		// &domain.ProductAsset{},

		// // Post related
		// &domain.PostMediaEntity{},
		// &domain.PostOrganization{},
		// &domain.PostUser{},

		// // Document related
		// &domain.IncomeDocument{},
		// &domain.AttachDocument{},

		// // Project related
		// &domain.Project{},
		// &domain.ProjectBuild{},
		// &domain.ProjectInventory{},
		// &domain.BuildRoom{},

		// // Apartment related
		// &domain.Apartment{},
		// &domain.ApartmentAttribute{},

		// Deal related
		// &domain.Deal{},
		// &domain.DealMember{},
		// &domain.DealMilestone{},
		// &domain.DealProduct{},
		// &domain.DealInvestment{},
		// &domain.DealOfOrganization{},
		// &domain.DealOfGroup{},
		// &domain.DealOfBranch{},
		// &domain.AttachDocument{},
		// &domain.DealInternalNote{},

		// Transaction related
		// &domain.Transaction{},
		// &domain.Tx{},
		// &domain.TxTransaction{},
		// &domain.TxTransactionHistory{},
		// &domain.TxTransactionProduct{},
		// &domain.TxTransactionGroup{},
		// &domain.TxTransactionCategory{},
		// &domain.TxAction{},
		// &domain.TxContractDeal{},
		// &domain.TxDealCost{},
		// &domain.TxCostType{},

		// Sharing and Access
		// &domain.SharingAccess{},
		// &domain.Sharing{},
		// &domain.RecordHistory{},

		// // Location entities
		// &domain.Province{},
		&domain.ProvinceV2{},
		// &domain.District{},
		// &domain.Ward{},
		&domain.WardV2{},
		// &domain.Region{},

		// // Other entities
		// &domain.Amenity{},
		// &domain.Developer{},
		// &domain.PropertyType{},
		// &domain.DocType{},
		// &domain.PaymentMethod{},
		// &domain.Task{},
		// &domain.Customer{},
		// &domain.BDSDomain{},
		// &domain.HouseInfo{},
		// &domain.Action{},

		// Property related
		// &domain.Product{},
		// &domain.Property{},
		// &domain.PropertyLandInfo{},
		// &domain.PropertyBuildingInfo{},
		// &domain.PropertyMedia{},
		// &domain.PropertyTagLink{},
		// &domain.PropertyRelation{},
		// &domain.Tag{},
		// &domain.MemberPlan{},
		// &domain.ProfileTransfer{},
		// &domain.BankAccount{},
		// &domain.Color{},
		// &domain.CommissionStats{},

		// // Attribute related
		// &domain.Attribute{},
		// &domain.AttributeValue{},
		// &domain.AttributeProduct{},
		// &domain.DistributionEntity{},
		// &domain.ProductUser{},
		// &domain.PropertyAmenity{},
		// &domain.PropertyRelation{},
		// &domain.ProvinceV2{},
		// &domain.WardV2{},
		// &domain.ProductNote{},
		// &domain.ProductNoteMention{},
		// &domain.ProductNoteFile{},

		// &domain.ProductUser{},
		// &domain.ProductNote{},
		// &domain.DistributionEntity{},
		// &domain.ProductPrice{},
		// &domain.Product{},
		// &domain.Property{},
		// &domain.PropertyLandInfo{},
		// &domain.PropertyMedia{},
		// &domain.PropertyUser{},
		// &domain.PropertyLineage{},
		// &domain.PropertyIdentify{},
		// &domain.PropertyInfo{},
		// &domain.PropertyLocation{},
		// &domain.PropertyType{},
		// &domain.PropertyLandInfo{},
		// &domain.PropertyBuildingInfo{},
		// &domain.PropertyMedia{},
		// &domain.PropertyExternalRef{},
		// &domain.PropertyEdvidence{},
		// &domain.PropertyTagLink{},
		// &domain.PropertyRelation{},
		// &domain.PropertyStatistic{},
		// &domain.PropertyLineage{},
	)
}
