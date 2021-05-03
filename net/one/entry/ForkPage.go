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

// start
func (fp *ForkPage) Start() int {
	start := (fp.Current - 1) * fp.PageSize
	if start > fp.Total {
		return fp.Total
	}
	return start
}

// end
func (fp *ForkPage) End() int {
	end := fp.Start() + fp.PageSize
	return end
}

// end
func (fp *ForkPage) NoOverflowEnd() int {
	end := fp.Start() + fp.PageSize
	if end > fp.Total {
		end = fp.Total
	}
	return end
}

func (fp *ForkPage) Overflow() bool {
	if fp.Current > fp.Over {
		return true
	}
	return false
}
