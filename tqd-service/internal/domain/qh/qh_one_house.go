package qh_domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type OneHouse struct {
	ID                   string        `gorm:"column:id;primaryKey"`
	Lat                  float64       `gorm:"column:lat"`
	Lng                  float64       `gorm:"column:lng"`
	Address              string        `gorm:"column:address"`
	District             string        `gorm:"column:district"`
	DistrictCode         string        `gorm:"column:district_code"`
	City                 string        `gorm:"column:city"`
	CityCode             string        `gorm:"column:city_code"`
	Province             string        `gorm:"column:province"`
	ProvinceCode         string        `gorm:"column:province_code"`
	Ward                 string        `gorm:"column:ward"`
	WardCode             string        `gorm:"column:ward_code"`
	MapNumber            string        `gorm:"column:map_number"`
	LandNumber           string        `gorm:"column:land_number"`
	Latitude             float64       `gorm:"column:latitude"`
	Longitude            float64       `gorm:"column:longitude"`
	Direction            string        `gorm:"column:direction"`
	Shape                string        `gorm:"column:shape"`
	NumberOfFacades      string        `gorm:"column:number_of_facades"`
	TotalArea            float64       `gorm:"column:total_area"`
	PlanningLandType     string        `gorm:"column:planning_land_type"`
	PlanningLandTypeCode string        `gorm:"column:planning_land_type_code"`
	PlanningInfos        PlanningInfos `gorm:"column:planning_infos;type:jsonb"`
	ShapeFileID          string        `gorm:"column:shape_file_id"`
	GeomGeoJSON          string        `gorm:"-"` // không lưu, chỉ nhận từ ST_AsGeoJSON
	IsSavedByUser        bool          `gorm:"column:is_saved_by_user"`
	PropertyCode         string        `gorm:"column:property_code"`
	PropertyUUID         string        `gorm:"column:property_uuid"`
	AcceptedTnc          bool          `gorm:"column:accepted_tnc"`
}

type PlanningInfos []QHPlanningUse

func (p PlanningInfos) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PlanningInfos) Scan(value interface{}) error {
	if value == nil {
		*p = PlanningInfos{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan PlanningInfos: not []byte")
	}
	return json.Unmarshal(bytes, p)
}
