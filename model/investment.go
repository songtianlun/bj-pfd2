package model

type Investments []Investment
type InvestmentMap map[string]Investment
type Investment struct {
    PID     string
    Name    string
    Money   float64
    Earning float64
    Year    int64
    Month   int64
    Day     int64
    Account string
    Type    string `default:"个人投资"`
}

func (is *Investments) ArrayToMap() *InvestmentMap {
    isMap := InvestmentMap{}
    for _, i := range *is {
        if i.PID != "" {
            isMap[i.PID] = i
        }
    }
    return &isMap
}

func (is *Investments) Compare(is2 *Investments) bool {
    if len(*is) != len(*is2) {
        return false
    }
    isMap := is.ArrayToMap()
    isMap2 := is2.ArrayToMap()
    for k, v := range *isMap {
        if v != (*isMap2)[k] {
            return false
        }
    }
    return true
}

func (nb *NotionBody) ParseInvestment() (is Investments) {
    res := nb.Results
    for i := 0; i < len(res); i++ {
        re := res[i]
        //utils.PrettyPrint(re)
        b := Investment{
            PID:     re.ID,
            Money:   re.Properties.Money1.Number,
            Year:    re.Properties.Year.Formula.Number,
            Month:   re.Properties.Month.Formula.Number,
            Day:     re.Properties.Day.Formula.Number,
            Earning: re.Properties.Earn.Number,
            Type:    "个人投资",
        }
        if len(re.Properties.Note1.Title) > 0 {
            b.Name = re.Properties.Note1.Title[0].PlainText
        }
        if len(re.Properties.RIAccount.Relation) > 0 {
            b.Account = re.Properties.RIAccount.Relation[0].ID
        }
        is = append(is, b)
    }
    return
}
