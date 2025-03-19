package services

import (
	_db "common/db"
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"noti/data"
	"noti/models"
	"noti/repo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationService struct {
	Repo *repo.NotificationRepository
}

func NewNotificationService(repo *repo.NotificationRepository) *NotificationService {
	return &NotificationService{Repo: repo}
}

// CreateNotification: Tạo thông báo mới
func (s *NotificationService) CreateNotification(c *gin.Context, request data.NewNotiRequest) (*data.NotificationResponse, error) {
	profileId := _jwt.GetProfileId(c)
	if profileId == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Thông tin không chính xác",
		}
	}

	notification := models.NotificationEntity{
		UserID:  request.UserID,
		Title:   request.Title,
		Message: request.Message,
		Link:    request.Link,
		Type:    request.Type,
	}

	if err := _db.SaveWithAudit(c, &notification); err != nil {
		return nil, err
	}

	response := s.mapToDTO(notification)
	return &response, nil
}

// Search: Tìm kiếm danh sách thông báo
func (s *NotificationService) Search(c *gin.Context, dto data.SearchNotiRequest, pageable _utils.Page) ([]data.NotificationResponse, error) {
	profileId := _jwt.GetProfileId(c)
	dto.UserID = &profileId
	notifications, err := s.Repo.Search(dto, pageable)
	if err != nil {
		return nil, err
	}

	var responses []data.NotificationResponse
	for _, noti := range notifications {
		responses = append(responses, s.mapToDTO(noti))
	}
	return responses, nil
}

// Read: Đánh dấu thông báo là đã đọc
func (s *NotificationService) Read(c *gin.Context, id uint64) (int64, error) {
	profileId := _jwt.GetProfileId(c)
	return s.Repo.MarkNotificationAsRead(id, profileId)
}

// Count: Đếm số thông báo chưa đọc
func (s *NotificationService) Count(c *gin.Context) (int64, error) {
	profileId := _jwt.GetProfileId(c)
	return s.Repo.Count(profileId)
}

// mapToDTO: Chuyển đổi NotificationEntity thành NotificationResponse
func (s *NotificationService) mapToDTO(entity models.NotificationEntity) data.NotificationResponse {
	return data.NotificationResponse{
		ID:        entity.ID,
		UserID:    entity.UserID,
		Title:     entity.Title,
		Message:   entity.Message,
		Link:      entity.Link,
		Type:      entity.Type,
		IsRead:    entity.IsRead,
		Timestamp: entity.CreatedAt,
	}
}

func (s *NotificationService) Delete(c *gin.Context) (int, error) {
	profileId := _jwt.GetProfileId(c)
	idParam := c.Param("id")
	notiId, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, err

	}
	return 1, s.Repo.DeleteByID(profileId, notiId)
}
