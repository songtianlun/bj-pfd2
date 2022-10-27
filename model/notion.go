package model

import (
    "encoding/json"
    "fmt"
    "io"
    "strings"
)

// NotionBodyPrams Notion 请求体参数
// ps：新增参数后记得更新 GetCacheKey()
type NotionBodyPrams struct {
    StartCursor string     `json:"start_cursor,omitempty"`
    PageSize    int32      `json:"page_size,omitempty"`
    Query       string     `json:"query,omitempty"`
    Filter      *NBPFilter `json:"filter,omitempty"`
}

type NBPFilter struct {
    Value    string `json:"value,omitempty"` // page or database
    Property string `json:"property,omitempty"`
}

//type NBPFProperty struct {
//    IsDeletedOnly          bool          `json:"isDeletedOnly,omitempty"`
//    ExcludeTemplates       bool          `json:"excludeTemplates,omitempty"`
//    IsNavigableOnly        bool          `json:"isNavigableOnly,omitempty"`
//    RequireEditPermissions bool          `json:"requireEditPermissions,omitempty"`
//    ExcludePathText        bool          `json:"excludePathText,omitempty"`
//    Ancestors              []interface{} `json:"ancestors,omitempty"`
//    CreatedBy              []interface{} `json:"createdBy,omitempty"`
//    EditedBy               []interface{} `json:"editedBy,omitempty"`
//    LastEditedTime         interface{}   `json:"lastEditedTime,omitempty"`
//    CreatedTime            interface{}   `json:"createdTime,omitempty"`
//    InTeams                []interface{} `json:"inTeams,omitempty"`
//}

func (nbp *NotionBodyPrams) GetCacheKey() (key string) {
    if nbp.StartCursor != "" {
        key += "start_cursor_" + nbp.StartCursor
    }
    if nbp.PageSize != 0 {
        key += "_page_size_" + fmt.Sprintf("%d", nbp.PageSize)
    }
    if nbp.Query != "" {
        key += "_query_" + nbp.Query
    }
    if nbp.Filter != nil {
        if nbp.Filter.Value != "" {
            key += "_filter_value_" + nbp.Filter.Value
        }
        if nbp.Filter.Property != "" {
            key += "_filter_property_" + nbp.Filter.Property
        }
    }
    return
}

func (nbp *NotionBodyPrams) GetJsonString() string {
    bs, err := json.Marshal(*nbp)
    if err != nil {
        return ""
    }
    //fmt.Println(string(bs))
    return string(bs)
}

func (nbp *NotionBodyPrams) GetReader() io.Reader {
    return strings.NewReader(nbp.GetJsonString())
}
