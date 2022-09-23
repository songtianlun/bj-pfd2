package model

import (
	"bj-pfd2/model/chart"
	"bj-pfd2/pkg/utils"
	"fmt"
)

type ChartData map[string]float64

type FullData struct {
	Token string

	HomePageUrl string // Bullet Journal Notion 地址

	SumMoney               float64 // 资产总额
	SumMoneyStr            string
	SumCredit              float64 // 信用总额
	SumCreditStr           string
	SumInvestment          float64 // 投资总额
	SumInvestmentStr       string
	SumInvestmentIncome    float64 // 投资收益总额
	SumInvestmentIncomeStr string
	SumBillsSpend          float64 // 账单支出总额
	SumBillsSpendStr       string
	SumBillsIncome         float64 // 账单收入总额
	SumBillsIncomeStr      string
	SumThisYearSpend       float64 // 本年支出总额
	SumThisYearSpendStr    string
	SumThisYearIncome      float64 // 本年收入总额
	SumThisYearIncomeStr   string
	SumThisMonthSpend      float64 // 本月支出总额
	SumThisMonthSpendStr   string
	SumThisMonthIncome     float64 // 本月收入总额
	SumThisMonthIncomeStr  string

	OverviewType ChartData // 账单支出类型概况
	Overview     ChartData // 账户概况

	WaterfallYearAll ChartData // 流水年度报表-总额
	WaterfallYearAdd ChartData // 流水年度报表-增量
	WaterfallYearSub ChartData // 流水年度报表-减量

	WaterfallMonthAll ChartData // 流水月度报表-总额
	WaterfallMonthAdd ChartData // 流水月度报表-增量
	WaterfallMonthSub ChartData // 流水月度报表-减量

	BudgetMonth ChartData // 预算月度报表-预算额
	BudgetReal  ChartData // 预算月度报表-实际额

	SpendYear  ChartData // 支出年度报表
	SpendMonth ChartData // 支出月度报表
	SpendDay   ChartData // 支出日报表

	InvestmentYear  ChartData // 投资年度报表
	InvestmentMonth ChartData // 投资月度报表

	Accounts    Accounts         // Accounts Raw Data
	Bills       Bills            // Bills Raw Data
	IAccounts   IAccounts        // Investment Accounts RawData
	Investments Investments      // Investments RawData
	Budgets     Budgets          // Budgets RawData
	Waterfall   chart.Waterfall  // Waterfall Statistic
	Spend       chart.Spend      // Spend Statistic
	Investment  chart.Investment // Investment Statistic

	AccountsReport    string
	BillsReport       string
	InvestmentsReport string
	IAccountsReport   string
	BudgetsReport     string
	WaterfallReport   string

	hasStatistic bool // 是否已经统计过
}

func (fd *FullData) Compare(tfd *FullData) bool {
	if fd.Token != tfd.Token {
		return false
	}

	if !fd.Accounts.Compare(&tfd.Accounts) {
		return false
	}

	if !fd.IAccounts.Compare(&tfd.IAccounts) {
		return false
	}

	if !fd.Bills.Compare(&tfd.Bills) {
		return false
	}

	if !fd.Investments.Compare(&tfd.Investments) {
		return false
	}

	if !fd.Budgets.Compare(&tfd.Budgets) {
		return false
	}

	return true
}

func (fd *FullData) StatisticAll() {
	if fd.hasStatistic {
		return
	}

	fd.StatisticAllRawData() // Step1: 统计所有原始数据

	fd.StatisticAllChartData() // Step2: 统计所有图表数据

	fd.hasStatistic = true // Step3: 标记已经统计过
}

