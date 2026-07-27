package seopublic

import (
	"net/http"

	"crm/internal/dto"
	"crm/internal/usecase"

	"github.com/gin-gonic/gin"
)

type OperationsHandler struct {
	administrativeUnitProvisioning *usecase.SeoAdministrativeUnitProvisioningUsecase
}

func NewOperationsHandler(
	administrativeUnitProvisioning *usecase.SeoAdministrativeUnitProvisioningUsecase,
) *OperationsHandler {
	return &OperationsHandler{administrativeUnitProvisioning: administrativeUnitProvisioning}
}

func (h *OperationsHandler) ProvisionAdministrativeUnit(c *gin.Context) {
	if h == nil || h.administrativeUnitProvisioning == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    http.StatusServiceUnavailable,
			"message": "Administrative Unit provisioning chưa sẵn sàng",
		})
		return
	}

	var req dto.ProvisionAdministrativeUnitSeoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Provisioning request không hợp lệ",
		})
		return
	}

	response, err := h.administrativeUnitProvisioning.Provision(c.Request.Context(), &req)
	if err != nil {
		writeError(c, err)
		return
	}

	statusCode := http.StatusOK
	if response != nil && response.Created {
		statusCode = http.StatusCreated
	}
	c.Header("Cache-Control", "private, no-store")
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("X-Content-Type-Options", "nosniff")
	c.JSON(statusCode, gin.H{
		"code": statusCode,
		"data": response,
	})
}
