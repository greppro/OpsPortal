package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ops-portal/config"
	"ops-portal/models"
	"fmt"
)

func GetTools(c *gin.Context) {
	fmt.Printf("Received request for environment: %s\n", c.Query("environment"))
	var tools []models.Tool
	
	env := c.Query("environment")
	if err := config.DB.Where("environment = ?", env).Find(&tools).Error; err != nil {
		fmt.Printf("Error fetching tools: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	fmt.Printf("Found %d tools\n", len(tools))
	c.JSON(http.StatusOK, tools)
}

func CreateTool(c *gin.Context) {
	var tool models.Tool
	fmt.Printf("Received tool data: %+v\n", tool)
	if err := c.ShouldBindJSON(&tool); err != nil {
		fmt.Printf("Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := config.DB.Create(&tool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, tool)
}

func UpdateTool(c *gin.Context) {
	id := c.Param("id")
	var tool models.Tool
	
	if err := config.DB.First(&tool, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tool not found"})
		return
	}
	
	if err := c.ShouldBindJSON(&tool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := config.DB.Save(&tool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, tool)
}

func DeleteTool(c *gin.Context) {
	id := c.Param("id")
	
	if err := config.DB.Delete(&models.Tool{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Tool deleted successfully"})
} 