package model

func StatisticSpend(accounts *Accounts, bills Bills) *Accounts {
	asm := accounts.ArrayToMap()
	for _, bill := range bills {
		if a, ok := (*asm)[bill.Account]; ok {
			a.Money += bill.Money
			(*asm)[bill.Account] = a
		}
	}
	return asm.MapToArray()
}

func StatisticInvestment(ias *IAccounts, is *Investments) *IAccounts {
	iam := ias.ArrayToMap()
	//(*iam)["all"] = IAccount{
	//	PID:     "all",
	//	Name:    "总计",
	//	Money:   0,
	//	Earning: 0,
	//	Type:    "个人投资",
	//}
	for _, iv := range *is {
		//all := (*iam)["all"]
		//all.Money += iv.Money
		//(*iam)["all"] = all
		if a, ok := (*iam)[iv.Account]; ok {
			a.Money += iv.Money
			(*iam)[iv.Account] = a
		}
	}
	//for k, ia := range *iam {
	//	if k != "all" {
	//		all := (*iam)["all"]
	//		all.Earning += ia.Earning
	//		(*iam)["all"] = all
	//	}
	//}
	return iam.MapToArray()
}

func StatisticAccountWithIAccount(as *Accounts, ias *IAccounts) *Accounts {
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

func StatisticBillsWithBudget(bs *Bills, bgs *Budgets) *Budgets {
	bgsm := bgs.ArrayToMap()
	for _, b := range *bs {
		if bg, ok := (*bgsm)[b.Budget]; ok {
			bg.Real += b.Money
			(*bgsm)[b.Budget] = bg
		}
	}
	return bgsm.MapToArray()
}