// StatisticAllRawData
// 统计所有原始数据（将相关数据结合）
func (fd *FullData) StatisticAllRawData() {
	if fd.hasStatistic {
		return
	}
	fd.Accounts = *StatisticBillsToAccounts(&fd.Accounts, fd.Bills)                  // Step1: 统计账单支出
	fd.IAccounts = *StatisticInvestmentWithIAccounts(&fd.IAccounts, &fd.Investments) // Step2： 然后统计投资状况
	fd.Accounts = *StatisticAccountWithIAccount(&fd.Accounts, &fd.IAccounts)         // Step3: 将投资状况统计到账户中
	fd.Budgets = *StatisticBillsWithBudget(&fd.Bills, &fd.Budgets)                   // Step4: 统计账单支出到预算中
	fd.Budgets.StatisticRemain()                                                     // Step5: 统计预算剩余金额
	fd.Waterfall = chart.Waterfall{
		Year:  make(map[int64]float64),
		Month: make(map[string]float64),
		Day:   make(map[string]float64),
	}
	fd.Spend = chart.Spend{
		Year:  make(map[int64]float64),
		Month: make(map[string]float64),
		Day:   make(map[string]float64),
	}
	StatisticBills(&fd.Bills, &fd.Waterfall, &fd.Spend) // Step6: 统计流水和纯支出
	fd.Investment = chart.Investment{
		Year:  make(map[int64]float64),
		Month: make(map[string]float64),
	}
	StatisticInvestment(&fd.Investments, &fd.Investment) // Step7: 统计投资
}

// StatisticAllChartData
// 统计所有图表数据（在完成原始数据统计的基础上）
func (fd *FullData) StatisticAllChartData() {
	if fd.hasStatistic {
		return
	}
	fd.AccountsReport,
		fd.SumMoney,
		fd.SumCredit,
		fd.SumInvestment = fd.Accounts.GenerateReport()

	fd.BillsReport,
		fd.SumBillsSpend,
		fd.SumBillsIncome,
		fd.SumThisYearSpend,
		fd.SumThisYearIncome,
		fd.SumThisMonthSpend,
		fd.SumThisMonthIncome = fd.Bills.GenerateReport()
	//fd.InvestmentsReport = fd.Investments.GenerateReport()

	fd.IAccountsReport,
		_, fd.SumInvestmentIncome = fd.IAccounts.GenerateReport()
	fd.BudgetsReport = fd.Budgets.GenerateReport()
	fd.WaterfallReport = fd.Waterfall.GenerateReport()

	fd.GenerateStrRMB() // Step7: 统计金额转换为人民币
	fd.GenerateChartData()

}

func (fd *FullData) GenerateStrRMB() {
	fd.SumMoneyStr = utils.Float64ToIntStrRMB(fd.SumMoney)
	fd.SumCreditStr = utils.Float64ToIntStrRMB(fd.SumCredit)
	fd.SumInvestmentStr = utils.Float64ToIntStrRMB(fd.SumInvestment)
	fd.SumInvestmentIncomeStr = utils.Float64ToIntStrRMB(fd.SumInvestmentIncome)
	fd.SumBillsSpendStr = utils.Float64ToIntStrRMB(fd.SumBillsSpend)
	fd.SumBillsIncomeStr = utils.Float64ToIntStrRMB(fd.SumBillsIncome)
	fd.SumThisYearSpendStr = utils.Float64ToIntStrRMB(fd.SumThisYearSpend)
	fd.SumThisYearIncomeStr = utils.Float64ToIntStrRMB(fd.SumThisYearIncome)
	fd.SumThisMonthSpendStr = utils.Float64ToIntStrRMB(fd.SumThisMonthSpend)
	fd.SumThisMonthIncomeStr = utils.Float64ToIntStrRMB(fd.SumThisMonthIncome)
}

