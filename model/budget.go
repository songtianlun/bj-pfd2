package model

import (
    "bj-pfd2/pkg/utils"
    "fmt"
    "sort"
)

type Budgets []Budget
type BudgetsMap map[string]Budget
type Budget struct {
    PID    string
    Money  float64
    Real   float64
    Remain float64
    Year   int64
    Month  int64
}

func (bgm *BudgetsMap) MapToArray() *Budgets {
    bgs := Budgets{}
    for _, v := range *bgm {
        bgs = append(bgs, v)
    }
    return &bgs
}

func (bgs *Budgets) ArrayToMap() *BudgetsMap {
    bgm := BudgetsMap{}
    for _, bg := range *bgs {
        if bg.PID != "" {
            bgm[bg.PID] = bg
        }
    }
    return &bgm
}

func (bs *Budgets) Compare(bs2 *Budgets) bool {
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

func (bgs *Budgets) StatisticRemain() {
    for i, b := range *bgs {
        ((*bgs)[i]).Remain = b.Money + b.Real
    }
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
    s += "年-月 \t 预算 \t 实际 \t 剩余 \t 日均\n"
    for _, bg := range *bgs {
        if bg.Year == 0 || bg.Month == 0 {
            continue
        }
        s += fmt.Sprintf(
            "%v \t %s \t %s \t %s \t %s\n",
            utils.EnDateWithYM(bg.Year, bg.Month),
            utils.PrintRMB(bg.Money),
            utils.PrintRMB(bg.Real),
            utils.PrintRMB(bg.Remain),
            utils.PrintRMB(bg.Real/float64(30)))
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
