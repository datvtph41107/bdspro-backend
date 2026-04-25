package dto

import (
	"encoding/json"
	tqdpb "pb/types/tqd"

	"gorm.io/datatypes"
)

func ProtoToSubscriptionScopeJSON(pb *tqdpb.SubscriptionScopeMessage) (datatypes.JSON, error) {
	if pb == nil {
		return nil, nil
	}
	data, err := json.Marshal(pb)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(data), nil
}

func SubscriptionScopeJSONToProto(raw datatypes.JSON) (*tqdpb.SubscriptionScopeMessage, error) {
	if raw == nil || len(raw) == 0 {
		return nil, nil
	}
	var pb tqdpb.SubscriptionScopeMessage
	if err := json.Unmarshal(raw, &pb); err != nil {
		return nil, err
	}
	return &pb, nil
}

type TriggerConfigDTO struct {
	MinSeverity         string   `json:"minSeverity,omitempty"`
	DebounceSeconds     int      `json:"debounceSeconds,omitempty"`
	AcceptedChangeTypes []string `json:"acceptedChangeTypes,omitempty"`
}

func TriggerConfigToJSON(cfg *TriggerConfigDTO) (datatypes.JSON, error) {
	if cfg == nil {
		return nil, nil
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(data), nil
}

func TriggerConfigFromJSON(raw datatypes.JSON) (*TriggerConfigDTO, error) {
	if raw == nil || len(raw) == 0 {
		return nil, nil
	}
	var cfg TriggerConfigDTO
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
