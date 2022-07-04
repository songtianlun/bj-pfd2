package model

import (
	"fmt"
	"sort"
)

type Budgets []Budget
type Budget struct {
	PID    string
	Money  float64
	Real   float64
	Remain float64
	Year   int64
	Month  int64
}

func (bgs *Budgets) Len() int {
	return len(*bgs)
}

func (bgs *Budgets) Less(i, j int) bool {
	if (*bgs)[i].Year != (*bgs)[j].Year {
		return (*bgs)[i].Year < (*bgs)[j].Year
	} else {
		return (*bgs)[i].Month < (*bgs)[j].Month
	}
}

func (bgs *Budgets) Swap(i, j int) {
	(*bgs)[i], (*bgs)[j] = (*bgs)[j], (*bgs)[i]
}

func (bgs *Budgets) GenerateReport() string {
	var s string
	sort.Sort(bgs)
	s += "预算报告：\n"
	s += "年-月\t预算\t实际\t剩余\n"
	for _, bg := range *bgs {
		if bg.Year == 0 || bg.Month == 0 {
			continue
		}
		s += fmt.Sprintf("%v-%v %.2f %.2f %.2f\n", bg.Year, bg.Month, bg.Money, bg.Real, bg.Remain)
	}
	return s
}

func (nb *NotionBody) ParseBudget() (bs Budgets) {
	res := nb.Results
	for i := 0; i < len(res); i++ {
		re := res[i]
		//utils.PrettyPrint(re)
		b := Budget{
			PID:    re.ID,
			Money:  re.Properties.Money2.Number,
			Real:   0,
			Remain: 0,
			Year:   re.Properties.Year.Formula.Number,
			Month:  re.Properties.Month.Formula.Number,
		}
		bs = append(bs, b)
	}
	return
}
