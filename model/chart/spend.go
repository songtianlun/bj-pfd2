package chart

import (
	"bj-pfd2/pkg/utils"
	"fmt"
	"sort"
	"strings"
)

type Spend struct {
	Year  WYear
	Month WMonth
	Day   WDay
}

type SYear map[int64]float64

func (sy *SYear) SortKey() []int64 {
	var keys []int64
	for k := range *sy {
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

func (sy *SYear) Get(k int64) float64 {
	if v, ok := (*sy)[k]; ok {
		return v
	}
	return 0
}

type SMonth map[string]float64

func (sm *SMonth) SortKey() []string {
	var keys []string
	for k := range *sm {
		ks := strings.Split(k, "-")
		if utils.StrToUInt64(ks[0]) == 0 ||
			utils.StrToUInt64(ks[1]) == 0 {
			continue // 过滤掉0字段
		}
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ki := strings.Split(keys[i], "-")
		kj := strings.Split(keys[j], "-")
		yi := utils.StrToUInt64(ki[0])
		yj := utils.StrToUInt64(kj[0])
		mi := utils.StrToUInt64(ki[1])
		mj := utils.StrToUInt64(kj[1])
		if yi != yj {
			return yi < yj
		} else {
			return mi < mj
		}
	})
	return keys
}

type SDay map[string]float64

func (sd *SDay) SortKey() []string {
	var keys []string
	for k := range *sd {
		ks := strings.Split(k, "-")
		if utils.StrToUInt64(ks[0]) == 0 ||
			utils.StrToUInt64(ks[1]) == 0 ||
			utils.StrToUInt64(ks[2]) == 0 {
			continue // 过滤掉0字段
		}
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ki := strings.Split(keys[i], "-")
		kj := strings.Split(keys[j], "-")
		yi := utils.StrToUInt64(ki[0])
		yj := utils.StrToUInt64(kj[0])
		mi := utils.StrToUInt64(ki[1])
		mj := utils.StrToUInt64(kj[1])
		di := utils.StrToUInt64(ki[2])
		dj := utils.StrToUInt64(kj[2])
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

func (sp *Spend) GenerateReport() string {
	var s string

	s += "年度：\n"
	yk := sp.Year.SortKey()
	for _, k := range yk {
		if k == 0 {
			continue
		}
		s += fmt.Sprintf("%d: %s\n", k, utils.PrintRMB(sp.Year[k]))
	}
	s += "月度：\n"
	mk := sp.Month.SortKey()
	for _, k := range mk {
		if sp.Month[k] == 0 {
			continue
		}
		s += fmt.Sprintf("%s: %s\n", k, utils.PrintRMB(sp.Month[k]))
	}

	//s += "日度：\n"
	//dk := w.Day.SortKey()
	//for _, k := range dk {
	//    s += fmt.Sprintf("%s\t%f\n", k, w.Day[k])
	//}
	return s
}
