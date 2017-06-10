package api

import (
	"encoding/json"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

func (a *Api) registerUser(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	user, err := registerNewUser(username, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": userJSON})
	return
}
