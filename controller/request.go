package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

const CtxUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("user not login")

func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func GetPageInfo(c *gin.Context) (page, size int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		zap.L().Debug("logic.strconv.ParseInt(offsetStr,10,64) failed", zap.Error(err))
		page = 1
	}

	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		zap.L().Debug("logic.strconv.ParseInt(limitStr,10,64) failed", zap.Error(err))
		size = 10
	}
	return
}
