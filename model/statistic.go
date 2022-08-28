package model

import (
	"bj-pfd2/com/log"
)

func StatisticSpend(accounts *Accounts, bills Bills) *Accounts {
	asm := accounts.ArrayToMap()
	for _, bill := range bills {
		if bill.Account == "" {
			bill.Account = "[DefaultAccount]"
		}
		if a, ok := (*asm)[bill.Account]; ok {
			a.Money += bill.Money
			(*asm)[bill.Account] = a
		}
	}
	return asm.MapToArray()
}

func StatisticInvestment(ias *IAccounts, is *Investments) *IAccounts {
	iam := ias.ArrayToMap()
	for _, iv := range *is {
		if a, ok := (*iam)[iv.Account]; ok {
			a.Money += iv.Money
			(*iam)[iv.Account] = a
		}
	}
	for k, v := range *iam {
		v.RMoney = v.Money + v.Earning
		(*iam)[k] = v
	}
	return iam.MapToArray()
}

func StatisticAccountWithIAccount(as *Accounts, ias *IAccounts) *Accounts {
	asm := as.ArrayToMap()
	for _, ia := range *ias {
		if a, ok := (*asm)[ia.RAID]; ok {
			a.IMoney += ia.RMoney
			a.IEarning += ia.Earning
			a.Money += ia.Earning
			(*asm)[ia.RAID] = a
		} else {
			log.DebugF("StatisticAccountWithIAccount: %s[%s] not found", ia.Name, ia.RAID)
		}
	}
	for k, a := range *asm {
		a.RMoney = a.Money - a.IMoney
		(*asm)[k] = a
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
