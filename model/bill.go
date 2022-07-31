package model

import (
	"bj-pfd2/com/utils"
	"fmt"
	"sort"
	"strings"
)

type Bills []Bill
type Bill struct {
	PID       string
	Name      string
	Money     float64
	Year      int64
	Month     int64
	Day       int64
	Trace     bool
	Account   string
	Budget    string
	Type      string
	UsageType string
}

func (bs *Bills) Less(i, j int) bool {
	if (*bs)[i].Year != (*bs)[j].Year {
		return (*bs)[i].Year < (*bs)[j].Year
	} else if (*bs)[i].Month != (*bs)[j].Month {
		return (*bs)[i].Month < (*bs)[j].Month
	} else {
		return (*bs)[i].Day < (*bs)[j].Day
	}
}
func (bs *Bills) Len() int {
	return len(*bs)
}
func (bs *Bills) Swap(i, j int) {
	(*bs)[i], (*bs)[j] = (*bs)[j], (*bs)[i]
}

func (bs *Bills) Waterfall() *Waterfall {
	w := &Waterfall{
		Year:  make(map[int64]float64),
		Month: make(map[string]float64),
		Day:   make(map[string]float64),
	}
	for _, b := range *bs {
		month := utils.EnDateWithYM(b.Year, b.Month)
		day := utils.EnDateWithYMD(b.Year, b.Month, b.Day)
		w.Year[b.Year] += b.Money
		w.Month[month] += b.Money
		w.Day[day] += b.Money

	}
	return w
}

func (nb *NotionBody) ParseBill() (bills Bills) {
	res := nb.Results
	for i := 0; i < len(res); i++ {
		re := res[i]
		//utils.PrettyPrint(re)
		b := Bill{
			PID:       re.ID,
			Money:     re.Properties.Money.Number,
			Year:      re.Properties.Year.Formula.Number,
			Month:     re.Properties.Month.Formula.Number,
			Day:       re.Properties.Day.Formula.Number,
			Trace:     re.Properties.IsTrans.Formula.Boolean,
			UsageType: re.Properties.BUsageType.Select.Name,
			Type:      "个人储蓄",
		}
		if len(re.Properties.Note.Title) > 0 {
			b.Name = re.Properties.Note.Title[0].PlainText
		}
		if len(re.Properties.RAccount.Relation) > 0 {
			b.Account = re.Properties.RAccount.Relation[0].ID
		}
		if len(re.Properties.RBudget.Relation) > 0 {
			b.Budget = re.Properties.RBudget.Relation[0].ID
		}
		bills = append(bills, b)
	}
	return
}

type TypeSpends []TypeSpend
type TypeSpendMap map[string]TypeSpend
type TypeSpend struct {
	Type  string
	Spend float64
}

func (tsm *TypeSpendMap) ToArray() (ts TypeSpends) {
	for _, t := range *tsm {
		ts = append(ts, t)
	}
	return
}
func (ts *TypeSpends) ToMap() (tsm TypeSpendMap) {
	tsm = make(TypeSpendMap)
	for _, t := range *ts {
		if _, ok := tsm[t.Type]; !ok {
			tsm[t.Type] = t
		} else {
			cm := tsm[t.Type]
			cm.Spend += t.Spend
			tsm[t.Type] = cm
		}
	}
	return
}
func (ts *TypeSpends) Len() int {
	return len(*ts)
}
func (ts *TypeSpends) Less(i, j int) bool {
	return (*ts)[i].Spend > (*ts)[j].Spend
}
func (ts *TypeSpends) Swap(i, j int) {
	(*ts)[i], (*ts)[j] = (*ts)[j], (*ts)[i]
}
func (ts *TypeSpends) GenerateReport() (rep string) {
	sort.Sort(ts)
	for _, t := range *ts {
		if t.Spend < 0 {
			rep += fmt.Sprintf("%s: %.2f\n", t.Type, t.Spend)
		}
	}
	return
}

