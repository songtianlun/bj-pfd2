package model

import (
	"bj-pfd2/com/utils"
	"fmt"
)

type FullData struct {
	SumMoney            float64 // 资产总额
	SumCredit           float64 // 信用总额
	SumInvestment       float64 // 投资总额
	SumInvestmentIncome float64 // 投资收益总额
	SumBillsSpend       float64 // 账单支出总额
	SumBillsIncome      float64 // 账单收入总额
	SumThisYearSpend    float64 // 本年支出总额
	SumThisYearIncome   float64 // 本年收入总额
	SumThisMonthSpend   float64 // 本月支出总额
	SumThisMonthIncome  float64 // 本月收入总额

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
