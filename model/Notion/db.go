package Notion

import "encoding/json"

type DBBody struct {
	Base
	Results []DBResults `json:"results"`
}

type DBResults struct {
	Object         string        `json:"object"`
	ID             string        `json:"id"`
	CreatedTime    string        `json:"created_time"`
	CreatedBy      TedBy         `json:"created_by"`
	LastEditedBy   TedBy         `json:"last_edited_by"`
	LastEditedTime string        `json:"last_edited_time"`
	Description    []interface{} `json:"description"`
	IsInline       bool          `json:"is_inline"`
	Parent         Parent        `json:"parent"`
	URL            string        `json:"url"`
	Archived       bool          `json:"archived"`
	Title          []Title       `json:"title"`
}

func (dbb *DBBody) ParseDBBody(body string) (err error) {
	return json.Unmarshal([]byte(body), &dbb)
}

//func ParseDBBody(body string) (DBBody, error) {
//	var dbb DBBody
//	err := json.Unmarshal([]byte(body), &dbb)
//	if err != nil {
//		return dbb, err
//	}
//	return dbb, nil
//}
