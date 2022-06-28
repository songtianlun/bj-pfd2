package model

import (
	"fmt"
	"sort"
)

type Accounts []Account
type AccountMap map[string]Account
type Account struct {
	PID   string
	Name  string
	Money float64
	Type  string
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
	s += "账户报告：\n"
	for _, a := range *as {
		s += fmt.Sprintf("%s:\t%.2f\n", a.Name, a.Money)
	}
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
		if n := re.Properties.Name.Title[0].PlainText; n != "" {
			accounts = append(accounts, Account{
				PID:   re.ID,
				Name:  n,
				Money: 0,
				Type:  "个人储蓄",
			})
		}
	}
	return accounts
}
