package model

type Bills []Bill
type Bill struct {
	PID     string
	Name    string
	Money   float64
	Year    int64
	Month   int64
	Day     int64
	Trace   bool
	Account string
	Budget  string
	Type    string
}

func (nb *NotionBody) ParseBill() (bills Bills) {
	res := nb.Results
	for i := 0; i < len(res); i++ {
		re := res[i]
		//utils.PrettyPrint(re)
		b := Bill{
			PID:   re.ID,
			Money: re.Properties.Money.Number,
			Year:  re.Properties.Year.Formula.Number,
			Month: re.Properties.Month.Formula.Number,
			Day:   re.Properties.Day.Formula.Number,
			Trace: re.Properties.IsTrans.Formula.Boolean,
			Type:  "个人储蓄",
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
