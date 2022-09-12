package handle

import (
	"bj-pfd2/com/cache"
	"bj-pfd2/com/log"
	"bj-pfd2/com/rest"
	"bj-pfd2/model"
	"bj-pfd2/model/notion"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var cacheTimeout = time.Minute * 30

func postToNotion(nUrl string, body model.NotionBodyPrams, notionToken string) (rs string, err error) {
	log.InfoF("Post To Notion - %v / %v", nUrl, body.GetJsonString())
	nUrl = "https://api.notion.com/v1" + nUrl
	client, err := rest.Client(nUrl, "POST", body.GetReader(),
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

func postNotionByCache(url string, body model.NotionBodyPrams, nToken string, noCache bool) (rs string, err error) {
	key := fmt.Sprintf("notion_%v_body_%v_%v", url, body.GetCacheKey(), nToken)
	rs = cache.Get(key)
	if rs != "" && !noCache {
		log.InfoF("Get by Cache - %v / %v", url, body.GetJsonString())
	} else {
		rs, err = postToNotion(url, body, nToken)
		if err != nil {
			err = fmt.Errorf("PostToNotion error: %v ", err)
			return
		}
		go func() {
			err = cache.Set(key, rs, cacheTimeout)
			if err != nil {
				log.Error("Set cache [%v] error: %v", key, err)
			}
		}()
	}
	return
}

func searchByNotion(name string, nToken string, noCache bool) (res string, err error) {
	res, err = postNotionByCache("/search", model.NotionBodyPrams{
		Query: name,
		Filter: &model.NBPFilter{
			Value:    "database",
			Property: "object",
		},
		PageSize: 100}, nToken, noCache)
	if err != nil {
		err = fmt.Errorf("seatchByNotion - %v ", err)
	}
	return
}

func searchDBIDByNotion(name string, nToken string) (id string) {
	res, err := searchByNotion(name, nToken, false)
	if err != nil {
		log.ErrorF("Err GetNDID - %v", err.Error())
		return
	}
	db := &notion.DBBody{}
	//fmt.Println(utils.PrettyJsonString(res))
	err = db.ParseDBBody(res)
	if err != nil {
		log.Error("Parse Notion Body error:", err)
		return
	}
	if len(db.Results) > 0 {
		for _, r := range db.Results {
			if len(r.Title) > 0 && r.Title[0].PlainText == name {
				id = strings.Replace(r.ID, "-", "", -1)
			}
			//fmt.Println(r.Title[0].PlainText, ": ", r.ID)
		}
	} else {
		log.ErrorF("Search DB [%v] failed, no result.", name)
	}
	return
}

func searchDbIdByCache(name string, nToken string, noCache bool) (id string) {
	key := fmt.Sprintf("notion_db_id_%s_%s", nToken, name)
	id = cache.Get(key)
	if id != "" && !noCache {
		log.DebugF("Search Notion DB [%v] ID [%v] by cache", key, id)
	} else {
		id = searchDBIDByNotion(name, nToken)
		if id != "" {
			go func() {
				err := cache.Set(key, id, cacheTimeout)
				if err != nil {
					log.Error("Set cache [%v] error: %v", key, err)
				}
			}()
		}
	}
	return
}

func GetDbId(name string, nToken string) string {
	return searchDbIdByCache(name, nToken, false)
}

func GetNotionDbByCache(dbID string, start string, size int32, nToken string, noCache bool, debug bool) (nb model.NotionBody, err error) {
	c, err := postNotionByCache("/databases/"+dbID+"/query", model.NotionBodyPrams{
		StartCursor: start,
		PageSize:    size,
	}, nToken, noCache)
	if err != nil {
		log.Error("GetNotionDB error:", err)
		return
	}
	if debug {
		fmt.Println(c)
	}
	nb, err = model.ParseNotionBody(c)
	return
}

func GetAllByNotion(aPID string, nToken string, noCache bool, debug bool, maxItem int32) (ns []model.NotionBody) {
	start := ""
	var count int32
	var pSize int32
	if maxItem > 0 && maxItem < 100 {
		pSize = maxItem
	} else {
		pSize = 100
	}
	if aPID == "" {
		log.ErrorF("Cannot to get all notion db with empty DB id.")
		return
	}
	for maxItem < 0 || count*pSize < maxItem {
		count++
		db, err := GetNotionDbByCache(aPID, start, pSize, nToken, noCache, debug)
		if err != nil {
			return
		}
		ns = append(ns, db)
		if db.HasMore {
			start = db.NextCursor
			continue
		} else {
			break
		}
	}
	return
}

func GetAllAccount(aPID string, nToken string, noCache bool, debug bool, maxItem int32) (as model.Accounts) {
	ns := GetAllByNotion(aPID, nToken, noCache, debug, maxItem)
	for _, n := range ns {
		as = append(as, n.ParseAccount()...)
	}
	log.InfoF("Get [%v] accounts.", len(as))
	return
}

func GetAllBills(billsPID string, nToken string, noCache bool, debug bool, maxItem int32) (bs model.Bills) {
	ns := GetAllByNotion(billsPID, nToken, noCache, debug, maxItem)
	for _, n := range ns {
		bs = append(bs, n.ParseBill()...)
	}
	log.InfoF("Get [%v] bills.", len(bs))
	return
}

func GetAllInvestmentAccount(investmentAccountPID string, nToken string, noCache bool, debug bool, maxItem int32) (ias model.IAccounts) {
	ns := GetAllByNotion(investmentAccountPID, nToken, noCache, debug, maxItem)
	for _, n := range ns {
		ias = append(ias, n.ParseInvestmentAccount()...)
	}
	log.InfoF("Get [%v] investment accounts.", len(ias))
	return
}

func GetAllInvestment(investmentPID string, nToken string, noCache bool, debug bool, maxItem int32) (is model.Investments) {
	ns := GetAllByNotion(investmentPID, nToken, noCache, debug, maxItem)
	for _, n := range ns {
		is = append(is, n.ParseInvestment()...)
	}
	log.InfoF("Get [%v] investments.", len(is))
	return
}

func GetAllBudget(budgetPid string, nToken string, noCache bool, debug bool, maxItem int32) (bs model.Budgets) {
	ns := GetAllByNotion(budgetPid, nToken, noCache, debug, maxItem)
	for _, n := range ns {
		bs = append(bs, n.ParseBudget()...)
	}
	log.InfoF("Get [%v] budgets.", len(bs))
	return
}

func GetAllData(nToken string, noCache bool) (fd model.FullData) {
	log.InfoF("Get All Data with token: %s", nToken)
	wg := sync.WaitGroup{}
	fd.Token = nToken

	wg.Add(5)
	go func() {
		aPID := GetDbId("BJPFD-账户-DB", nToken)
		if aPID != "" {
			fd.Accounts = GetAllAccount(aPID, nToken, noCache, false, -1)
		}
		wg.Done()
	}()
	go func() {
		bPID := GetDbId("BJPFD-账本-DB", nToken)
		if bPID != "" {
			fd.Bills = GetAllBills(bPID, nToken, noCache, false, -1)
		}
		wg.Done()
	}()
	go func() {
		iaPID := GetDbId("BJPFD-投资账户-DB", nToken)
		if iaPID != "" {
			fd.IAccounts = GetAllInvestmentAccount(iaPID, nToken, noCache, false, -1)
		}
		wg.Done()
	}()
	go func() {
		ibPID := GetDbId("BJPFD-投资账本-DB", nToken)
		if ibPID != "" {
			fd.Investments = GetAllInvestment(ibPID, nToken, noCache, false, -1)
		}
		wg.Done()
	}()
	go func() {
		bgPID := GetDbId("BJPFD-预算-DB", nToken)
		if bgPID != "" {
			fd.Budgets = GetAllBudget(bgPID, nToken, noCache, false, -1)
		}
		wg.Done()
	}()

	wg.Wait()
	return
}

func TokenValid(token string) bool {
	aPID := GetDbId("BJPFD-账户-DB", token)
	if aPID == "" {
		return false
	} else {
		return true
	}
}

func ReportWithToken(token string) {
	fd := GetAllData(token, false)
	fd.Report()
}
