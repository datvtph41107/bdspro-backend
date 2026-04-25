package mapper

import (
	"fmt"
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
				Raw: string(p.Geometry),
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
		Id: p.ID,
		// Name:      p.Name,
		Latitude:  p.Lat,
		Longitude: p.Lng,

		// MinLat: p.MinLat,
		// MinLng: p.MinLng,
		// MaxLat: p.MaxLat,
		// MaxLng: p.MaxLng,
	}
	if val, err := p.Geometry.Value(); err == nil && val != nil {
		switch v2 := val.(type) {
		case []byte:
			v.Geometry = v2
		case string:
			v.Geometry = []byte(v2)
		default:
			// fallback (hiếm)
			v.Geometry = []byte(fmt.Sprintf("%v", v2))
		}
	}

	return v
}
