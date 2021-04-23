package entry

import (
	"math"
)

/*
	"Current": 当前页码,
	"Over": 结束页码,
	"PageSize": 叶大小,
	"Total": 总页数,
	"Content": 数据内容
*/
type ForkPage struct {
	Current  int
	Over     int
	PageSize int
	Total    int
	Content  interface{}
}

// 设置分页信息
func (fp *ForkPage) Set(current int, size int, total int) {
	fp.Current = current
	fp.PageSize = size
	fp.Total = total
	ceil := math.Ceil(float64(total) / float64(size))
	fp.Over = int(ceil)
}

func (fp *ForkPage) Overflow() bool {
	if fp.Current > fp.Over {
		return true
	}
	return false
}
