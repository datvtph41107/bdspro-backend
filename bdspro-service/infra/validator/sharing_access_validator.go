package validator

import (
	"errors"
	pb_bdspro "pb/types/bdspro"
)

type SharingAccessValidator struct {
}

func NewSharingAccessValidator() *SharingAccessValidator {
	return &SharingAccessValidator{}
}

func (v *SharingAccessValidator) ValidateBulkSave(dto *pb_bdspro.SharingAccessBulkRequest) error {
	if dto.DomainId == 0 || dto.DomainType == 0 {
		return errors.New("domainId and domainType are required")
	}

	// if dto.FromType == 0 {
	// 	return errors.New("fromType is required")
	// }

	if dto.Datas != nil {
		for _, data := range dto.Datas {
			if data.ToType == 0 ||
				data.ToId == 0 {
				return errors.New("datas.toType, datas.toId is required")
			}
		}
	}

	return nil
}

func (v *SharingAccessValidator) ValidateSearchAsset(dto *pb_bdspro.SharingAccessRequest) error {
	if dto.DomainId == 0 {
		return errors.New("domainID is required")
	}

	return nil
}

func (v *SharingAccessValidator) ValidateSetCommission(dto *pb_bdspro.SharingAccessByDomainRequest) error {
	if dto.DomainId == 0 || dto.Commission < 0 {
		return errors.New("domainID and commission are required")
	}

	return nil
}
