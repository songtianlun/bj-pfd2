package model

import "encoding/json"

func ParseNotionBody(body string) (NotionBody, error) {
	var nbp NotionBody
	err := json.Unmarshal([]byte(body), &nbp)
	if err != nil {
		return nbp, err
	}
	return nbp, nil
}

func (nb *NotionBody) GetJsonString() string {
	bs, err := json.Marshal(nb)
	if err != nil {
		return ""
	}
	return string(bs)
}

type NotionBody struct {
	HasMore    bool     `json:"has_more"`
	NextCursor string   `json:"next_cursor"`
	Object     string   `json:"object"`
	Results    []Result `json:"results"`
}

type Result struct {
	Archived       bool        `json:"archived"`
	Cover          interface{} `json:"cover"`
	CreatedBy      TedBy       `json:"created_by"`
	CreatedTime    string      `json:"created_time"`
	Icon           interface{} `json:"icon"`
	ID             string      `json:"id"`
	LastEditedBy   TedBy       `json:"last_edited_by"`
	LastEditedTime string      `json:"last_edited_time"`
	Object         string      `json:"object"`
	Parent         Parent      `json:"parent"`
	Properties     Properties  `json:"properties"`
	URL            string      `json:"url"`
}

type Properties struct {
	CreatedTime    CreatedTime    `json:"Created time"`
	Day            Day            `json:"Day"`
	DayOfWeek      Day            `json:"DayOfWeek"`
	LastEditedTime LastEditedTime `json:"Last edited time"`
	Month          Day            `json:"Month"`
	Week           Week           `json:"Week"`
	Year           Day            `json:"Year"`
	IsTrans        IsTrans        `json:"isTrans"`
	Name           PName          `json:"名称"`
	Note           PName          `json:"备注"`
	Money          PNumber        `json:"数额"`
	RAccount       RAccount       `json:"关联账户"`
	// for account
	AType AType `json:"类型"`
	// for investment account
	Earn Earning `json:"收益"`
	// for investment
	Note1     PName    `json:"Note"`
	RIAccount RAccount `json:"关联投资账户"`
	Money1    PNumber  `json:"本金"`
	// for budget
	Money2  PNumber `json:"预算"`
	RlMoney RNumber `json:"实际花销"`
	REMoney FNumber `json:"剩余"`
}

type AType struct {
	ID     string `json:"id"`
	Select Select `json:"select"`
	Type   string `json:"type"`
}

type Select struct {
	Color string `json:"color"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type FNumber struct {
	Formula Formula `json:"formula"`
	ID      string  `json:"id"`
	Type    string  `json:"type"`
}

type Formula struct {
	Number float64 `json:"number"`
	Type   string  `json:"type"`
}

type RNumber struct {
	ID     string `json:"id"`
	Rollup Rollup `json:"rollup"`
	Type   string `json:"type"`
}

type Rollup struct {
	Function string  `json:"function"`
	Number   float64 `json:"number"`
	Type     string  `json:"type"`
}

type Earning struct {
	ID     string  `json:"id"`
	Number float64 `json:"number"`
	Type   string  `json:"type"`
}

type RAccount struct {
	ID       string     `json:"id"`
	Relation []Relation `json:"relation"`
	Type     string     `json:"type"`
}

type Relation struct {
	ID string `json:"id"`
}

type PNumber struct {
	ID     string  `json:"id"`
	Number float64 `json:"number"`
	Type   string  `json:"type"`
}

type PName struct {
	ID    string  `json:"id"`
	Title []Title `json:"title"`
	Type  string  `json:"type"`
}

type Title struct {
	Annotations map[string]interface{} `json:"annotations"`
	Href        interface{}            `json:"href"`
	PlainText   string                 `json:"plain_text"`
	Text        TiText                 `json:"text"`
	Type        string                 `json:"type"`
}

type TiText struct {
	Context string `json:"content"`
	Link    string `json:"link"`
}

type TedBy struct {
	ID     string `json:"id"`
	Object string `json:"object"`
}

type Parent struct {
	DatabaseID string `json:"database_id"`
	Type       string `json:"type"`
}

type CreatedTime struct {
	CreatedTime string `json:"created_time"`
	ID          string `json:"id"`
	Type        string `json:"type"`
}

type Day struct {
	Formula DayFormula `json:"formula"`
	ID      string     `json:"id"`
	Type    DayType    `json:"type"`
}

type DayFormula struct {
	Number int64       `json:"number"`
	Type   FormulaType `json:"type"`
}

type IsTrans struct {
	Formula IsTransFormula `json:"formula"`
	ID      string         `json:"id"`
	Type    DayType        `json:"type"`
}

type IsTransFormula struct {
	Boolean bool   `json:"boolean"`
	Type    string `json:"type"`
}

type LastEditedTime struct {
	ID             string `json:"id"`
	LastEditedTime string `json:"last_edited_time"`
	Type           string `json:"type"`
}

type Week struct {
	Formula WeekFormula `json:"formula"`
	ID      string      `json:"id"`
	Type    DayType     `json:"type"`
}

type WeekFormula struct {
	String string `json:"string"`
	Type   string `json:"type"`
}

type FormulaType string
type DayType string
