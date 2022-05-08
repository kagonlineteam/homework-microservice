package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kagonlineteam/homework-microservice/middleware"
	"github.com/kagonlineteam/homework-microservice/models"
)

func ListHomework(c *gin.Context) {
	var homeworks []models.Homework

	// Do only allow access with correct befugnis
	if !c.GetBool(middleware.SHOWALL_ATTRIBUTE_GIN_NAME) && !c.GetBool(middleware.TEACHER_ATTRIBUTE_GIN_NAME) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not allowed to access this endpoint."})
		return
	}

	query := models.DB.Begin()

	if c.Query("class") != "" {
		query = query.Where("class = ?", c.Query("class"))
	}

	if c.Query("grade") != "" {
		query = query.Where("grade = ?", c.Query("grade"))
	}

	if c.Query("deadline-geq") != "" {
		query = query.Where("deadline >= ?", c.Query("deadline-geq"))
	}

	if c.Query("deadline-leq") != "" {
		query = query.Where("deadline <= ?", c.Query("deadline-leq"))
	}

	if c.Query("reported") != "" {
		if c.Query("reported") == "1" {
			query = query.Where("reported IS NOT NULL")
		} else {
			query = query.Where("reported IS NULL")
		}
	}

	if c.Query("reportedBy") != "" {
		query = query.Where("reported = ?", c.Query("reportedBy"))
	}

	if c.Query("author") != "" {
		query = query.Where("author = ?", c.Query("author"))
	}

	if c.Query("course") != "" {
		query = query.Where("course = ?", c.Query("course"))
	}

	if c.Query("limit") != "" {
		query = query.Limit(c.Query("limit"))
	}

	if c.Query("offset") != "" {
		query = query.Offset(c.Query("offset"))
	}

	query.Find(&homeworks)

	c.JSON(http.StatusOK, gin.H{"message": "Hausaufgaben geladen.", "entities": homeworks})
}

func DeleteHomework(c *gin.Context) {
	var homework models.Homework
	if err := models.DB.Where("id = ?", c.Param("id")).First(&homework).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Homework not found"})
		return
	}

	// Do only allow people with homework-delete to delete
	if !c.GetBool(middleware.ALLOWDELETE_ATTRIBUTE_GIN_NAME) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You can not delete this homework"})
		return
	}

	models.DB.Model(&homework).Delete(&homework)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
