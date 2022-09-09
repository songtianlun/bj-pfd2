package model

import (
	"bj-pfd2/com/utils"
	"bj-pfd2/model/chart"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Bills []Bill
type BillMap map[string]Bill
type Bill struct {
	PID       string
	Name      string
	Money     float64
	Year      int64
	Month     int64
	Day       int64
	IsTrace   bool // 非实际支出为true
	Account   string
	Budget    string
	Type      string
	UsageType string
}

func (bs *Bills) Compare(bs2 *Bills) bool {
	if len(*bs) != len(*bs2) {
		return false
	}
	bsm := bs.ArrayToMap()
	bsm2 := bs2.ArrayToMap()
	for k, v := range *bsm {
		if v != (*bsm2)[k] {
			return false
		}
	}
	return true
}

func (bs *Bills) ArrayToMap() *BillMap {
	bsm := BillMap{}
	for _, b := range *bs {
		if b.PID != "" {
			bsm[b.PID] = b
		}
	}
	return &bsm
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

func IsFunctionBillsType(bType string) (pass bool) {
	switch bType {
	case "转移", "还信用卡":
		pass = true
	default:
		pass = false
	}
	return
}

func (b *Bill) isFunctionType() bool {
	return IsFunctionBillsType(b.UsageType)
}

func (bs *Bills) Waterfall() *chart.Waterfall {
	w := &chart.Waterfall{
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
			IsTrace:   re.Properties.IsTrans.Formula.Boolean,
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
			rep += fmt.Sprintf("%s: %s\n", t.Type, utils.PrintRMB(t.Spend))
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
	yi := utils.StrToUInt64(ki[0])
	yj := utils.StrToUInt64(kj[0])
	mi := utils.StrToUInt64(ki[1])
	mj := utils.StrToUInt64(kj[1])
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
		rep += fmt.Sprintf("===>%s: %s\n", m.Month, utils.PrintRMB(m.Spend))
		array := m.TypesMap.ToArray()
		rep += array.GenerateReport()
	}
	return
}

// GenerateReport 花销统计
func (bs *Bills) GenerateReport() (rep string,
	sumSpend float64, sumIncome float64,
	sumYearSpend float64, sumYearIncome float64,
	sumMonthSpend float64, sumMonthIncome float64) {
	sort.Sort(bs) // 根据日期排序
	var BTypes = TypeSpends{}
	var BMonths = MonthSpends{}
	//var BMonth map[string]float64
	//var BDay map[string]float64
	cYear, cMonth, _ := time.Now().Date()

	BTypesMap := BTypes.ToMap()
	BMonthMap := BMonths.ToMap()
	for _, b := range *bs {
		if !b.isFunctionType() {
			if b.Money <= 0 {
				sumSpend += b.Money
				if int(b.Year) == cYear {
					sumYearSpend += b.Money
					if int(b.Month) == int(cMonth) {
						sumMonthSpend += b.Money
					}
				}
			} else {
				sumIncome += b.Money
				if int(b.Year) == cYear {
					sumYearIncome += b.Money
					if int(b.Month) == int(cMonth) {
						sumMonthIncome += b.Money
					}
				}
			}
		}

		if b.IsTrace {
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

	rep += "== 消费分类 ==\n"
	rep += BTypes.GenerateReport()
	rep += "== 消费分月 ==\n"
	rep += BMonths.GenerateReport()
	return
}
