package cli

import (
	"fmt"
	"github.com/spf13/pflag"
)

type CLI struct {
	Abbr   string      // 缩写
	Dft    interface{} // 默认值
	Desc   string      // 描述
	Value  *bool       // 值
	Handle HandleCLI
}

var MapCLI = make(map[string]*CLI)

type HandleCLI func()

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
		Value:  pflag.BoolP(k, abbr, false, desc),
		Handle: handle,
	})
	//flag.BoolVar(&value, k, false, desc)
}

//func RegisterStringCLI(k string, addr string, desc string, handle HandleCLI) {
//	RegisterCLI(k, &CLI{
//		Abbr:   "",
//		Dft:    addr,
//		Desc:   desc,
//		Value:  pflag.StringP(k, "", addr, desc),
//		Handle: handle,
//	})
//}

func HandleBool(cli *CLI) () {

}

func CheckCLI() (isCli bool) {
	pflag.Parse()
	//flag.Parse()
	for _, v := range MapCLI {
		if *v.Value {
			v.Handle()
			isCli = true
			break
		}
	}
	return
}
