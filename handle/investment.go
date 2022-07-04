package handle

import "bj-pfd2/model"

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
