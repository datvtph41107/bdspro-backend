package _utils

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Page represents pagination information
type Page struct {
	Limit int    `json:"limit" form:"limit,default=50"`
	Skip  int    `json:"skip" form:"skip,default=0"`
	Sort  string `json:"sort"`
}

// NewPage creates a new Page instance
func NewPage(skip, limit int, sort string) *Page {
	return &Page{Limit: limit, Skip: skip, Sort: sort}
}

// NewPage creates a new Page instance
func GetPage(c *gin.Context) Page {
	limit, _ := strconv.ParseInt(c.DefaultQuery("size", "20"), 10, 32)
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "0"), 10, 32)

	skip := page * limit
	sortStr := c.DefaultQuery("sort", "created_at,desc")
	sort := strings.Replace(sortStr, ",", " ", 1)

	return Page{Limit: int(limit), Skip: int(skip), Sort: sort}
}
