package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"mdn4/modules/bdspro/enums"
)

// PublicService owns public catalog/enum flows for modules/bdspro.
type PublicService interface {
	GetEnums(ctx context.Context, timestamp int64) (*bdspropb.GetEnumsResponse, error)
	GetEnumByName(ctx context.Context, name string) (*bdspropb.GetEnumResponse, error)
}

// LocalPublicService serves public enum catalogs from in-process maps + SyncUtil.
// @bind: mdn4/modules/bdspro/services.PublicService
type LocalPublicService struct {
	sync      *_utils.SyncUtil
	startTime int64
}

func NewLocalPublicService(sync *_utils.SyncUtil) *LocalPublicService {
	return &LocalPublicService{
		sync:      sync,
		startTime: time.Now().UnixMilli(),
	}
}

func (s *LocalPublicService) GetEnums(ctx context.Context, timestamp int64) (*bdspropb.GetEnumsResponse, error) {
	catalogTS := s.startTime
	if s.sync != nil {
		catalogTS = s.sync.EnsureTimestamp(ctx, _utils.SyncKeyBdsproEnums)
	}

	updated := true
	if s.sync != nil {
		updated = s.sync.HasUpdated(ctx, _utils.SyncKeyBdsproEnums, timestamp)
	} else if timestamp > 0 && timestamp >= catalogTS {
		updated = false
	}
	if !updated {
		return &bdspropb.GetEnumsResponse{
			Timestamp: catalogTS,
			Lastest:   false,
		}, nil
	}

	result := make([]*bdspropb.GetEnumResponse, 0, len(enumProviders))
	for k, v := range enumProviders {
		result = append(result, &bdspropb.GetEnumResponse{
			Key:  k,
			Data: v,
		})
	}
	return &bdspropb.GetEnumsResponse{
		Enums:     result,
		Timestamp: catalogTS,
		Lastest:   true,
	}, nil
}

func (s *LocalPublicService) GetEnumByName(_ context.Context, name string) (*bdspropb.GetEnumResponse, error) {
	if values, ok := enumProviders[name]; ok {
		return &bdspropb.GetEnumResponse{Data: values}, nil
	}
	return nil, _errors.NotFoundException(fmt.Sprintf("enum not found: %s", name))
}

var enumProviders = map[string][]*bdspropb.EnumItem{
	"product_sale_status":     enumMapToList(enums.EProductSaleStatusNames),
	"product_rent_status":     enumMapToList(enums.EProductRentStatusNames),
	"asset_history":           enumMapToList(enums.AssetHistoryNames),
	"target_type":             enumMapToList(enums.TargetTypeNames),
	"visibility":              enumMapToList(enums.EVisibilityNames),
	"producthistory":          enumMapToList(enums.EProductHistoryNames),
	"poststatus":              enumMapToList(enums.EPostStatusNames),
	"doctype":                 enumMapToList(enums.EDocTypeNames),
	"assetstatus":             enumMapToList(enums.EAssetStatusNames),
	"house_orient":            enumMapToList(enums.EHouseOrientNames),
	"house_certificate":       enumMapToList(enums.EHouseCertificateNames),
	"property_status":         enumMapToList(enums.EPropertyStatusNames),
	"building_type":           enumMapToList(enums.BuildingTypeNames),
	"transaction_status":      enumMapToList(enums.TransactionStatusNames),
	"transaction_type":        enumMapToList(enums.TransactionTypeNames),
	"owner_type":              enumMapToList(enums.OwnerTypeNames),
	"post_transaction_status": enumMapToList(enums.EPostTransactionStatusNames),
	"distribution_status":     enumMapToList(enums.DistributionStatusNames),
	"distribution_reason":     enumMapToList(enums.DistributionReasonNameMap),
	"deal_type":               enumMapToList(enums.DealTypeNames),
	"source_status":           enumMapToList(enums.SourceStatusNames),
	"source_type":             enumMapToList(enums.SourceTypeMap),
	"mining_mode":             enumMapToList(enums.MiningModeNames),
	"mining_scope":            enumMapToList(enums.MiningScopeNames),
	"classify_property":       enumMapToList(enums.EClassifyPropertyNames),
	"purpose_used":            enumMapToList(enums.EPurposeUsedNames),
}

func enumMapToList[T ~uint32 | ~uint | ~int](enumMap map[T]string) []*bdspropb.EnumItem {
	keys := make([]T, 0, len(enumMap))
	for k := range enumMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	result := make([]*bdspropb.EnumItem, 0, len(keys))
	for _, k := range keys {
		result = append(result, &bdspropb.EnumItem{
			Value: uint32(k),
			Name:  enumMap[k],
		})
	}
	return result
}
