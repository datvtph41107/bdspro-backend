package routes

// RegisterRoutes đăng ký các routes cho API
// func RegisterRoutes(router *gin.Engine, service *service.MapService) {
// 	r := router.Group("/map/locations")

// 	// @Summary Tìm các địa điểm lân cận
// 	// @Description Trả về danh sách các điểm gần tọa độ đã cho trong phạm vi bán kính
// 	// @Tags Locations
// 	// @Accept json
// 	// @Produce json
// 	// @Param lat query float64 true "Vĩ độ"
// 	// @Param lng query float64 true "Kinh độ"
// 	// @Param radius query float64 false "Bán kính tìm kiếm (mặc định: 5km)" default(5)
// 	// @Router /map/locations/nearby [get]
// 	r.GET("/nearby", func(c *gin.Context) {
// 		keys := []string{"lat", "lng", "radius"}
// 		values := make([]float64, len(keys))

// 		for i, key := range keys {
// 			valueStr := c.Query(key)
// 			value, err := strconv.ParseFloat(valueStr, 64)
// 			if err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid " + key + " value"})
// 				return
// 			}
// 			values[i] = value
// 		}

// 		// Convert values and handle errors (omitted for brevity)
// 		locations, err := service.FindNearbyLocations(values[0], values[1], values[2])
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(200, locations)
// 	})

// 	// @Summary Thêm địa điểm mới
// 	// @Description Thêm một điểm mới vào cơ sở dữ liệu
// 	// @Tags Locations
// 	// @Accept json
// 	// @Produce json
// 	// @Param location body MapPointDTO true "Thông tin địa điểm"
// 	// @Success 200 {object} map[string]string
// 	// @Failure 400 {object} map[string]string
// 	// @Router /map/locations/new [post]
// 	r.POST("/new", func(c *gin.Context) {
// 		var dto models.MapPointDTO
// 		if err := c.ShouldBindJSON(&dto); err != nil {
// 			c.JSON(400, gin.H{"error": "Invalid request body"})
// 			return
// 		}
// 		location, err := service.NewLocation(&dto)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(200, location)
// 	})

// 	// @Summary Tìm địa điểm theo vùng (polygon)
// 	// @Description Tìm các điểm nằm trong một vùng đa giác
// 	// @Tags Locations
// 	// @Accept json
// 	// @Produce json
// 	// @Param polygon body []MapPointDTO true "Danh sách tọa độ tạo thành vùng đa giác"
// 	// @Success 200 {object} []MapPointDTO
// 	// @Failure 400 {object} map[string]string
// 	// @Router /map/locations/find-by-polygon [post]
// 	r.POST("/find-by-polygon", func(c *gin.Context) {
// 		var polygon []models.MapPointDTO
// 		if err := c.ShouldBindJSON(&polygon); err != nil {
// 			c.JSON(400, gin.H{"error": "Invalid request body"})
// 			return
// 		}
// 		locations, err := service.FindLocationByPolygon(polygon)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(200, locations)
// 	})
// }
