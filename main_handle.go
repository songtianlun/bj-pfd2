package main

import (
	"bj-pfd2/handle"
	"bj-pfd2/pkg/cache"
	"bj-pfd2/pkg/cfg"
	"bj-pfd2/pkg/cli"
	"bj-pfd2/pkg/log"
	"bj-pfd2/pkg/v"
	"bj-pfd2/pkg/web"
	"embed"
	"fmt"
	"github.com/arl/statsviz"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
	"net/http"
	"strconv"
)

//go:embed public
var efsStatic embed.FS

//go:embed templates
var tplEFS embed.FS

var r *chi.Mux

// Step1 - 初始化配置
func initCfg() {
	// 首先完成配置项的注册
	cfg.RegisterCfg("Port", 6010, "int64")
	cfg.RegisterCfg("ReadTimeout", 10, "int64")
	cfg.RegisterCfg("WriteTimeout", 600, "int64")
	cfg.RegisterCfg("SessionTimeoutHour", 6, "int64")
	// log
	cfg.RegisterCfg("log.level", "info", "string")
	cfg.RegisterCfg("log.file_name", "log/minegin.log", "string")
	cfg.RegisterCfg("log.max_size_mb", 1, "int")
	cfg.RegisterCfg("log.max_file_num", 64, "int")
	cfg.RegisterCfg("log.max_file_day", 7, "int")
	cfg.RegisterCfg("log.compress", false, "bool")
	cfg.RegisterCfg("log.stdout", true, "bool")
	cfg.RegisterCfg("log.only_stdout", true, "bool")
	// cache
	cfg.RegisterCfg("cache.enable", true, "bool")
	cfg.RegisterCfg("cache.type", "memory", "string")
	cfg.RegisterCfg("cache.addr", "127.0.0.1:6379", "string")
	cfg.RegisterCfg("cache.passwd", "", "string")
	cfg.RegisterCfg("cache.db", 0, "int")
	// bjpfd
	cfg.RegisterCfg("bjpfd.notion_token", "", "string")
	cfg.RegisterCfg("bjpfd.account_pid", "", "string")
	cfg.RegisterCfg("bjpfd.bills_pid", "", "string")
	cfg.RegisterCfg("bjpfd.i_account_pid", "", "string")
	cfg.RegisterCfg("bjpfd.investment_pid", "", "string")
	cfg.RegisterCfg("bjpfd.budget_pid", "", "string")

	// 之后再进行初始化
	err := cfg.Init("")
	if err != nil {
		panic(fmt.Sprintf("init cfg failed: %s", err))
	}
}

// Step2 - 初始化日志
func initLog() {
	log.InitGlobal(log.NewLogrus(&log.CfgLog{
		FileName:   cfg.GetString("log.file_name"),
		Level:      cfg.GetString("log.level"),
		MaxSizeMB:  cfg.GetInt("log.max_size_mb"),
		MaxFileNum: cfg.GetInt("log.max_file_num"),
		MaxFileDay: cfg.GetInt("log.max_file_day"),
		Compress:   cfg.GetBool("log.compress"),
		Stdout:     cfg.GetBool("log.stdout"),
		OnlyStdout: cfg.GetBool("log.only_stdout"),
	}))
}

// Step3 - 初始化缓存
func initCacheDB() {
	if cfg.GetBool("cache.enable") {
		cache.Init(
			true,
			cfg.GetString("cache.type"),
			cfg.GetString("cache.addr"),
			cfg.GetString("cache.passwd"),
			cfg.GetInt("cache.db"))
	} else {
		cache.Init(
			false,
			"", "", "", 0)
	}
}

// Step4 - 初始化 CLI 命令
func runCLI() (isCli bool) {
	cli.RegisterBoolCLI("version", "V", "show version info.", func(mapCli cli.MapCli) {
		fmt.Println(v.GetVersionStr())
	})
	cli.RegisterStringCLI("token", "T", "", "Get Report With Notion Token.", func(mapCli cli.MapCli) {
		handle.ReportWithToken(*mapCli["token"].SValue)
	})
	cli.RegisterBoolCLI("router", "R", "Generate Router Doc.", func(mapCli cli.MapCli) {
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "bj-pfd2",
			Intro:       "Welcome to the bjpfd2 router docs.",
		}))
	})
	return cli.CheckCLI()
}

// Step5 - 初始化 web 服务
func initRESTHandle() {
	r = chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	web.RegisterTplEmbedFs(&tplEFS)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(efsStatic))))
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			return
		}
	})
	r.Group(func(r chi.Router) {
		r.Use(handle.Auth)
		r.Get("/", handle.Index)
		r.Get("/home", handle.Home)
	})
	r.Group(func(r chi.Router) {
		r.Get("/err", handle.Err)
		r.Get("/login", handle.Login)
		r.Post("/authenticate", handle.Authenticate)
		r.Get("/logout", handle.Logout)
	})
	// register statsviz
	r.Get("/debug/statsviz/ws", statsviz.Ws)
	r.Get("/debug/statsviz", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/debug/statsviz/", 301)
	})
	r.Handle("/debug/statsviz/*", statsviz.Index)
}

// Step6 - 启动 web 服务
func runGlobalWebServer() {
	if r == nil {
		panic("please init web handle first.")
	}
	Addr := ":" + strconv.FormatInt(cfg.GetInt64("Port"), 10)
	log.Infof("BJ-PFD2[%v] is running on %v", v.GetVersionStr(), Addr)
	err := http.ListenAndServe(Addr, r)
	if err != nil {
		panic("web server error: " + err.Error())
	}
}
