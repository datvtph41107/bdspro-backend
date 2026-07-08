package utils

import (
	"crm/internal/enums"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOwnerFromContext(c *gin.Context) (uint64, enums.EOwnerOf, error) {
	ownerId, err := strconv.ParseUint(c.Param("ownerId"), 10, 64)
	if err != nil {
		return 0, enums.EOwnerOf(0), err
	}
	ownerTypeInt, err := strconv.ParseInt(c.Param("ownerType"), 10, 64)
	if err != nil {
		return 0, enums.EOwnerOf(0), err
	}
	ownerType := enums.EOwnerOf(ownerTypeInt)
	return ownerId, ownerType, nil
}