package mapper

import (
	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
)

// ParcelMapper maps between domain and DTO for POI
type ParcelMapper struct{}

func NewParcelMapper() *ParcelMapper {
	return &ParcelMapper{}
}

func (m *ParcelMapper) ProtoToDomains(req *tqdpb.ParcelListDTO) []qh_domain.Parcel {
	parcels := make([]qh_domain.Parcel, 0, len(req.Data))

	for _, p := range req.Data {

		// var geom postgis.Geometry // 👈 interface type

		// 👉 chỉ assign raw string, KHÔNG parse
		// if p.Geometry != nil {
		// 	geom =
		// }

		parcels = append(parcels, qh_domain.Parcel{
			// Name:     p.Name,
			Lat: p.Latitude,
			Lng: p.Longitude,
			Geometry: qh_domain.RawGeometry{
				Raw: p.Geometry,
			},
		})
	}

	return parcels
}

func (m *ParcelMapper) ToProtoFromDomainList(parcels *dto.ParcelList) []*tqdpb.Parcel {
	var protoParcels []*tqdpb.Parcel

	for _, p := range parcels.Data {
		v := m.DomainToProto(&p)
		protoParcels = append(protoParcels, v)
	}

	return protoParcels
}

func (m *ParcelMapper) DomainToProto(p *qh_domain.Parcel) *tqdpb.Parcel {
	v := &tqdpb.Parcel{
		Id:        p.ID,
		Latitude:  p.Lat,
		Longitude: p.Lng,
		// RefId:     p.RefID,
		// RefType:   p.RefType,
		Geometry: p.Geometry.Raw, // trực tiếp, không cần switch
	}
	return v
}

func (m *ParcelMapper) ParcelInfoToProto(p *qh_domain.QHParcelInfo) *tqdpb.ParcelInfo {

	return &tqdpb.ParcelInfo{
		Id:         p.ID,
		Lat:        p.Latitude,
		Lng:        p.Longitude,
		AddressStr: p.AddressStr,
		// MapNumber:  p.MapNumber,
		Geometry:   p.Geometry.Raw,
		TotalArea:  p.TotalArea,
		MapNumber:  p.MapNumber,
		LandNumber: p.LandNumber,
	}
}
