package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kagonlineteam/homework-microservice/middleware"
	"github.com/kagonlineteam/homework-microservice/models"
)

func GetOwnHomeworks(c *gin.Context) {
	var homeworks []models.Homework

	// Do not allow my Endpoint to be used
	if c.GetString(middleware.CLASS_ATTRIBUTE_GIN_NAME) == "" || c.GetString(middleware.GRADE_ATTRIBUTE_GIN_NAME) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Your account has no grade or class."})
		return
	}

	models.DB.
		Where("class = ?", c.GetString(middleware.CLASS_ATTRIBUTE_GIN_NAME)).
		Where("grade = ?", c.GetString(middleware.GRADE_ATTRIBUTE_GIN_NAME)).
		Where("deadline >= ?", time.Now().Unix()).
		Where("reported is NULL").
		Find(&homeworks)

	c.JSON(http.StatusOK, gin.H{"message": "Hausaufgaben geladen.", "entities": homeworks})
}

func CreateHomework(c *gin.Context) {
	// Validate input
	var input models.Homework
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if input.Deadline < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deadline has to be in the future."})
		return
	}

	// Force non teachers to only allow creating
	// homeworks in there own class
	if !c.GetBool(middleware.TEACHER_ATTRIBUTE_GIN_NAME) {
		input.Class = c.GetString(middleware.CLASS_ATTRIBUTE_GIN_NAME)
		input.Grade = c.GetString(middleware.GRADE_ATTRIBUTE_GIN_NAME)
	}

	input.Author = c.GetString(middleware.USERNAME_ATTRIBUTE_GIN_NAME)
	input.Reported = nil // Do not allow to create a reported Homework

	if models.DB.Create(&input).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "entity": input})
}

func EditHomework(c *gin.Context) {
	var homework models.Homework
	if err := models.DB.Where("id = ?", c.Param("id")).First(&homework).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Homework not found"})
		return
	}

	// Do only allow teachers or students in the class to edit
	if !c.GetBool(middleware.TEACHER_ATTRIBUTE_GIN_NAME) && (homework.Class != c.GetString(middleware.CLASS_ATTRIBUTE_GIN_NAME) || homework.Grade != c.GetString(middleware.GRADE_ATTRIBUTE_GIN_NAME)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You can not edit this"})
		return
	}

	if homework.Reported != nil && !c.GetBool(middleware.TEACHER_ATTRIBUTE_GIN_NAME) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You can not edit a homework that has been reported"})
		return
	}

	// Validate input
	var input models.Homework
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	input.Author = c.GetString(middleware.USERNAME_ATTRIBUTE_GIN_NAME)
	// Force non teachers to only allow creating
	// homeworks in there own class
	if !c.GetBool(middleware.TEACHER_ATTRIBUTE_GIN_NAME) {
		input.Class = c.GetString(middleware.CLASS_ATTRIBUTE_GIN_NAME)
		input.Grade = c.GetString(middleware.GRADE_ATTRIBUTE_GIN_NAME)
	}

	// This code is needed to reset reported to nil if wished
	if input.Reported == nil && homework.Reported != nil {
		models.DB.Model(&homework).Update("reported", nil)
	}

	input.Reported = nil // Do not allow custom reporting attributes

	models.DB.Model(&homework).Updates(input)

	c.JSON(http.StatusOK, gin.H{"entity": homework})
}

func ReportHomework(c *gin.Context) {
	var homework models.Homework
	if err := models.DB.Where("id = ?", c.Param("id")).First(&homework).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Homework not found"})
		print(err.Error())
		return
	}

	// Do only allow teachers or students in the class to edit
	if !c.GetBool(middleware.TEACHER_ATTRIBUTE_GIN_NAME) && (homework.Class != c.GetString(middleware.CLASS_ATTRIBUTE_GIN_NAME) || homework.Grade != c.GetString(middleware.GRADE_ATTRIBUTE_GIN_NAME)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You can not report this"})
		return
	}

	if homework.Reported != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is already reported"})
		return
	}

	// Validate input
	var input models.Homework
	var username = c.GetString(middleware.USERNAME_ATTRIBUTE_GIN_NAME)
	input.Reported = &username

	models.DB.Model(&homework).Updates(input)

	c.JSON(http.StatusOK, gin.H{"entity": homework})
}
