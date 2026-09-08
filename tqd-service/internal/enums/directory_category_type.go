package enums

// DirectoryCategoryType represents different types of directory categories
type DirectoryCategoryType uint32

const (
	DirectoryCategoryTypeMain     DirectoryCategoryType = iota // Main category
	DirectoryCategoryTypeSub                                   // Sub category
	DirectoryCategoryTypeTag                                   // Tag category
	DirectoryCategoryTypeLocation                              // Location-based category
)
