package web

import "embed"

type Web interface {
	RegisterEmbedFs(path string, fs *embed.FS)
}

func NewWebServer() {
}
