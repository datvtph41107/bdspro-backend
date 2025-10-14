package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Core entities
type UserProfileEntity struct {
	ProfileID       uint64     `gorm:"primaryKey;column:profile_id" json:"profileId"`
	Phone           string     `gorm:"size:20" json:"phone"`
	Phone2          string     `gorm:"size:20" json:"phone2"`
	Email           string     `gorm:"size:50" json:"email"`
	TaxCode         string     `gorm:"size:13" json:"taxCode"`
	FullName        string     `gorm:"size:255" json:"fullName"`
	Address         string     `gorm:"size:255" json:"address"`
	Gender          uint8      `json:"gender"`
	Avatar          string     `gorm:"size:255" json:"avatar"`
	TickVerified    bool       `gorm:"column:tick_verified" json:"tickVerified"`
	BackgroundImage string     `gorm:"size:255" json:"backgroundImage"`
	ReferralCode    string     `gorm:"size:255" json:"referralCode"`
	ReferralBy      string     `gorm:"size:255" json:"referralBy"`
	RoleType        uint8      `json:"roleType"`
	RoleKey         uint32     `json:"roleKey"`
	Position        string     `gorm:"size:255" json:"position"`
	DepartmentID    uint64     `gorm:"column:department_id" json:"departmentId"`
	Birth           *time.Time `gorm:"column:birth" json:"birth,omitempty"`
	StartDate       *time.Time `gorm:"column:start_date" json:"startDate,omitempty"`
	PlanID          *uint64    `gorm:"column:plan_id" json:"planId,omitempty"`
	PlanAt          *time.Time `gorm:"column:plan_at" json:"planAt,omitempty"`
	Slogan          string     `gorm:"column:slogan" json:"slogan"`
	Website         string     `gorm:"column:website" json:"website"`
	Facebook        string     `gorm:"column:facebook" json:"facebook"`
	Instagram       string     `gorm:"column:instagram" json:"instagram"`
	Twitter         string     `gorm:"column:twitter" json:"twitter"`
	Linkedin        string     `gorm:"column:linkedin" json:"linkedin"`
	Youtube         string     `gorm:"column:youtube" json:"youtube"`
	Introduction    string     `gorm:"column:introduction" json:"introduction"`
	Visibility      uint8      `gorm:"column:visibility" json:"visibility"`
	StatusOnline    uint8      `gorm:"column:status_online" json:"statusOnline"`
	RoleID          *uint64    `gorm:"column:role_id" json:"roleId,omitempty"`
	Status          uint8      `gorm:"column:status" json:"status"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deletedAt,omitempty"`
}

func (UserProfileEntity) TableName() string {
	return "user_profile"
}

type AdminProfile struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID uint64     `gorm:"column:profile_id" json:"profileId"`
	Email     string     `gorm:"size:100" json:"email"`
	FullName  string     `gorm:"size:255" json:"fullName"`
	Phone     string     `gorm:"size:20" json:"phone"`
	Avatar    string     `gorm:"size:255" json:"avatar"`
	Status    uint8      `gorm:"default:10" json:"status"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (AdminProfile) TableName() string {
	return "admin_profile"
}

type Plan struct {
	ID                 uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PackageName        string     `gorm:"column:package_name" json:"packageName"`
	Duration           int64      `gorm:"column:duration" json:"duration"`
	Price              int64      `gorm:"column:price" json:"price"`
	CreateProduct      int        `gorm:"column:create_product" json:"createProduct"`
	LimitPost          int64      `gorm:"column:post_limit" json:"limitPost"`
	PublicProductLimit int        `gorm:"column:public_product_limit" json:"publicProductLimit"`
	CreateProductAI    int        `gorm:"column:create_product_ai" json:"createProductAI"`
	FriendLimit        int        `gorm:"column:friend_limit" json:"friendLimit"`
	CustomerLimit      int        `gorm:"column:customer_limit" json:"customerLimit"`
	JoinOrg            int        `gorm:"column:join_org" json:"joinOrg"`
	CreateOrganization int        `gorm:"column:create_organization" json:"createOrganization"`
	OrgMemberLimit     int        `gorm:"column:org_member_limit" json:"orgMemberLimit"`
	OrgProductLimit    int        `gorm:"column:org_product_limit" json:"orgProductLimit"`
	Color              string     `gorm:"column:color" json:"color"`
	Visibility         bool       `gorm:"column:visibility" json:"visibility"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt          *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (Plan) TableName() string {
	return "plans"
}

