package handle

import "bj-pfd2/model"

func StatisticSpend(accounts *model.Accounts, bills model.Bills) *model.Accounts {
	asm := accounts.ArrayToMap()
	//(*asm)["all"] = model.Account{
	//	PID:   "all",
	//	Name:  "总计",
	//	Money: 0,
	//}
	//(*asm)["savings_all"] = model.Account{
	//	PID:   "savings_all",
	//	Name:  "储蓄账户合计",
	//	Money: 0,
	//}
	//(*asm)["credit_all"] = model.Account{
	//	PID:   "credit_all",
	//	Name:  "信用账户合计",
	//	Money: 0,
	//}
	for _, bill := range bills {
		//all := (*asm)["all"]
		//all.Money += bill.Money
		//(*asm)["all"] = all
		if a, ok := (*asm)[bill.Account]; ok {
			a.Money += bill.Money
			(*asm)[bill.Account] = a
			//if (*asm)[bill.Account].Type == "信用账户" {
			//	ca := (*asm)["credit_all"]
			//	ca.Money += bill.Money
			//	(*asm)["credit_all"] = ca
			//} else if (*asm)[bill.Account].Type == "储蓄账户" {
			//	sa := (*asm)["savings_all"]
			//	sa.Money += bill.Money
			//	(*asm)["savings_all"] = sa
			//}
		}
	}
	return asm.MapToArray()
}

func StatisticInvestment(ias *model.IAccounts, is *model.Investments) *model.IAccounts {
	iam := ias.ArrayToMap()
	(*iam)["all"] = model.IAccount{
		PID:     "all",
		Name:    "总计",
		Money:   0,
		Earning: 0,
		Type:    "个人投资",
	}
	for _, iv := range *is {
		all := (*iam)["all"]
		all.Money += iv.Money
		(*iam)["all"] = all
		if a, ok := (*iam)[iv.Account]; ok {
			a.Money += iv.Money
			(*iam)[iv.Account] = a
		}
	}
	for k, ia := range *iam {
		if k != "all" {
			all := (*iam)["all"]
			all.Earning += ia.Earning
			(*iam)["all"] = all
		}
	}
	return iam.MapToArray()
}

func StatisticAccountWithIAccount(as *model.Accounts, ias *model.IAccounts) *model.Accounts {
	asm := as.ArrayToMap()
	for _, ia := range *ias {
		if a, ok := (*asm)[ia.RAID]; ok {
			a.IMoney += ia.Money
			a.IEarning += ia.Earning
			(*asm)[ia.RAID] = a
		}
	}
	return asm.MapToArray()
}
