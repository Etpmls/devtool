package d_gin

import (
	d "github.com/etpmls/devtool"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 获取PageNumber设置的字段名
func GetPageNumberKey() string {
	pageNumberKey, err := d.GetField("pageNumber")
	if err != nil {
		pageNumberKey = "pageNumber"
	}
	return pageNumberKey
}

// 获取PageSize设置的字段名
func GetPageSizeKey() string {
	pageSizeKey, err := d.GetField("pageSize")
	if err != nil {
		pageSizeKey = "pageSize"
	}
	return pageSizeKey
}

// 根据PageNO和PageSize获取分页
func GetPageByQuery(c *gin.Context) (limit int, offset int) {
	pageNumberKey := GetPageNumberKey()
	pageSizeKey := GetPageSizeKey()

	limit = -1
	offset = -1

	pn := c.Query(pageNumberKey)
	pageNo, err1 := strconv.Atoi(pn)

	ps := c.Query(pageSizeKey)
	pageSize, err2 := strconv.Atoi(ps)

	if err1 == nil && err2 == nil && pageSize > 0 && pageNo > 0 {
		limit = pageSize
		offset = (pageNo - 1) * limit
	}

	return limit, offset
}
