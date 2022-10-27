package main

func main() {
    // Step1 - 初始化配置
    initCfg()
    // Step2 - 初始化日志
    initLog()
    // Step3 - 初始化缓存数据库
    initCacheDB()
    // Step4 - 初始化web处理器
    initRESTHandle()
    // Step5 - 检查是否作为终端程序并运行
    if runCLI() {
        return
    }
    // Step6 - 启动web服务
    runGlobalWebServer()
}
