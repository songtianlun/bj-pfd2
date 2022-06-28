package model

type IAccounts []IAccount
type IAccount struct {
	PID     string
	RAID    string // 对应的账户
	Name    string
	Money   float64
	Earning float64
	Type    string
}

func (nb *NotionBody) ParseInvestmentAccount() (iAccounts IAccounts) {
	res := nb.Results
	for i := 0; i < len(res); i++ {
		re := res[i]
		//utils.PrettyPrint(re)
		ia := IAccount{
			PID:     re.ID,
			Money:   0,
			Earning: re.Properties.Earn.Number,
			Type:    "个人投资",
		}
		if len(re.Properties.Note.Title) > 0 {
			ia.Name = re.Properties.Note.Title[0].PlainText
		}
		if len(re.Properties.RAccount.Relation) > 0 {
			ia.RAID = re.Properties.RAccount.Relation[0].ID
		}
		iAccounts = append(iAccounts, ia)
	}
	return
}