type Membership struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID    uint64     `gorm:"column:profile_id" json:"profileId"`
	PlanID       uint64     `gorm:"column:plan_id" json:"planId"`
	FinishedDate *time.Time `gorm:"column:finished_date" json:"finishedDate,omitempty"`
	OrderID      uint64     `gorm:"column:order_id" json:"orderId"`
	PaidAt       *time.Time `gorm:"column:paid_at" json:"paidAt,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (Membership) TableName() string {
	return "memberships"
}

// Social entities
type GroupEntity struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"not null" json:"name"`
	TextColor       string     `gorm:"type:varchar(12)" json:"textColor"`
	BackgroundColor string     `gorm:"type:varchar(12)" json:"backgroundColor"`
	CreatedBy       uint64     `json:"createdBy"`
	UpdatedBy       uint64     `json:"updatedBy"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt       *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (GroupEntity) TableName() string {
	return "db_group"
}

type FriendEntity struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID        uint64     `gorm:"column:sender_id" json:"senderId"`
	ReceiverID      uint64     `gorm:"column:receiver_id" json:"receiverId"`
	GroupID         *uint64    `gorm:"column:group_id" json:"groupId,omitempty"`
	GroupReceiverID *uint64    `gorm:"column:group_receiver_id" json:"groupReceiverId,omitempty"`
	Status          uint8      `json:"status"`
	Text            string     `json:"text"`
	RespondedAt     *time.Time `gorm:"column:responded_at" json:"respondedAt,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt       *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (FriendEntity) TableName() string {
	return "db_friend"
}

type FollowEntity struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowerID uint64     `gorm:"column:follower_id" json:"followerId"`
	FollowedID uint64     `gorm:"column:followed_id" json:"followedId"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	DeletedAt  *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (FollowEntity) TableName() string {
	return "db_follow"
}

type BlockEntity struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BlockerID uint64     `gorm:"column:blocker_id" json:"blockerId"`
	BlockedID uint64     `gorm:"column:blocked_id" json:"blockedId"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (BlockEntity) TableName() string {
	return "db_block"
}

type ContactEntity struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID uint64     `gorm:"column:profile_id" json:"profileId"`
	ContactID uint64     `gorm:"column:contact_id" json:"contactId"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (ContactEntity) TableName() string {
	return "db_contact"
}

type BookmarkUserEntity struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID  uint64     `gorm:"column:profile_id" json:"profileId"`
	BookmarkID uint64     `gorm:"column:bookmark_id" json:"bookmarkId"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	DeletedAt  *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (BookmarkUserEntity) TableName() string {
	return "db_bookmark_user"
}

type ProfessionEntity struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID uint64     `gorm:"column:profile_id" json:"profileId"`
	Name      string     `gorm:"size:255" json:"name"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (ProfessionEntity) TableName() string {
	return "db_profession"
}

type CertificationEntity struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID uint64     `gorm:"column:profile_id" json:"profileId"`
	Name      string     `gorm:"size:255" json:"name"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (CertificationEntity) TableName() string {
	return "db_certification"
}

