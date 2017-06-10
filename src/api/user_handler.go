package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

func (a *Api) registerUser(c *gin.Context) {
	user := &User{}
	err := c.Bind(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Couldn't parse form: %s", err.Error())})
		return
	}

	user, err = registerNewUser(a.DB, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Couldn't register user: %s", err.Error())})
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
