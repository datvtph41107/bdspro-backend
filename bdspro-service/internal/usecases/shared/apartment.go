package shared_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	_routes "common/routes"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ApartmentUsecase struct {
	Repo          repo.ApartmentRepo
	AttributeRepo repo.ApartmentAttrRepo
}

func NewApartmentUsecase(apartmentRepo repo.ApartmentRepo) *ApartmentUsecase {
	return &ApartmentUsecase{Repo: apartmentRepo}
}

func (s *ApartmentUsecase) UpdateApartments(c *gin.Context) (*[]domain.Apartment, error) {
	var body dto.ApartmentRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Dữ liệu không hợp lệ",
		}
	}

	attribute := body.Attribute
	attribute.BuildID = body.BuildID
	if err := s.AttributeRepo.Create(c, &attribute); err != nil {
		return nil, err
	}

	apartments := body.Elements
	for idx := range apartments {
		apartments[idx].BuildID = body.BuildID
	}
	if err := s.Repo.UpdateApartments(c, apartments, attribute.ID); err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return &apartments, nil
}

func (s *ApartmentUsecase) UpdateStatusApartments(c *gin.Context) (*[]domain.Apartment, error) {
	var body dto.ApartmentStatusRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Dữ liệu không hợp lệ",
		}
	}

	apartments := body.Apartments
	for idx := range apartments {
		apartments[idx].BuildID = body.BuildID
	}
	if err := s.Repo.UpdateStatusApartments(c, apartments, body.Data.Status, body.Data.Archived); err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return &apartments, nil

}

func (s *ApartmentUsecase) ImportApartments(c *gin.Context) (*[]domain.Apartment, error) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}

	buildIdStr := c.PostForm("buildId")
	buildId, err := strconv.Atoi(buildIdStr)

	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}
	defer file.Close()

	// Đọc file Excel
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("lỗi đọc file Excel: %w", err)
	}

	// Giả sử dữ liệu nằm trong sheet "Sheet1"
	sheetName := "Sheet1"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("lỗi đọc sheet: %w", err)
	}

	numOrdinal := len(rows[0]) - 2
	for idxOrdinal := 0; idxOrdinal < numOrdinal; idxOrdinal++ {
		var apartments []domain.Apartment
		// attribute := &domain.ApartmentAttribute{
		// 	Area: rows[1][2],
		// }
		column := idxOrdinal + 2
		attribute := domain.ApartmentAttribute{
			Area:        new(float64),
			NumBathroom: new(uint32),
			NumBedroom:  new(uint32),
		}
		fmt.Sscanf(rows[1][column], "%d", attribute.Area)
		fmt.Sscanf(rows[2][column], "%d", attribute.NumBedroom)
		fmt.Sscanf(rows[3][column], "%d", attribute.NumBathroom)
		attribute.Furniture = strings.TrimSpace(rows[4][column])

		fmt.Println(rows[2][column], rows[3][column], rows[4][column], "RESULT")

		attribute.BuildID = uint64(buildId)
		if err := s.AttributeRepo.Create(c, &attribute); err != nil {
			return nil, err
		}
		// Bắt đầu từ hàng thứ 2 để bỏ qua tiêu đề
		for i, row := range rows[5:] {
			// if len(row) < 5 {
			// 	log.Printf("Row %d bị thiếu dữ liệu, bỏ qua\n", i+2)
			// 	continue
			// }

			var apt domain.Apartment
			// _, err := fmt.Sscanf(row[0], "%d", &apt.BuildID)
			// fmt.Println(row)
			// fmt.Println(i, column)

			// _, err = fmt.Sscanf(row[1], "%d", &apt.Floor)
			if column < len(row) {
				_, err = fmt.Sscanf(row[column], "%d", &apt.Status)
			}
			slog.Default().Debug(
				"apartment import row parsed",
				slog.Int("column", column),
				slog.Int("row_length", len(row)),
				slog.Uint64("status", uint64(apt.Status)),
			)

			// _, err = fmt.Sscanf(row[3], "%d", &apt.Status)
			// _, err = fmt.Sscanf(row[4], "%d", &apt.Archived)
			floor := (i + 1)
			apt.Floor = &floor
			apt.BuildID = uint64(buildId)
			apt.Ordinal = idxOrdinal + 1

			if err != nil {
				slog.Default().Warn(
					"apartment import row parse failed",
					slog.Int("row", i+2),
					slog.Any("error", err),
				)
				// continue
			}

			apartments = append(apartments, apt)
		}
		s.Repo.UpdateApartment2(c, apartments, attribute.ID)
	}

	return nil, nil
}
