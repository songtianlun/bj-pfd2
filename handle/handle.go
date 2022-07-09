package handle

import (
	cache2 "bj-pfd2/com/cache"
	"bj-pfd2/com/cfg"
	"bj-pfd2/com/log"
	"bj-pfd2/com/rest"
	"bj-pfd2/com/utils"
	"bj-pfd2/model"
	"fmt"
	"net/http"
)

func GetNotionDB(dbID string, start string, size int32) (rs string, err error) {
	notionToken := cfg.GetString("bjpfd.notion_token")
	url := "https://api.notion.com/v1/databases/" + dbID + "/query"
	log.DebugF("GetNotionDB [%v_%v] By Notion", dbID, start)
	body := model.NotionBodyPrams{
		StartCursor: start,
		PageSize:    size,
	}
	client, err := rest.Client(url, "POST", body.GetReader(),
		http.Header{
			"Authorization":  {"Bearer " + notionToken},
			"Notion-Version": {"2021-08-16"},
		})
	if err != nil {
		return
	}
	//fmt.Println(string(client))
	rs = string(client)
	return
}

func GetNotionDbByCache(dbID string, start string, size int32, noCache bool) (nb model.NotionBody, err error) {
	key := "notion_" + dbID + "_" + utils.Int32ToString(size) + "_" + start
	cache := cache2.Get(key)
	if cache != "" && !noCache {
		log.DebugF("Get url [%v] by cache", key)
	} else {
		cache, err = GetNotionDB(dbID, start, size)
		if err != nil {
			log.Error("GetNotionDB error:", err)
			return
		}
		//utils.PrettyPrint(nb)
		err = cache2.Set(key, cache)
		if err != nil {
			log.Error("Set cache [%v] error: %v", key, err)
		}
	}
	//fmt.Println(cache)
	nb, err = model.ParseNotionBody(cache)
	return
}

func GetAllAccount() (as model.Accounts) {
	start := ""
	accountPID := cfg.GetString("bjpfd.account_pid")
	for true {
		db, err := GetNotionDbByCache(accountPID, start, 100, false)
		if err != nil {
			return
		}
		as = append(as, db.ParseAccount()...)
		if db.HasMore {
			start = db.NextCursor
			continue
		} else {
			break
		}
	}
	log.InfoF("Success to Get [%v] Accounts.", len(as))
	return
}

func GetAllBills() (bs model.Bills) {
	start := ""
	billsPID := cfg.GetString("bjpfd.bills_pid")
	for true {
		db, err := GetNotionDbByCache(billsPID, start, 100, false)
		if err != nil {
			return
		}
		bs = append(bs, db.ParseBill()...)
		//utils.PrettyPrint(bs)
		if db.HasMore {
			start = db.NextCursor
			continue
		} else {
			break
		}
	}
	log.InfoF("Success to Get [%v] Bills.", len(bs))
	return
}

func GetAllInvestmentAccount() (ias model.IAccounts) {
	start := ""
	investmentAccountPID := cfg.GetString("bjpfd.i_account_pid")
	for true {
		db, err := GetNotionDbByCache(investmentAccountPID, start, 100, false)
		if err != nil {
			return
		}
		ias = append(ias, db.ParseInvestmentAccount()...)
		if db.HasMore {
			start = db.NextCursor
			continue
		} else {
			break
		}
	}
	log.InfoF("Success to Get [%v] InvestmentAccount.", len(ias))
	return
}

func GetAllInvestment() (is model.Investments) {
	start := ""
	investmentPID := cfg.GetString("bjpfd.investment_pid")
	for true {
		db, err := GetNotionDbByCache(investmentPID, start, 100, false)
		if err != nil {
			return
		}
		is = append(is, db.ParseInvestment()...)
		if db.HasMore {
			start = db.NextCursor
			continue
		} else {
			break
		}
	}
	log.InfoF("Success to Get [%v] Investments.", len(is))

	return
}

func GetAllBudget() (bs model.Budgets) {
	start := ""
	budgetPid := cfg.GetString("bjpfd.budget_pid")
	for true {
		db, err := GetNotionDbByCache(budgetPid, start, 100, false)
		if err != nil {
			return
		}
		bs = append(bs, db.ParseBudget()...)
		if db.HasMore {
			start = db.NextCursor
			continue
		} else {
			break
		}
	}
	log.InfoF("Success to Get [%v] Budgets.", len(bs))
	return
}

func TestCode() {
	GetAllAccount()
	abs := GetAllAccount()
	bs := GetAllBills()
	abs = *StatisticSpend(&abs, bs)

	ias := GetAllInvestmentAccount()
	is := GetAllInvestment()
	ias = *StatisticInvestment(&ias, &is)
	abs = *StatisticAccountWithIAccount(&abs, &ias)

	bgs := GetAllBudget()
	bgs = *StatisticBillsWithBudget(&bs, &bgs)
	bgs.StatisticRemain()

	w := bs.Waterfall()

	fmt.Println(abs.GenerateReport())
	fmt.Println(ias.GenerateReport())
	fmt.Println(bgs.GenerateReport())
	fmt.Println(w.GenerateReport())
}
