package useraccess

import (
	"testing"
	"time"

	userpb "pb/types/user"
)

func TestFromProtoPreservesMeteringEvidence(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	result, err := fromProto(&userpb.GetAccessResponse{
		SubjectType:     "profile",
		SubjectId:       "42",
		Operation:       "workspace.report.generate",
		Allowed:         true,
		FeatureCode:     "workspace.report.generate",
		MeterCode:       "workspace.report_generation.accepted",
		UnitsPerAction:  1,
		PolicyVersion:   "1.0.0",
		Limit:           100,
		Period:          "subscription_cycle",
		PeriodStartUnix: start.Unix(),
		PeriodEndUnix:   end.Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Metering.MeterCode != "workspace.report_generation.accepted" || result.Metering.UnitsPerAction != 1 {
		t.Fatalf("metering = %+v", result.Metering)
	}
}