type MonthSpends []MonthSpend
type MonthSpendMap map[string]MonthSpend
type MonthSpend struct {
	Month    string
	Spend    float64
	TypesMap TypeSpendMap
}

func (msm *MonthSpendMap) ToArray() (ms MonthSpends) {
	for _, m := range *msm {
		ms = append(ms, m)
	}
	return
}
func (ms *MonthSpends) ToMap() (msm MonthSpendMap) {
	msm = make(MonthSpendMap)
	for _, m := range *ms {
		if _, ok := msm[m.Month]; !ok {
			msm[m.Month] = m
		} else {
			cm := msm[m.Month]
			cm.Spend += m.Spend
			cm.TypesMap = m.TypesMap
			msm[m.Month] = cm
		}
	}
	return
}
func (ms *MonthSpends) Len() int {
	return len(*ms)
}
func (ms *MonthSpends) Less(i, j int) bool {
	ki := strings.Split((*ms)[i].Month, "-")
	kj := strings.Split((*ms)[j].Month, "-")
	yi := utils.StrToInt64(ki[0])
	yj := utils.StrToInt64(kj[0])
	mi := utils.StrToInt64(ki[1])
	mj := utils.StrToInt64(kj[1])
	if yi != yj {
		return yi < yj
	} else {
		return mi < mj
	}
}
func (ms *MonthSpends) Swap(i, j int) {
	(*ms)[i], (*ms)[j] = (*ms)[j], (*ms)[i]
}
func (ms *MonthSpends) GenerateReport() (rep string) {
	sort.Sort(ms)
	for _, m := range *ms {
		rep += fmt.Sprintf("===>%s: %.2f\n", m.Month, m.Spend)
		array := m.TypesMap.ToArray()
		rep += array.GenerateReport()
	}
	return
}

// GenerateReport 花销统计
func (bs *Bills) GenerateReport() string {
	var s string
	sort.Sort(bs) // 根据日期排序
	var BTypes = TypeSpends{}
	var BMonths = MonthSpends{}
	//var BMonth map[string]float64
	//var BDay map[string]float64

	s += "===== 消费统计 =====\n"

	BTypesMap := BTypes.ToMap()
	BMonthMap := BMonths.ToMap()
	for _, b := range *bs {
		if b.Trace {
			continue
		}
		if b.UsageType != "" {
			if _, ok := BTypesMap[b.UsageType]; !ok {
				BTypesMap[b.UsageType] = TypeSpend{
					Type:  b.UsageType,
					Spend: b.Money,
				}
			} else {
				cm := BTypesMap[b.UsageType]
				cm.Spend += b.Money
				BTypesMap[b.UsageType] = cm
			}
		}
		if b.Year != 0 && b.Month != 0 {
			k := utils.EnDateWithYM(b.Year, b.Month)
			if b.UsageType == "" {
				b.UsageType = "未分类"
			}
			if _, ok := BMonthMap[k]; !ok {
				BMonthMap[k] = MonthSpend{
					Month:    k,
					Spend:    b.Money,
					TypesMap: TypeSpendMap{},
				}
			} else {
				cm := BMonthMap[k]
				cm.Spend += b.Money
				if b.UsageType != "" {
					if _, ok := cm.TypesMap[b.UsageType]; !ok {
						cm.TypesMap[b.UsageType] = TypeSpend{
							Type:  b.UsageType,
							Spend: b.Money,
						}
					} else {
						ctm := cm.TypesMap[b.UsageType]
						ctm.Spend += b.Money
						cm.TypesMap[b.UsageType] = ctm
					}
				}
				BMonthMap[k] = cm
			}
		}
	}
	BTypes = BTypesMap.ToArray()
	BMonths = BMonthMap.ToArray()

	s += "== 消费分类 ==\n"
	s += BTypes.GenerateReport()
	s += "== 消费分月 ==\n"
	s += BMonths.GenerateReport()
	return s
}
