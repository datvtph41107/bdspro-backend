package routes

import (
	_jwt "common/jwt"
	"net/http"
	"relay/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetListMessage retrieves messages for a group with pagination and filtering by time.
// @Summary Get chat history
// @Description Retrieve messages from a group chat
// @Tags Chat
// @Accept json
// @Produce json
// @Param group query string true "Group ID"
// @Param from_time query string false "Filter messages from this time (ISO 8601 format: 2006-01-02T15:04:05)"
// @Param size query int true "Number of messages per page"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /chat/message/history [get]
func GetChatHistory(c *gin.Context) {
	fromTimeStr, _ := c.GetQuery("from_time")
	fromTime, err := time.Parse("2006-01-02T15:04:05", fromTimeStr) // Định dạng ISO 8601
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	sizeStr, _ := c.GetQuery("size")
	size, _ := strconv.ParseInt(sizeStr, 10, 32)
	groupIdStr, _ := c.GetQuery("group")
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 32)
	userId := _jwt.GetProfileId(c)
	ok, err := models.CheckUserInGroup(uint64(userId), groupId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Bạn không thuộc nhóm chat"})
		return
	}
	// var webhookData payos.WebhookType
	response, err := models.GetListMessage(groupId, fromTime, int(size))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"code": 0, "data": response})
}

// GetGroupsByMemberID retrieves the list of groups a user belongs to with pagination.
// @Summary Get user's group list
// @Description Retrieve a list of groups the user is a member of
// @Tags Group
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param size query int true "Number of groups per page"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /chat/group/list [get]
func GetGroupList(c *gin.Context) {
	sizeStr, _ := c.GetQuery("size")
	size, _ := strconv.ParseInt(sizeStr, 10, 32)
	pageStr, _ := c.GetQuery("page")
	page, _ := strconv.ParseInt(pageStr, 10, 32)
	// groupIdStr, _ := c.GetQuery("group")
	// groupId, _ := strconv.ParseUint(groupIdStr, 10, 32)
	userId := _jwt.GetProfileId(c)
	groups, err := models.GetGroupsByMemberID(uint64(userId), int(page), int(size))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"code": 0, "data": groups})
}
