package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// parseID 读取路径里的 :id。必须是正整数：GORM 收到非数字字符串时会把它当成
// 原始 SQL 条件拼进 WHERE，直接把路径参数传给 First/Delete 有注入风险。
func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return 0, false
	}
	return uint(id), true
}

// tooLong 判断去掉首尾空白后的字符数是否超过上限。
func tooLong(s string, max int) bool {
	return utf8.RuneCountInString(strings.TrimSpace(s)) > max
}
