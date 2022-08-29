package model

import (
	"bj-pfd2/com/utils"
	"fmt"
)

type ChartData map[string]float64

type FullData struct {
	Token string

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

	WaterfallDayAll ChartData // 流水日报表-总额
	WaterfallDayAdd ChartData // 流水日报表-增量
	WaterfallDaySub ChartData // 流水日报表-减量

	SpendYear   ChartData // 支出年度报表
	SpendMonth  ChartData // 支出月度报表
	BudgetMonth ChartData // 预算月度报表
	SpendDay    ChartData // 支出日报表

	InvestmentYear  ChartData // 投资年度报表
	InvestmentMonth ChartData // 投资月度报表

	Accounts    Accounts
	Investments Investments
	Bills       Bills
	IAccounts   IAccounts
	Budgets     Budgets
	Waterfall   Waterfall

	AccountsReport    string
	BillsReport       string
	InvestmentsReport string
	IAccountsReport   string
	BudgetsReport     string
	WaterfallReport   string

	hasStatistic bool // 是否已经统计过
}

func (fd *FullData) GenerateStrRMB() {
	fd.SumMoneyStr = utils.Float64ToRMB(fd.SumMoney)
	fd.SumCreditStr = utils.Float64ToRMB(fd.SumCredit)
	fd.SumInvestmentStr = utils.Float64ToRMB(fd.SumInvestment)
	fd.SumInvestmentIncomeStr = utils.Float64ToRMB(fd.SumInvestmentIncome)
	fd.SumBillsSpendStr = utils.Float64ToRMB(fd.SumBillsSpend)
	fd.SumBillsIncomeStr = utils.Float64ToRMB(fd.SumBillsIncome)
	fd.SumThisYearSpendStr = utils.Float64ToRMB(fd.SumThisYearSpend)
	fd.SumThisYearIncomeStr = utils.Float64ToRMB(fd.SumThisYearIncome)
	fd.SumThisMonthSpendStr = utils.Float64ToRMB(fd.SumThisMonthSpend)
	fd.SumThisMonthIncomeStr = utils.Float64ToRMB(fd.SumThisMonthIncome)
}

func (fd *FullData) StatisticAll() {
	if fd.hasStatistic {
		return
	}

	fd.Accounts = *StatisticSpend(&fd.Accounts, fd.Bills)                    // Step1: 统计账单支出
	fd.IAccounts = *StatisticInvestment(&fd.IAccounts, &fd.Investments)      // Step2： 然后统计投资状况
	fd.Accounts = *StatisticAccountWithIAccount(&fd.Accounts, &fd.IAccounts) // Step3: 将投资状况统计到账户中
	fd.Budgets = *StatisticBillsWithBudget(&fd.Bills, &fd.Budgets)           // Step4: 统计账单支出到预算中
	fd.Budgets.StatisticRemain()                                             // Step5: 统计预算剩余金额
	fd.Waterfall = *fd.Bills.Waterfall()                                     // Step6: 根据账单支出统计流水

	fd.hasStatistic = true // Step7: 标记已统计

	// Step8: 统计报表
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

	fd.Report()
	fd.ShowChartData()

	//fmt.Println(fd.Accounts.GenerateReport())
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
