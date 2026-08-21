// common/quota/mock/scenarios.go
package mock

type TestScenario struct {
	Name                string
	Service             string
	Action              string
	UserID              string
	InputID             string // parcel id hoặc user id
	ExpectError         bool
	ExpectQuotaDecrease bool
}

var Scenarios = []TestScenario{
	{
		Name:                "Normal flow - get existing parcel",
		Service:             "tqd",
		Action:              "get_parcel",
		UserID:              "user_normal",
		InputID:             "PARCEL_001",
		ExpectError:         false,
		ExpectQuotaDecrease: true,
	},
	{
		Name:                "Parcel not found",
		Service:             "tqd",
		Action:              "get_parcel",
		UserID:              "user_normal",
		InputID:             "INVALID",
		ExpectError:         true, // 404
		ExpectQuotaDecrease: false,
	},
	{
		Name:                "Database error",
		Service:             "tqd",
		Action:              "get_parcel",
		UserID:              "user_normal",
		InputID:             "DB_ERROR",
		ExpectError:         true,
		ExpectQuotaDecrease: false,
	},
	{
		Name:                "Quota exhausted",
		Service:             "tqd",
		Action:              "get_parcel",
		UserID:              "user_quota_test", // limit=5
		InputID:             "PARCEL_001",
		ExpectError:         false,
		ExpectQuotaDecrease: true,
		// Sẽ chạy 6 lần, lần cuối fail
	},
}
