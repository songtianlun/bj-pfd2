package handle

import "bj-pfd2/model"

func StatisticSpend(accounts *model.Accounts, bills model.Bills) *model.Accounts {
	asm := accounts.ArrayToMap()
	(*asm)["all"] = model.Account{
		PID:   "all",
		Name:  "总计",
		Money: 0,
		Type:  "个人储蓄",
	}
	for _, bill := range bills {
		all := (*asm)["all"]
		all.Money += bill.Money
		(*asm)["all"] = all
		if a, ok := (*asm)[bill.Account]; ok {
			a.Money += bill.Money
			(*asm)[bill.Account] = a
		}
	}
	return asm.MapToArray()
}
