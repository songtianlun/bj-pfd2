package handle

import "bj-pfd2/model"

func StatisticSpend(accounts *model.Accounts, bills model.Bills) *model.Accounts {
	asm := accounts.ArrayToMap()
	(*asm)["all"] = model.Account{
		PID:   "all",
		Name:  "总计",
		Money: 0,
	}
	(*asm)["savings_all"] = model.Account{
		PID:   "savings_all",
		Name:  "储蓄账户合计",
		Money: 0,
	}
	(*asm)["credit_all"] = model.Account{
		PID:   "credit_all",
		Name:  "信用账户合计",
		Money: 0,
	}
	for _, bill := range bills {
		all := (*asm)["all"]
		all.Money += bill.Money
		(*asm)["all"] = all
		if a, ok := (*asm)[bill.Account]; ok {
			a.Money += bill.Money
			(*asm)[bill.Account] = a
			if (*asm)[bill.Account].Type == "信用账户" {
				ca := (*asm)["credit_all"]
				ca.Money += bill.Money
				(*asm)["credit_all"] = ca
			} else if (*asm)[bill.Account].Type == "储蓄账户" {
				sa := (*asm)["savings_all"]
				sa.Money += bill.Money
				(*asm)["savings_all"] = sa
			}
		}
	}
	return asm.MapToArray()
}
