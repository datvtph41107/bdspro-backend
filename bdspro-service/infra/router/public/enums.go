package public_router

import (
	"bdspro/internal/enums"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func EnumMapToList[T ~uint32 | ~uint | ~int](enumMap map[T]string) []map[string]interface{} {
	keys := make([]T, 0, len(enumMap))
	for k := range enumMap {
		keys = append(keys, k)
	}

	// Sắp xếp keys theo giá trị tăng dần
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	result := make([]map[string]interface{}, 0, len(keys))
	for _, k := range keys {
		result = append(result, map[string]interface{}{
			"name":  enumMap[k],
			"value": k,
		})
	}
	return result
}

// @Summary Lấy enum dạng name-value
// @Tags Enum
// @Param type path string true "Tên loại enum (VD: saleStatus)"
// @Produce json
// @Success 200 {object} []map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /enums/{type} [get]
func (r *PublicRouter) GetEnumByType(c *gin.Context) {
	enumType := c.Param("type")

	// Khai báo enumProviders nội bộ
	var enumProviders = map[string]interface{}{
		"product_sale_status": EnumMapToList(enums.EProductSaleStatusNames),
		"product_rent_status": EnumMapToList(enums.EProductRentStatusNames),
		"asset_history":       EnumMapToList(enums.AssetHistoryNames),
		"target_type":         EnumMapToList(enums.TargetTypeNames),
		"visibility":          EnumMapToList(enums.EVisibilityNames),
		// "transaction_status": EnumMapToList(enums.TransactionStatusNames),
		"producthistory": EnumMapToList(enums.EProductHistoryNames),
		"poststatus":     EnumMapToList(enums.EPostStatusNames),
		"doctype":        EnumMapToList(enums.EDocTypeNames),
		"assetstatus":    EnumMapToList(enums.EAssetStatusNames),

		"house_orient":      EnumMapToList(enums.EHouseOrientNames),
		"house_certificate": EnumMapToList(enums.EHouseCertificateNames),
		// Thêm các enum khác ở đây nếu cần
	}

	if values, ok := enumProviders[enumType]; ok {
		c.JSON(http.StatusOK, gin.H{
			"data": values,
		})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "enum not found",
	})
}