func (fd *FullData) GenerateChartData() {
	fd.Overview = make(ChartData)
	fd.OverviewType = make(ChartData)

	// 遍历账户清单统计储蓄账户总额
	for _, a := range fd.Accounts {
		if a.Name == "" {
			continue
		}
		kn := fmt.Sprintf("(%s)%s", a.Type, a.Name)
		if _, ok := fd.Overview[kn]; !ok {
			fd.Overview[kn] = 0
		}
		fd.Overview[kn] += a.RMoney

		if a.Type == "" {
			a.Type = "储蓄账户"
		}

		if _, ok := fd.OverviewType[a.Type]; !ok {
			fd.OverviewType[a.Type] = 0
		}
		fd.OverviewType[a.Type] += a.RMoney
	}
	// 遍历投资清单统计投资总额
	for _, i := range fd.IAccounts {
		if i.Name == "" {
			continue
		}
		kn := fmt.Sprintf("(%s)%s", i.Type, i.Name)
		if _, ok := fd.Overview[kn]; !ok {
			fd.Overview[kn] = 0
		}
		fd.Overview[kn] += i.Money + i.Earning // 利息仅统计到投资账户

		if _, ok := fd.OverviewType[i.Type]; !ok {
			fd.OverviewType[i.Type] = 0
		}
		fd.OverviewType[i.Type] += i.Money + i.Earning
	}

	// 遍历 Waterfall.Year 生成瀑布图数据
	fd.WaterfallYearAll = make(ChartData)
	fd.WaterfallYearAdd = make(ChartData)
	fd.WaterfallYearSub = make(ChartData)
	var waterfallYearAll float64
	waterfallYearAll = 0
	for _, k := range fd.Waterfall.Year.SortKey() {
		if k == 0 {
			continue
		}
		key := fmt.Sprintf("%d", k)
		value := fd.Waterfall.Year[k]

		waterfallYearAll += value

		if fd.Waterfall.Year[k] > 0 {
			fd.WaterfallYearAll[key] = waterfallYearAll - value
			fd.WaterfallYearAdd[key] = value
			fd.WaterfallYearSub[key] = 0
		} else {
			fd.WaterfallYearAll[key] = waterfallYearAll + value
			fd.WaterfallYearAdd[key] = 0
			fd.WaterfallYearSub[key] = -value
		}

	}

	// 遍历 Waterfall.Month 生成瀑布图数据
	fd.WaterfallMonthAll = make(ChartData)
	fd.WaterfallMonthAdd = make(ChartData)
	fd.WaterfallMonthSub = make(ChartData)
	var waterfallMonthAll float64
	waterfallMonthAll = 0
	for _, k := range fd.Waterfall.Month.SortKey() {
		key := fmt.Sprintf("%s", k)
		value := fd.Waterfall.Month[k]

		if value == 0 {
			continue
		}

		waterfallMonthAll += value

		if fd.Waterfall.Month[k] > 0 {
			fd.WaterfallMonthAll[key] = waterfallMonthAll - value
			fd.WaterfallMonthAdd[key] = value
			fd.WaterfallMonthSub[key] = 0
		} else {
			fd.WaterfallMonthAll[key] = waterfallMonthAll + value
			fd.WaterfallMonthAdd[key] = 0
			fd.WaterfallMonthSub[key] = -value
		}
	}

	// 遍历消费和预算生成月度消费、预算数据
	fd.BudgetMonth = make(ChartData)
	fd.BudgetReal = make(ChartData)
	for _, v := range fd.Budgets {
		if v.Year == 0 || v.Month == 0 {
			continue
		}
		key := fmt.Sprintf("%d-%02d", v.Year, v.Month)
		fd.BudgetMonth[key] = v.Money
		fd.BudgetReal[key] = v.Real
	}

	fd.SpendYear = make(ChartData)
	for _, v := range fd.Spend.Year.SortKey() {
		if v == 0 {
			continue
		}
		key := fmt.Sprintf("%d", v)
		fd.SpendYear[key] = fd.Spend.Year[v]
	}
	fd.SpendMonth = make(ChartData)
	for _, v := range fd.Spend.Month.SortKey() {
		if v == "" {
			continue
		}
		key := fmt.Sprintf("%s", v)
		fd.SpendMonth[key] = fd.Spend.Month[v]
	}
	fd.SpendDay = make(ChartData)
	for _, v := range fd.Spend.Day.SortKey() {
		if v == "" {
			continue
		}
		key := fmt.Sprintf("%s", v)
		fd.SpendDay[key] = fd.Spend.Day[v]
	}

	fd.InvestmentYear = make(ChartData)
	for _, v := range fd.Investment.Year.SortKey() {
		if v == 0 {
			continue
		}
		key := fmt.Sprintf("%d", v)
		fd.InvestmentYear[key] = fd.Investment.Year[v]
	}
	fd.InvestmentMonth = make(ChartData)
	for _, v := range fd.Investment.Month.SortKey() {
		if v == "" {
			continue
		}
		key := fmt.Sprintf("%s", v)
		fd.InvestmentMonth[key] = fd.Investment.Month[v]
	}

}

