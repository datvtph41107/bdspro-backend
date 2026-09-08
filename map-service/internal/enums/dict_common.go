package enums

// WorkScheduleType represents different types of work schedules
type WorkScheduleType int

const (
	WorkScheduleTypeRegular   WorkScheduleType = iota // Regular business hours
	WorkScheduleTypeHoliday                           // Holiday schedule
	WorkScheduleTypeSpecial                           // Special events
	WorkScheduleTypeEmergency                         // Emergency hours
)

// WorkScheduleStatus represents the status of a work schedule
type WorkScheduleStatus int

const (
	WorkScheduleStatusActive    WorkScheduleStatus = iota // Active schedule
	WorkScheduleStatusInactive                            // Inactive schedule
	WorkScheduleStatusSuspended                           // Temporarily suspended
)

// HolidayType represents different types of holidays
type HolidayType int

const (
	HolidayTypeNational    HolidayType = iota // National holiday
	HolidayTypeReligious                      // Religious holiday
	HolidayTypeCultural                       // Cultural holiday
	HolidayTypeCompany                        // Company holiday
	HolidayTypeMaintenance                    // Maintenance day
	HolidayTypeEmergency                      // Emergency holiday
	HolidayTypeCustom                         // Custom holiday
)

// HolidayStatus represents the status of a holiday
type HolidayStatus int

const (
	HolidayStatusActive    HolidayStatus = iota // Active holiday
	HolidayStatusInactive                       // Inactive holiday
	HolidayStatusExpired                        // Expired holiday
	HolidayStatusCancelled                      // Cancelled holiday
)

// RecurringType represents the type of recurring holiday
type RecurringType int

const (
	RecurringTypeNone    RecurringType = iota // No recurrence
	RecurringTypeYearly                       // Yearly
	RecurringTypeMonthly                      // Monthly
	RecurringTypeWeekly                       // Weekly
	RecurringTypeCustom                       // Custom
)

// ReviewRating represents different rating levels
type ReviewRating int

const (
	ReviewRating1Star ReviewRating = 1
	ReviewRating2Star ReviewRating = 2
	ReviewRating3Star ReviewRating = 3
	ReviewRating4Star ReviewRating = 4
	ReviewRating5Star ReviewRating = 5
)

// ReviewStatus represents the status of a review
type ReviewStatus int

const (
	ReviewStatusPending  ReviewStatus = iota // Pending approval
	ReviewStatusApproved                     // Approved and visible
	ReviewStatusRejected                     // Rejected
	ReviewStatusHidden                       // Hidden by admin
	ReviewStatusDeleted                      // Soft deleted
)

// ReviewReportStatus represents the status of a review report
type ReviewReportStatus int

const (
	ReviewReportStatusPending  ReviewReportStatus = iota // Pending review
	ReviewReportStatusApproved                           // Report approved
	ReviewReportStatusRejected                           // Report rejected
)

// CategoryType represents different types of categories
type CategoryType int

const (
	CategoryTypeMain CategoryType = iota // Main category
	CategoryTypeSub                      // Sub category
	CategoryTypeTag                      // Tag category
)

// CategoryStatus represents the status of a category
type CategoryStatus int

const (
	CategoryStatusActive   CategoryStatus = iota // Active category
	CategoryStatusInactive                       // Inactive category
	CategoryStatusArchived                       // Archived category
)

// AmenityCategory represents different categories of amenities
type AmenityCategory int

const (
	AmenityCategoryParking        AmenityCategory = iota // Parking related
	AmenityCategoryFood                                  // Food and dining
	AmenityCategoryAccessibility                         // Accessibility features
	AmenityCategoryEntertainment                         // Entertainment
	AmenityCategoryServices                              // Services
	AmenityCategorySafety                                // Safety features
	AmenityCategoryTransportation                        // Transportation
	AmenityCategoryOther                                 // Other amenities
)

// AmenityStatus represents the status of an amenity
type AmenityStatus int

const (
	AmenityStatusActive      AmenityStatus = iota // Active amenity
	AmenityStatusInactive                         // Inactive amenity
	AmenityStatusMaintenance                      // Under maintenance
)

// EventType represents different types of events
type EventType int

const (
	EventTypeConcert    EventType = iota // Concert
	EventTypeFestival                    // Festival
	EventTypeMeeting                     // Meeting
	EventTypeWorkshop                    // Workshop
	EventTypeConference                  // Conference
	EventTypeExhibition                  // Exhibition
	EventTypeSports                      // Sports
	EventTypeCultural                    // Cultural
	EventTypeSocial                      // Social
	EventTypeOther                       // Other
)

// EventStatus represents the status of an event
type EventStatus int

const (
	EventStatusDraft     EventStatus = iota // Draft
	EventStatusPublished                    // Published
	EventStatusCancelled                    // Cancelled
	EventStatusCompleted                    // Completed
)

