package mapper

import (
	"encoding/json"
	"fmt"
	qh_domain "tqd/internal/domain/qh"

	pb "pb/types/tqd"
)

func ToProtoOneHouse(oh *qh_domain.OneHouse) (*pb.OneHouse, error) {
	if oh == nil {
		return nil, nil
	}

	planningInfos := make([]*pb.PlanningInfo, len(oh.PlanningInfos))
	for i, info := range oh.PlanningInfos {
		planningInfos[i] = &pb.PlanningInfo{
			PlanningTypeCode: info.PlanningTypeCode,
			PlanningTypeName: info.PlanningTypeName,
			// PlanningArea:     info.PlanningArea,
		}
	}

	var geometry *pb.Geometry
	if oh.GeomGeoJSON != "" {
		var geoJSON struct {
			Type        string          `json:"type"`
			Coordinates json.RawMessage `json:"coordinates"`
		}
		if err := json.Unmarshal([]byte(oh.GeomGeoJSON), &geoJSON); err != nil {
			return nil, fmt.Errorf("parse GeoJSON failed: %w", err)
		}
		if geoJSON.Type != "MultiPolygon" {
			return nil, fmt.Errorf("unsupported geometry type: %s", geoJSON.Type)
		}

		var multiCoords [][][][]float64
		if err := json.Unmarshal(geoJSON.Coordinates, &multiCoords); err != nil {
			return nil, fmt.Errorf("unmarshal coordinates failed: %w", err)
		}

		geometry = &pb.Geometry{
			MultiPolygon: make([]*pb.MultiPolygon, len(multiCoords)),
		}
		for i, polygons := range multiCoords {
			rings := make([]*pb.LineString, len(polygons))
			for j, ringCoords := range polygons {
				points := make([]*pb.PointTile, len(ringCoords))
				for k, coord := range ringCoords {
					if len(coord) < 2 {
						continue
					}
					points[k] = &pb.PointTile{X: coord[0], Y: coord[1]}
				}
				rings[j] = &pb.LineString{Points: points}
			}
			geometry.MultiPolygon[i] = &pb.MultiPolygon{
				Polygons: []*pb.Polygon{{Rings: rings}},
			}
		}
	}

	return &pb.OneHouse{
		Id:                   oh.ID,
		Lat:                  oh.Lat,
		Lng:                  oh.Lng,
		Address:              oh.Address,
		District:             oh.District,
		DistrictCode:         oh.DistrictCode,
		City:                 oh.City,
		CityCode:             oh.CityCode,
		Province:             oh.Province,
		ProvinceCode:         oh.ProvinceCode,
		Ward:                 oh.Ward,
		WardCode:             oh.WardCode,
		MapNumber:            oh.MapNumber,
		LandNumber:           oh.LandNumber,
		Latitude:             oh.Latitude,
		Longitude:            oh.Longitude,
		Direction:            oh.Direction,
		Shape:                oh.Shape,
		NumberOfFacades:      oh.NumberOfFacades,
		TotalArea:            oh.TotalArea,
		PlanningLandType:     oh.PlanningLandType,
		PlanningLandTypeCode: oh.PlanningLandTypeCode,
		PlanningInfos:        planningInfos,
		ShapeFileId:          oh.ShapeFileID,
		Geometry:             geometry,
		IsSavedByUser:        oh.IsSavedByUser,
		PropertyCode:         oh.PropertyCode,
		PropertyUuid:         oh.PropertyUUID,
		AcceptedTnc:          oh.AcceptedTnc,
	}, nil
}
