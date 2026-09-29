package handlers

import (
	_routes "common/routes"
	"net/http"

	"search/models"
	"search/services"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.IndexDocument(doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to index"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Indexed"})
}

func SearchHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query"})
		return
	}

	docs, err := services.SearchDocuments(query)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}
	_routes.RouteResult(c, docs, nil)
	// c.JSON(http.StatusOK, gin.H{"results": docs})
}

func CraeteBatchHandler(c *gin.Context) {
	var doc []string
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.SaveKeywordsBatchBulk(doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to index"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Indexed"})
}

func DeleteHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID"})
		return
	}

	if err := services.DeleteDocument(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
