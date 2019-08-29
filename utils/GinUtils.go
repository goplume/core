package utils

import "github.com/gin-gonic/gin"

func GetParam(c *gin.Context, key string) (string, bool) {
	values := c.Param(key)
	if StringIsEmpty(values) {
		return "", false
	}
	return values, true

}
