package model

type Investments []investment
type investment struct {
	PID     string
	Name    string
	Money   float64
	Year    int64
	Month   int64
	Day     int64
	Account string
	Type    string `default:"个人投资"`
}

func (nb *NotionBody) ParseInvestment() (is Investments) {
	res := nb.Results
	for i := 0; i < len(res); i++ {
		re := res[i]
		//utils.PrettyPrint(re)
		b := investment{
			PID:   re.ID,
			Money: re.Properties.Money1.Number,
			Year:  re.Properties.Year.Formula.Number,
			Month: re.Properties.Month.Formula.Number,
			Day:   re.Properties.Day.Formula.Number,
			Type:  "个人投资",
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
