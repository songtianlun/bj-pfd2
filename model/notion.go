package model

import (
	"encoding/json"
	"io"
	"strings"
)

type NotionBodyPrams struct {
	StartCursor string `json:"start_cursor,omitempty"`
	PageSize    int32  `json:"page_size,omitempty"`
}

func (nbp *NotionBodyPrams) GetJsonString() string {
	bs, err := json.Marshal(nbp)
	if err != nil {
		return ""
	}
	return string(bs)
}

func (nbp *NotionBodyPrams) GetReader() io.Reader {
	return strings.NewReader(nbp.GetJsonString())
}
