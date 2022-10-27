package cli

import (
	"fmt"
	"github.com/spf13/pflag"
)

type CLI struct {
	Abbr   string      // 缩写
	Dft    interface{} // 默认值
	Desc   string      // 描述
	Type   string      // 类型 bool,string
	BValue *bool       // bool 值
	SValue *string     // string 值
	Handle HandleCLI
}

type MapCli map[string]*CLI

var MapCLI = make(MapCli)

type HandleCLI func(param MapCli)

func RegisterCLI(k string, cfg *CLI) {
	_, ok := MapCLI[k]
	if ok {
		panic(fmt.Sprintf("%s is already registered", k))
	}
	MapCLI[k] = cfg
}

func RegisterBoolCLI(k string, abbr string, desc string, handle HandleCLI) {
	RegisterCLI(k, &CLI{
		Abbr:   abbr,
		Dft:    false,
		Desc:   desc,
		Type:   "bool",
		BValue: pflag.BoolP(k, abbr, false, desc),
		Handle: handle,
	})
}

func RegisterStringCLI(k string, addr string, dft string, desc string, handle HandleCLI) {
	RegisterCLI(k, &CLI{
		Abbr:   addr,
		Dft:    dft,
		Desc:   desc,
		Type:   "string",
		SValue: pflag.StringP(k, addr, dft, desc),
		Handle: handle,
	})
}

func CheckCLI() (isCli bool) {
	pflag.Parse()
	for _, v := range MapCLI {
		if v.Type == "bool" && *v.BValue {
			v.Handle(MapCLI)
			isCli = true
			break
		} else if v.Type == "string" && *v.SValue != "" {
            v.Handle(MapCLI)
            isCli = true
            break
        }
    }
    return
}