type SystemConfigEntity struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Key       string     `gorm:"size:255;unique" json:"key"`
	Value     string     `gorm:"type:text" json:"value"`
	GroupKey  string     `gorm:"size:100" json:"groupKey"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (SystemConfigEntity) TableName() string {
	return "system_config"
}

type ExportData struct {
	// Core entities
	UserProfiles  []UserProfileEntity `json:"userProfiles"`
	AdminProfiles []AdminProfile      `json:"adminProfiles"`
	Plans         []Plan              `json:"plans"`
	Memberships   []Membership        `json:"memberships"`

	// Social entities
	Groups         []GroupEntity         `json:"groups"`
	Friends        []FriendEntity        `json:"friends"`
	Follows        []FollowEntity        `json:"follows"`
	Blocks         []BlockEntity         `json:"blocks"`
	Contacts       []ContactEntity       `json:"contacts"`
	BookmarkUsers  []BookmarkUserEntity  `json:"bookmarkUsers"`
	Professions    []ProfessionEntity    `json:"professions"`
	Certifications []CertificationEntity `json:"certifications"`

	// System
	SystemConfigs []SystemConfigEntity `json:"systemConfigs"`

	// Metadata
	ExportedAt       time.Time `json:"exportedAt"`
	TotalUsers       int       `json:"totalUsers"`
	TotalPlans       int       `json:"totalPlans"`
	TotalMemberships int       `json:"totalMemberships"`
	TotalFriends     int       `json:"totalFriends"`
}

func main() {
	// Đọc connection string từ environment hoặc sử dụng default
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default connection string từ config/app.yml
		dsn = "host=14.225.210.29 user=hieubap password=Matkhau.1 dbname=db_user port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"
	}

	// Kết nối database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to database successfully!")

	// Khởi tạo export data
	exportData := ExportData{
		ExportedAt: time.Now(),
	}

	// Fetch User Profiles (không bị xóa)
	fmt.Println("Fetching user profiles...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.UserProfiles).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch user profiles: %v", err)
	} else {
		exportData.TotalUsers = len(exportData.UserProfiles)
		fmt.Printf("✓ Fetched %d user profiles\n", exportData.TotalUsers)
	}

	// Fetch Admin Profiles
	fmt.Println("Fetching admin profiles...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.AdminProfiles).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch admin profiles: %v", err)
	} else {
		fmt.Printf("✓ Fetched %d admin profiles\n", len(exportData.AdminProfiles))
	}

	// Fetch Plans
	fmt.Println("Fetching plans...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.Plans).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch plans: %v", err)
	} else {
		exportData.TotalPlans = len(exportData.Plans)
		fmt.Printf("✓ Fetched %d plans\n", exportData.TotalPlans)
	}

	// Fetch Memberships
	fmt.Println("Fetching memberships...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.Memberships).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch memberships: %v", err)
	} else {
		exportData.TotalMemberships = len(exportData.Memberships)
		fmt.Printf("✓ Fetched %d memberships\n", exportData.TotalMemberships)
	}

	// Fetch Groups
	fmt.Println("Fetching groups...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.Groups).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch groups: %v", err)
	} else {
		fmt.Printf("✓ Fetched %d groups\n", len(exportData.Groups))
	}

	// Fetch Friends
	fmt.Println("Fetching friends...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.Friends).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch friends: %v", err)
	} else {
		exportData.TotalFriends = len(exportData.Friends)
		fmt.Printf("✓ Fetched %d friends\n", exportData.TotalFriends)
	}

	// Fetch Follows
	fmt.Println("Fetching follows...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.Follows).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch follows: %v", err)
	} else {
		fmt.Printf("✓ Fetched %d follows\n", len(exportData.Follows))
	}

	// Fetch Blocks
	fmt.Println("Fetching blocks...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.Blocks).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch blocks: %v", err)
	} else {
		fmt.Printf("✓ Fetched %d blocks\n", len(exportData.Blocks))
	}

	// Fetch Contacts
	fmt.Println("Fetching contacts...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.Contacts).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch contacts: %v", err)
	} else {
		fmt.Printf("✓ Fetched %d contacts\n", len(exportData.Contacts))
	}

	// Fetch Bookmark Users
	fmt.Println("Fetching bookmark users...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.BookmarkUsers).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch bookmark users: %v", err)
	} else {
		fmt.Printf("✓ Fetched %d bookmark users\n", len(exportData.BookmarkUsers))
	}

	// Fetch Professions
	fmt.Println("Fetching professions...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.Professions).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch professions: %v", err)
	} else {
		fmt.Printf("✓ Fetched %d professions\n", len(exportData.Professions))
	}

	// Fetch Certifications
	fmt.Println("Fetching certifications...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.Certifications).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch certifications: %v", err)
	} else {
		fmt.Printf("✓ Fetched %d certifications\n", len(exportData.Certifications))
	}

	// Fetch System Configs
	fmt.Println("Fetching system configs...")
	err = db.Where("deleted_at IS NULL").Find(&exportData.SystemConfigs).Error
	if err != nil {
		log.Printf("Warning: Failed to fetch system configs: %v", err)
	} else {
		fmt.Printf("✓ Fetched %d system configs\n", len(exportData.SystemConfigs))
	}

	// Export ra file JSON
	outputFile := "user_service_export.json"
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("\n✅ Export completed successfully!\n")
	fmt.Printf("📁 Output file: %s\n", outputFile)
	fmt.Printf("📊 Statistics:\n")
	fmt.Printf("   - User Profiles: %d\n", exportData.TotalUsers)
	fmt.Printf("   - Admin Profiles: %d\n", len(exportData.AdminProfiles))
	fmt.Printf("   - Plans: %d\n", exportData.TotalPlans)
	fmt.Printf("   - Memberships: %d\n", exportData.TotalMemberships)
	fmt.Printf("   - Groups: %d\n", len(exportData.Groups))
	fmt.Printf("   - Friends: %d\n", exportData.TotalFriends)
	fmt.Printf("   - Follows: %d\n", len(exportData.Follows))
	fmt.Printf("   - Blocks: %d\n", len(exportData.Blocks))
	fmt.Printf("   - Contacts: %d\n", len(exportData.Contacts))
	fmt.Printf("   - Bookmark Users: %d\n", len(exportData.BookmarkUsers))
	fmt.Printf("   - Professions: %d\n", len(exportData.Professions))
	fmt.Printf("   - Certifications: %d\n", len(exportData.Certifications))
	fmt.Printf("   - System Configs: %d\n", len(exportData.SystemConfigs))
	fmt.Printf("   - Exported at: %s\n", exportData.ExportedAt.Format("2006-01-02 15:04:05"))
}
