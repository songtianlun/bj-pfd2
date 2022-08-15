package model

import (
	"bj-pfd2/com/utils"
	"fmt"
	"sort"
	"strings"
)

// 资产瀑布统计

type Waterfall struct {
	Year  WYear
	Month WMonth
	Day   WDay
}

type WYear map[int64]float64

func (wy *WYear) SortKey() []int64 {
	var keys []int64
	for k := range *wy {
		if k == 0 {
			continue // 过滤掉0字段
		}
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

type WMonth map[string]float64

func (wm *WMonth) SortKey() []string {
	var keys []string
	for k := range *wm {
		ks := strings.Split(k, "-")
		if utils.StrToInt64(ks[0]) == 0 ||
			utils.StrToInt64(ks[1]) == 0 {
			continue // 过滤掉0字段
		}
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ki := strings.Split(keys[i], "-")
		kj := strings.Split(keys[j], "-")
		yi := utils.StrToInt64(ki[0])
		yj := utils.StrToInt64(kj[0])
		mi := utils.StrToInt64(ki[1])
		mj := utils.StrToInt64(kj[1])
		if yi != yj {
			return yi < yj
		} else {
			return mi < mj
		}
	})
	return keys
}

type WDay map[string]float64

func (wd *WDay) SortKey() []string {
	var keys []string
	for k := range *wd {
		ks := strings.Split(k, "-")
		if utils.StrToInt64(ks[0]) == 0 ||
			utils.StrToInt64(ks[1]) == 0 ||
			utils.StrToInt64(ks[2]) == 0 {
			continue // 过滤掉0字段
		}
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ki := strings.Split(keys[i], "-")
		kj := strings.Split(keys[j], "-")
		yi := utils.StrToInt64(ki[0])
		yj := utils.StrToInt64(kj[0])
		mi := utils.StrToInt64(ki[1])
		mj := utils.StrToInt64(kj[1])
		di := utils.StrToInt64(ki[2])
		dj := utils.StrToInt64(kj[2])
		if yi != yj {
			return yi < yj
		} else {
			if mi != mj {
				return mi < mj
			} else {
				return di < dj
			}
		}
	})
	return keys
}

func (w *Waterfall) GenerateReport() string {
	var s string

	s += "年度：\n"
	yk := w.Year.SortKey()
	for _, k := range yk {
		if k == 0 {
			continue
		}
		s += fmt.Sprintf("%d: %s\n", k, utils.PrintRMB(w.Year[k]))
	}
	s += "月度：\n"
	mk := w.Month.SortKey()
	for _, k := range mk {
		if w.Month[k] == 0 {
			continue
		}
		s += fmt.Sprintf("%s: %s\n", k, utils.PrintRMB(w.Month[k]))
	}

	//s += "日度：\n"
	//dk := w.Day.SortKey()
	//for _, k := range dk {
	//	s += fmt.Sprintf("%s\t%f\n", k, w.Day[k])
	//}
	return s
}
