package services

import (
	_routes "common/routes"
	"errors"
	"fmt"
	"map/models"
	"map/repo"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type MapService struct {
	MapRepo *repo.MapRepo
}

func NewMapService(repo *repo.MapRepo) *MapService {
	return &MapService{MapRepo: repo}
}

func (s *MapService) FindNearbyLocations(lat, lng, radius float64) ([]models.MapPoint, error) {
	return s.MapRepo.FindLocationNearby(lat, lng, radius)
}

func (s *MapService) NewLocation(dto *models.MapPointDTO) (*models.MapPointResponse, error) {
	err := s.MapRepo.CreateNewPoint(dto.Lat, dto.Lng, dto.Name)
	if err != nil {
		return nil, err
	}
	return &models.MapPointResponse{
		Lat:  dto.Lat,
		Lng:  dto.Lng,
		Name: dto.Name,
	}, nil
}

func (s *MapService) findLocationByPolygon(coordinates []models.MapPointDTO) ([]models.MapPointDTO, error) {
	if len(coordinates) < 3 {
		return nil, errors.New("Polygon must have at least 3 points.")
	}

	// Chuyển danh sách tọa độ thành chuỗi POLYGON((x1 y1, x2 y2, ..., x1 y1))
	var coords []string
	for _, coord := range coordinates {
		coords = append(coords, fmt.Sprintf("%f %f", coord.Lng, coord.Lat))
	}
	coords = append(coords, fmt.Sprintf("%f %f", coordinates[0].Lng, coordinates[0].Lat)) // Đóng polygon
	polygon := fmt.Sprintf("POLYGON((%s))", strings.Join(coords, ", "))

	fmt.Println(polygon)
	return s.MapRepo.FindLocationInPolygon(polygon)
}

// @Summary Tìm địa điểm theo vùng (polygon)
// @Description Tìm các điểm nằm trong một vùng đa giác
// @Tags Locations
// @Accept json
// @Produce json
// @Param polygon body []models.MapPointDTO true "Danh sách tọa độ tạo thành vùng đa giác"
// @Success 200 {object} []models.MapPointDTO
// @Router /locations/find-by-polygon [post]
func (service *MapService) HandleNearbyLocations(c *gin.Context) {
	keys := []string{"lat", "lng", "radius"}
	values := make([]float64, len(keys))

	for i, key := range keys {
		valueStr := c.Query(key)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid " + key + " value"})
			return
		}
		values[i] = value
	}

	locations, err := service.FindNearbyLocations(values[0], values[1], values[2])
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, locations)
}

// @Summary Thêm địa điểm mới
// @Description Thêm một điểm mới vào cơ sở dữ liệu
// @Tags Locations
// @Accept json
// @Produce json
// @Param location body models.MapPointDTO true "Thông tin địa điểm"
// @Router /locations/new [post]
// @Security BearerAuth
func (service *MapService) HandleNewLocation(c *gin.Context) {
	var dto models.MapPointDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	location, err := service.NewLocation(&dto)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, location)
}

// @Summary Tìm các địa điểm lân cận
// @Description Trả về danh sách các điểm gần tọa độ đã cho trong phạm vi bán kính
// @Tags Locations
// @Accept json
// @Produce json
// @Param lat query float64 true "Vĩ độ"
// @Param lng query float64 true "Kinh độ"
// @Param radius query float64 false "Bán kính tìm kiếm (mặc định: 5km)" default(5)
// @Router /locations/nearby [get]
// @Security BearerAuth
func (service *MapService) HandleFindByPolygon(c *gin.Context) {
	var polygon []models.MapPointDTO
	if err := c.ShouldBindJSON(&polygon); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	locations, err := service.findLocationByPolygon(polygon)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	_routes.RouteResult(c, locations, nil)
}
