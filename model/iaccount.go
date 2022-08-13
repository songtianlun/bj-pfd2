package model

import (
	"fmt"
	"sort"
)

type IAccounts []IAccount
type IAccountMap map[string]IAccount
type IAccount struct {
	PID     string
	RAID    string // 对应的账户
	Name    string
	Money   float64
	Earning float64
	Type    string
}

func (ias *IAccounts) ArrayToMap() *IAccountMap {
	iam := IAccountMap{}
	for _, a := range *ias {
		if a.PID != "" {
			iam[a.PID] = a
		}
	}
	return &iam
}

func (iam *IAccountMap) MapToArray() *IAccounts {
	accounts := IAccounts{}
	for _, v := range *iam {
		accounts = append(accounts, v)
	}
	return &accounts
}

func (ias *IAccounts) Len() int {
	return len(*ias)
}

func (ias *IAccounts) Swap(i, j int) {
	(*ias)[i], (*ias)[j] = (*ias)[j], (*ias)[i]
}

func (ias *IAccounts) Less(i, j int) bool {
	return (*ias)[i].Money < (*ias)[j].Money
}

func (ias *IAccounts) GenerateReport() string {
	var s string
	sort.Sort(ias)
	s += "===== 投资账户报告 =====\n"
	//s += "账户名称\t账户余额\t收益\n"
	for _, a := range *ias {
		s += fmt.Sprintf("%s: %.2f (%.2f)\n", a.Name, a.Money, a.Earning)
	}
	return s
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
		if len(re.Properties.Name.Title) > 0 {
			ia.Name = re.Properties.Name.Title[0].PlainText
		}
		if len(re.Properties.RAccount.Relation) > 0 {
			ia.RAID = re.Properties.RAccount.Relation[0].ID
		}
		iAccounts = append(iAccounts, ia)
	}
	return
}
