package _category

//	type ISearchDTO interface {
//		GetPage() int
//		SetPage(int)
//		GetSize() int
//		SetSize(int)
//	}
type SearchDTO[DTO any] struct {
	Page int `form:"page"`
	Size int `form:"size"`
	DTO  DTO
}

// func (dto SearchDTO[DTO]) GetPage() int {
// 	return dto.Page
// }
// func (dto SearchDTO[DTO]) SetPage(p int) {
// 	dto.Page = p
// }

// func (dto SearchDTO[DTO]) GetSize() int {
// 	if dto.Size == 0 {
// 		return 20
// 	}
// 	return dto.Size
// }
// func (dto SearchDTO[DTO]) SetSize(s int) {
// 	dto.Size = s
// }
