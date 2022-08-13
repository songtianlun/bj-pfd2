package model

import (
	"fmt"
	"sort"
)

type Accounts []Account
type AccountMap map[string]Account
type Account struct {
	PID      string
	Name     string
	Money    float64
	IMoney   float64 // Investment money
	IEarning float64 // Investment earning
	Type     string
}

func (asm *AccountMap) MapToArray() *Accounts {
	accounts := Accounts{}
	for _, v := range *asm {
		accounts = append(accounts, v)
	}
	return &accounts
}

func (as *Accounts) ArrayToMap() *AccountMap {
	asm := AccountMap{}
	for _, a := range *as {
		if a.PID != "" {
			asm[a.PID] = a
		}
	}
	return &asm
}

func (as *Accounts) GenerateReport() string {
	var s string
	sort.Sort(as)
	s += "===== 账户报告 =====\n"
	//s += "账户\t账户余额\t投资总额\n"
	var cas float64 // 信用账户总额
	var sas float64 // 储蓄账户总额
	var im float64  // 投资总额
	for _, a := range *as {
		//if a.Name == "信用账户合计" {
		//	cas = a.Money
		//	continue
		//} else if a.Name == "储蓄账户合计" {
		//	s += fmt.Sprintf("%s:\t%.2f (%.2f)\n", "账户合计", a.Money, cas)
		//	continue
		//} else if a.Name == "总计" {
		//	continue
		//}
		im += a.IMoney
		sas += a.IEarning
		if a.Type == "信用账户" {
			cas += a.Money
		} else {
			sas += a.Money
		}
		s += fmt.Sprintf("%s:%.2f (投资：%.2f)\n", a.Name, a.Money+a.IEarning, a.IMoney)
	}
	s += fmt.Sprintf("%s: %.2f (%.2f)\n", "账户合计", sas, cas)
	return s
}

func (as *Accounts) Len() int {
	return len(*as)
}

func (as *Accounts) Swap(i, j int) {
	(*as)[i], (*as)[j] = (*as)[j], (*as)[i]
}

func (as *Accounts) Less(i, j int) bool {
	return (*as)[i].Money < (*as)[j].Money
}

func (nb *NotionBody) ParseAccount() (accounts Accounts) {
	res := nb.Results
	for i := 0; i < len(res); i++ {
		re := res[i]
		//utils.PrettyPrint(re)
		a := Account{
			PID:   re.ID,
			Money: 0,
			Type:  re.Properties.AType.Select.Name,
		}
		if len(re.Properties.Name.Title) > 0 {
			a.Name = re.Properties.Name.Title[0].PlainText
		}
		accounts = append(accounts, a)
	}
	return accounts
}
