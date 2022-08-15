package model

import (
	"bj-pfd2/com/utils"
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

// GenerateReport
// rep - 报告内容
// sas - 储蓄账户总额
// cas - 信用账户总额
// im - 投资总额
func (as *Accounts) GenerateReport() (rep string, sas float64, cas float64, im float64) {
	sort.Sort(as)
	rep += "===== 账户报告 =====\n"
	for _, a := range *as {
		im += a.IMoney
		sas += a.IEarning
		if a.Type == "信用账户" {
			cas += a.Money
		} else {
			sas += a.Money
		}
		rep += fmt.Sprintf("%s:%s (投资：%s)\n", a.Name,
			utils.PrintRMB(a.Money+a.IEarning), utils.PrintRMB(a.IMoney))
	}
	rep += fmt.Sprintf("%s: %s (%s)\n", "账户合计", utils.PrintRMB(sas), utils.PrintRMB(cas))
	return
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
