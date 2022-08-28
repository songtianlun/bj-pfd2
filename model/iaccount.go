package model

import (
	"bj-pfd2/com/utils"
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
	RMoney  float64
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

func (ias *IAccounts) Compare(ias2 *IAccounts) bool {
	if len(*ias) != len(*ias2) {
		return false
	}
	iam := ias.ArrayToMap()
	iam2 := ias2.ArrayToMap()
	for k, v := range *iam {
		if v != (*iam2)[k] {
			return false
		}
	}
	return true
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

// GenerateReport
// rep - 报告
// tis - 投资总额
// tes - 收益总额
func (ias *IAccounts) GenerateReport() (rep string, tis float64, tes float64) {
	sort.Sort(ias)
	//s += "账户名称\t账户余额\t收益\n"
	for _, a := range *ias {
		rep += fmt.Sprintf("%s: %s (%s)\n", a.Name, utils.PrintRMB(a.RMoney), utils.PrintRMB(a.Earning))
		tis += a.Money
		tes += a.Earning
	}
	rep += fmt.Sprintf("%s: %s (+ %s)\n", "账户合计", utils.PrintRMB(tis), utils.PrintRMB(tes))
	return
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
