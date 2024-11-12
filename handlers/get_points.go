package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler for returning the number of points for a particular ID
func GetPoints(c *gin.Context) {
	id := c.Param("id")
	points, exists := Scores[id]
	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"description": "No receipt found for that id"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"points": points})
}
