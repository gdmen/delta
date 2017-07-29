package api

import (
	"github.com/golang/glog"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"
)

func (a *Api) registerUser(c *gin.Context) {
	logPrefix := GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)
	user := &User{}
	err := c.Bind(user)
	if err != nil {
		msg := "Couldn't parse input form"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}
	glog.Infof("%s Username: %s", logPrefix, user.Username)

	// Salt and hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		msg := "Couldn't hash password"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}

	// Write to database
	result, err := a.DB.Exec(InsertUserSQL, user.Username, hashedPass)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			msg := "Username isn't available"
			glog.Errorf("%s %s: %v", logPrefix, msg, err)
			c.JSON(http.StatusBadRequest, gin.H{"message": msg})
			return
		}
		msg := "Couldn't add user to database"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}
	// Get the User.Id that was just auto-written to the database
	id, err := result.LastInsertId()
	if err != nil {
		// If the db doesn't support LastInsertId(), throw an error for now
		msg := "Internal configuration mishap"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}
	user.Id = id
	user.Password = ""

	glog.Infof("%s Success: %+v", logPrefix, user)
	c.JSON(http.StatusCreated, gin.H{"user": user})
	return
}