// AttendeeStatus represents the status of an attendee
type AttendeeStatus int

const (
	AttendeeStatusPending   AttendeeStatus = iota // Pending
	AttendeeStatusConfirmed                       // Confirmed
	AttendeeStatusCancelled                       // Cancelled
)

// NotificationType represents different types of notifications
type NotificationType int

const (
	NotificationTypeEvent     NotificationType = iota // Event related
	NotificationTypeReview                            // Review related
	NotificationTypeLocation                          // Location related
	NotificationTypeSystem                            // System related
	NotificationTypeMarketing                         // Marketing related
)

// NotificationStatus represents the status of a notification
type NotificationStatus int

const (
	NotificationStatusUnread   NotificationStatus = iota // Unread
	NotificationStatusRead                               // Read
	NotificationStatusArchived                           // Archived
	NotificationStatusDeleted                            // Deleted
)

// FavoriteType represents different types of favorites
type FavoriteType int

const (
	FavoriteTypeLocation FavoriteType = iota // Location favorite
	FavoriteTypeEvent                        // Event favorite
	FavoriteTypeCategory                     // Category favorite
)

// FavoriteStatus represents the status of a favorite
type FavoriteStatus int

const (
	FavoriteStatusActive   FavoriteStatus = iota // Active favorite
	FavoriteStatusArchived                       // Archived favorite
	FavoriteStatusDeleted                        // Deleted favorite
)

// DayOfWeek represents days of the week
type DayOfWeek int

const (
	DayOfWeekSunday    DayOfWeek = iota // Sunday
	DayOfWeekMonday                     // Monday
	DayOfWeekTuesday                    // Tuesday
	DayOfWeekWednesday                  // Wednesday
	DayOfWeekThursday                   // Thursday
	DayOfWeekFriday                     // Friday
	DayOfWeekSaturday                   // Saturday
)

// LocationType represents different types of locations
type LocationType int

const (
	LocationTypeRestaurant LocationType = iota // Restaurant
	LocationTypeHotel                          // Hotel
	LocationTypeShop                           // Shop
	LocationTypeOffice                         // Office
	LocationTypeHospital                       // Hospital
	LocationTypeSchool                         // School
	LocationTypeBank                           // Bank
	LocationTypeGasStation                     // Gas Station
	LocationTypeParking                        // Parking
	LocationTypeOther                          // Other
)

// MapPointType represents different types of map points
type MapPointType int

const (
	MapPointTypeLocation MapPointType = iota // Location point
	MapPointTypeEvent                        // Event point
	MapPointTypeLandmark                     // Landmark point
	MapPointTypeCustom                       // Custom point
)

// String methods for better readability
func (w WorkScheduleType) String() string {
	switch w {
	case WorkScheduleTypeRegular:
		return "regular"
	case WorkScheduleTypeHoliday:
		return "holiday"
	case WorkScheduleTypeSpecial:
		return "special"
	case WorkScheduleTypeEmergency:
		return "emergency"
	default:
		return "unknown"
	}
}

func (h HolidayType) String() string {
	switch h {
	case HolidayTypeNational:
		return "national"
	case HolidayTypeReligious:
		return "religious"
	case HolidayTypeCultural:
		return "cultural"
	case HolidayTypeCompany:
		return "company"
	case HolidayTypeMaintenance:
		return "maintenance"
	case HolidayTypeEmergency:
		return "emergency"
	case HolidayTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

func (r ReviewRating) String() string {
	switch r {
	case ReviewRating1Star:
		return "1_star"
	case ReviewRating2Star:
		return "2_star"
	case ReviewRating3Star:
		return "3_star"
	case ReviewRating4Star:
		return "4_star"
	case ReviewRating5Star:
		return "5_star"
	default:
		return "unknown"
	}
}

func (e EventType) String() string {
	switch e {
	case EventTypeConcert:
		return "concert"
	case EventTypeFestival:
		return "festival"
	case EventTypeMeeting:
		return "meeting"
	case EventTypeWorkshop:
		return "workshop"
	case EventTypeConference:
		return "conference"
	case EventTypeExhibition:
		return "exhibition"
	case EventTypeSports:
		return "sports"
	case EventTypeCultural:
		return "cultural"
	case EventTypeSocial:
		return "social"
	case EventTypeOther:
		return "other"
	default:
		return "unknown"
	}
}

func (d DayOfWeek) String() string {
	switch d {
	case DayOfWeekSunday:
		return "sunday"
	case DayOfWeekMonday:
		return "monday"
	case DayOfWeekTuesday:
		return "tuesday"
	case DayOfWeekWednesday:
		return "wednesday"
	case DayOfWeekThursday:
		return "thursday"
	case DayOfWeekFriday:
		return "friday"
	case DayOfWeekSaturday:
		return "saturday"
	default:
		return "unknown"
	}
}
