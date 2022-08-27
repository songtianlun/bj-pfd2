package notion

type Base struct {
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
	Object     string `json:"object"`
}

type TedBy struct {
	ID     string `json:"id"`
	Object string `json:"object"`
}

type Parent struct {
	DatabaseID string `json:"database_id,omitempty"`
	PageID     string `json:"page_id,omitempty"`
	Type       string `json:"type"`
}

type Title struct {
	Annotations Annotations `json:"annotations"`
	Href        interface{} `json:"href"`
	PlainText   string      `json:"plain_text"`
	Text        Text        `json:"text"`
	Type        string      `json:"type"`
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
}

type Text struct {
	Content string      `json:"content"`
	Link    interface{} `json:"link"`
}