func (fd *FullData) ShowChartData() {
	fmt.Println("====> Overview:")
	utils.PrettyPrint(fd.Overview)

	fmt.Println("====> OverviewType:")
	utils.PrettyPrint(fd.OverviewType)

	fmt.Println("====> WaterfallYearAll:")
	utils.PrettyPrint(fd.WaterfallYearAll)
	fmt.Println("====> WaterfallYearAdd:")
	utils.PrettyPrint(fd.WaterfallYearAdd)
	fmt.Println("====> WaterfallYearSub:")
	utils.PrettyPrint(fd.WaterfallYearSub)

	fmt.Println("====> WaterfallMonthAll:")
	utils.PrettyPrint(fd.WaterfallMonthAll)
	fmt.Println("====> WaterfallMonthAdd:")
	utils.PrettyPrint(fd.WaterfallMonthAdd)
	fmt.Println("====> WaterfallMonthSub:")
	utils.PrettyPrint(fd.WaterfallMonthSub)

	fmt.Println("====> BudgetMonth:")
	utils.PrettyPrint(fd.BudgetMonth)
	fmt.Println("====> BudgetReal:")
	utils.PrettyPrint(fd.BudgetReal)

	fmt.Println("====> SpendYear:")
	utils.PrettyPrint(fd.SpendYear)
	fmt.Println("====> SpendMonth:")
	utils.PrettyPrint(fd.SpendMonth)
	fmt.Println("====> SpendDay:")
	utils.PrettyPrint(fd.SpendDay)

	fmt.Println("====> InvestmentYear:")
	utils.PrettyPrint(fd.InvestmentYear)
	fmt.Println("====> InvestmentMonth:")
	utils.PrettyPrint(fd.InvestmentMonth)
}

func (fd *FullData) Report() {
	if !fd.hasStatistic {
		fd.StatisticAll()
	}

	fmt.Println()
	fmt.Println("======== BJ-PFD2 Report ========")
	fmt.Printf("资产总额：%s (%s)\n", utils.PrintRMB(fd.SumMoney), utils.PrintRMB(fd.SumCredit))
	fmt.Println("投资总额：", utils.PrintRMB(fd.SumInvestment))
	fmt.Println("投资收益总额：", utils.PrintRMB(fd.SumInvestmentIncome))
	fmt.Println("收支总况：", utils.PrintRMB(fd.SumBillsIncome), " / ",
		utils.PrintRMB(fd.SumBillsSpend))
	fmt.Println("本年收支：", utils.PrintRMB(fd.SumThisYearIncome), " / ",
		utils.PrintRMB(fd.SumThisYearSpend))
	fmt.Println("本月收支：", utils.PrintRMB(fd.SumThisMonthIncome), " / ",
		utils.PrintRMB(fd.SumThisMonthSpend))

	fmt.Println()

	fmt.Println("======== 资产瀑布统计 ========")
	fmt.Println(fd.WaterfallReport)
	fmt.Println("======== 预算消费统计 ========")
	fmt.Println(fd.BudgetsReport)
	fmt.Println("======== 投资账户统计 ========")
	fmt.Println(fd.IAccountsReport)
	fmt.Println("======== 账户统计 ========")
	fmt.Println(fd.AccountsReport)
	fmt.Println("======== 消费统计 ========")
	fmt.Println(fd.BillsReport)
}
