package enums

type ESettingTab int

const (
	ETabProfile       ESettingTab = 10
	ETabPost          ESettingTab = 20
	ETabActivity      ESettingTab = 30
	ETabSetting       ESettingTab = 40
	EIntrodution      ESettingTab = 50
	EEvaluation       ESettingTab = 60
	EReaction         ESettingTab = 70
	EShare            ESettingTab = 80
	ESharePost        ESettingTab = 90
	EShareActivity    ESettingTab = 100
	EShareSetting     ESettingTab = 110
	EShareIntrodution ESettingTab = 120
	EShareEvaluation  ESettingTab = 130
)

var SettingTabMap = map[ESettingTab]string{
	ETabProfile:    "Hồ sơ",
	ETabPost:       "Bài viết",
	ETabActivity:   "Hoạt động",
	ETabSetting:    "Cài đặt",
	EIntrodution:   "Giới thiệu",
	EEvaluation:    "Đánh giá",
	EReaction:      "Phản hồi",
	EShare:         "Chia sẻ",
	ESharePost:     "Chia sẻ bài viết",
	EShareActivity: "Chia sẻ hoạt động",
}
