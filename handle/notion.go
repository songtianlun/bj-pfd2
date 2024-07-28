package handle

import (
	"bj-pfd2/model"
	"bj-pfd2/model/notion"
	"bj-pfd2/pkg/cache"
	"bj-pfd2/pkg/constvar"
	"bj-pfd2/pkg/log"
	"bj-pfd2/pkg/rest"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

func postToNotion(nUrl string, body model.NotionBodyPrams, notionToken string) (rs string, err error) {
	nUrl = "https://api.notion.com/v1" + nUrl
	client, err := rest.Client(nUrl, "POST", body.GetReader(),
		http.Header{
			"Authorization":  {"Bearer " + notionToken},
			"Notion-Version": {"2022-02-22"},
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

	if !noCache {
		rs = cache.Get(key)
	}
	if rs != "" {
		log.Infof("Get by Cache(cCache:%v) - %v / %v", noCache, url, body.GetJsonString())
	} else {
		log.Infof("Post To Notion(nCache:%v) - %v / %v", noCache, url, body.GetJsonString())
		rs, err = postToNotion(url, body, nToken)
		if err != nil {
			err = fmt.Errorf("PostToNotion error: %v ", err)
			return
		}
		go func() {
			err = cache.Set(key, rs, constvar.CacheTimeout)
			if err != nil {
				log.Error("Set cache [%v] error: %v", key, err)
			}
		}()
	}
	return
}

func searchNotion(name string, nToken string, filter *model.NBPFilter, size int32, noCache bool) (res string, err error) {
	res, err = postNotionByCache("/search", model.NotionBodyPrams{
		Query:    name,
		Filter:   filter,
		PageSize: size}, nToken, noCache)
	if err != nil {
		err = fmt.Errorf("seatchByNotion - %v ", err)
	}
	return
}

func searchDBByNotion(name string, nToken string, noCache bool) (string, error) {
	return searchNotion(name, nToken, &model.NBPFilter{
		Value:    "database",
		Property: "object",
	}, 100, noCache)
}

func searchPageByNotion(name string, nToken string, noCache bool) (string, error) {
	return searchNotion(name, nToken, &model.NBPFilter{
		Value:    "page",
		Property: "object",
	}, 1, noCache)
}

func searchPageUrlByNotion(name string, nToken string, noCache bool) (url string) {
	res, err := searchPageByNotion(name, nToken, noCache)
	if err != nil {
		log.Errorf("Err Get Page Url - %v", err.Error())
	}

	pg := &notion.DBBody{}
	err = pg.ParseDBBody(res)
	if err != nil {
		log.Error("Parse Notion Body error:", err)
		return
	}
	//utils.PrettyPrint(pg)
	if len(pg.Results) > 0 {
		url = pg.Results[0].URL
	} else {
		log.Errorf("Search Page [%v] failed, no result.", name)
	}
	return
}

func searchDBIDByNotion(name string, nToken string, noCache bool) (id string) {
	res, err := searchDBByNotion(name, nToken, noCache)
	if err != nil {
		log.Errorf("Err Get DB ID - %v", err.Error())
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
		log.Errorf("Search DB [%v] failed, no result.", name)
	}
	return
}

func GetDbId(name string, nToken string, noCache bool) (id string) {
	key := fmt.Sprintf("notion_db_id_%s_%s", nToken, name)
	id = cache.Get(key)
	if id != "" && !noCache {
		log.Debugf("Search NDB(noCache:%v) [%v] ID [%v] by cache", noCache, key, id)
	} else {
		id = searchDBIDByNotion(name, nToken, noCache)
		if id != "" {
			go func() {
				err := cache.Set(key, id, constvar.CacheTimeout)
				if err != nil {
					log.Error("Set cache [%v] error: %v", key, err)
				}
			}()
		}
	}
	return
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
		log.Errorf("Cannot to get all notion db with empty DB id.")
		return
	}
	for maxItem < 0 || count*pSize < maxItem {
		count++
		db, err := GetNotionDbByCache(aPID, start, pSize, nToken, noCache, debug)
		if err != nil {
			log.Errorf("Cannot to get all notion db: %v", err)
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
	log.Infof("Get [%v] accounts.", len(as))
	return
}

func GetAllBills(billsPID string, nToken string, noCache bool, debug bool, maxItem int32) (bs model.Bills) {
	ns := GetAllByNotion(billsPID, nToken, noCache, debug, maxItem)
	for _, n := range ns {
		bs = append(bs, n.ParseBill()...)
	}
	log.Infof("Get [%v] bills.", len(bs))
	return
}

func GetAllInvestmentAccount(investmentAccountPID string, nToken string, noCache bool, debug bool, maxItem int32) (ias model.IAccounts) {
	ns := GetAllByNotion(investmentAccountPID, nToken, noCache, debug, maxItem)
	for _, n := range ns {
		ias = append(ias, n.ParseInvestmentAccount()...)
	}
	log.Infof("Get [%v] investment accounts.", len(ias))
	return
}

func GetAllInvestment(investmentPID string, nToken string, noCache bool, debug bool, maxItem int32) (is model.Investments) {
	ns := GetAllByNotion(investmentPID, nToken, noCache, debug, maxItem)
	for _, n := range ns {
		is = append(is, n.ParseInvestment()...)
	}
	log.Infof("Get [%v] investments.", len(is))
	return
}

func GetAllBudget(budgetPid string, nToken string, noCache bool, debug bool, maxItem int32) (bs model.Budgets) {
	ns := GetAllByNotion(budgetPid, nToken, noCache, debug, maxItem)
	for _, n := range ns {
		bs = append(bs, n.ParseBudget()...)
	}
	log.Infof("Get [%v] budgets.", len(bs))
	return
}

func GetAllData(nToken string, noCache bool) (fd model.FullData) {
	log.Infof("Get All Data with token: %s", nToken)
	wg := sync.WaitGroup{}
	fd.Token = nToken

	log.Infof("Notion BJ Url: %s", fd.HomePageUrl)

	wg.Add(6)
	go func() {
		fd.HomePageUrl = searchPageUrlByNotion("Bullet Journal", nToken, noCache)
		wg.Done()
	}()
	go func() {
		aPID := GetDbId("BJPFD-账户-DB", nToken, noCache)
		if aPID != "" {
			fd.Accounts = GetAllAccount(aPID, nToken, noCache, false, -1)
		}
		wg.Done()
	}()
	go func() {
		bPID := GetDbId("BJPFD-账本-DB", nToken, noCache)
		if bPID != "" {
			fd.Bills = GetAllBills(bPID, nToken, noCache, false, -1)
		}
		wg.Done()
	}()
	go func() {
		iaPID := GetDbId("BJPFD-投资账户-DB", nToken, noCache)
		if iaPID != "" {
			fd.IAccounts = GetAllInvestmentAccount(iaPID, nToken, noCache, false, -1)
		}
		wg.Done()
	}()
	go func() {
		ibPID := GetDbId("BJPFD-投资账本-DB", nToken, noCache)
		if ibPID != "" {
			fd.Investments = GetAllInvestment(ibPID, nToken, noCache, false, -1)
		}
		wg.Done()
	}()
	go func() {
		bgPID := GetDbId("BJPFD-预算-DB", nToken, noCache)
		if bgPID != "" {
			fd.Budgets = GetAllBudget(bgPID, nToken, noCache, false, -1)
		}
		wg.Done()
	}()

	wg.Wait()
	return
}

func TokenValid(token string) bool {
	aPID := GetDbId("BJPFD-账户-DB", token, false)
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
