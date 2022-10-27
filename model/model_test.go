package model

import (
    "bj-pfd2/model/chart"
    "testing"
)

var testRawFullData = FullData{
    Token:                  "",
    SumMoney:               0,
    SumMoneyStr:            "",
    SumCredit:              0,
    SumCreditStr:           "",
    SumInvestment:          0,
    SumInvestmentStr:       "",
    SumInvestmentIncome:    0,
    SumInvestmentIncomeStr: "",
    SumBillsSpend:          0,
    SumBillsSpendStr:       "",
    SumBillsIncome:         0,
    SumBillsIncomeStr:      "",
    SumThisYearSpend:       0,
    SumThisYearSpendStr:    "",
    SumThisYearIncome:      0,
    SumThisYearIncomeStr:   "",
    SumThisMonthSpend:      0,
    SumThisMonthSpendStr:   "",
    SumThisMonthIncome:     0,
    SumThisMonthIncomeStr:  "",
    OverviewType:           nil,
    Overview:               nil,
    WaterfallYearAll:       nil,
    WaterfallYearAdd:       nil,
    WaterfallYearSub:       nil,
    WaterfallMonthAll:      nil,
    WaterfallMonthAdd:      nil,
    WaterfallMonthSub:      nil,
    BudgetMonth:            nil,
    SpendYear:              nil,
    SpendMonth:             nil,
    SpendDay:               nil,
    InvestmentYear:         nil,
    InvestmentMonth:        nil,

    // 账户数据，记录账户名称、账户类型、
    // 统计需得到，账户余额、投资总额、收益总额、活期总额
    Accounts: Accounts{
        Account{PID: "0", Name: "测试账户1", Money: 0, IMoney: 0, IEarning: 0, RMoney: 0, Type: ""},
        Account{PID: "1", Name: "测试账户2", Money: 0, IMoney: 0, IEarning: 0, RMoney: 0, Type: ""},
    },
    Bills: Bills{
        Bill{
            PID: "0", Name: "测试流水1", Money: 1000, Year: 2020, Month: 10, Day: 10,
            IsTrace: false, Account: "0", Budget: "0", Type: "个人储蓄", UsageType: "工资",
        },
        Bill{
            PID: "0", Name: "测试流水1", Money: -100, Year: 2020, Month: 10, Day: 10,
            IsTrace: false, Account: "0", Budget: "0", Type: "个人储蓄", UsageType: "日食",
        },
        Bill{
            PID: "1", Name: "测试流水2", Money: -200, Year: 2020, Month: 10, Day: 10,
            IsTrace: false, Account: "0", Budget: "0", Type: "个人储蓄", UsageType: "日食",
        },
    },
    IAccounts:   IAccounts{},
    Investments: Investments{},
    Budgets:     Budgets{},
    Waterfall:   chart.Waterfall{},

    AccountsReport:    "",
    BillsReport:       "",
    InvestmentsReport: "",
    IAccountsReport:   "",
    BudgetsReport:     "",
    WaterfallReport:   "",
    hasStatistic:      false,
}

var teatResult = FullData{
    Token:                  "",
    SumMoney:               0,
    SumMoneyStr:            "",
    SumCredit:              0,
    SumCreditStr:           "",
    SumInvestment:          0,
    SumInvestmentStr:       "",
    SumInvestmentIncome:    0,
    SumInvestmentIncomeStr: "",
    SumBillsSpend:          0,
    SumBillsSpendStr:       "",
    SumBillsIncome:         0,
    SumBillsIncomeStr:      "",
    SumThisYearSpend:       0,
    SumThisYearSpendStr:    "",
    SumThisYearIncome:      0,
    SumThisYearIncomeStr:   "",
    SumThisMonthSpend:      0,
    SumThisMonthSpendStr:   "",
    SumThisMonthIncome:     0,
    SumThisMonthIncomeStr:  "",
    OverviewType:           nil,
    Overview:               nil,
    WaterfallYearAll:       nil,
    WaterfallYearAdd:       nil,
    WaterfallYearSub:       nil,
    WaterfallMonthAll:      nil,
    WaterfallMonthAdd:      nil,
    WaterfallMonthSub:      nil,
    BudgetMonth:            nil,
    SpendYear:              nil,
    SpendMonth:             nil,
    SpendDay:               nil,
    InvestmentYear:         nil,
    InvestmentMonth:        nil,

    Accounts: Accounts{
        Account{PID: "0", Name: "测试账户1", Money: 700, IMoney: 0, IEarning: 0, RMoney: 700, Type: ""},
        Account{PID: "1", Name: "测试账户2", Money: 0, IMoney: 0, IEarning: 0, RMoney: 0, Type: ""},
    },
    Bills: Bills{
        Bill{
            PID: "0", Name: "测试流水1", Money: 1000, Year: 2020, Month: 10, Day: 10,
            IsTrace: false, Account: "0", Budget: "0", Type: "个人储蓄", UsageType: "工资",
        },
        Bill{
            PID: "0", Name: "测试流水1", Money: -100, Year: 2020, Month: 10, Day: 10,
            IsTrace: false, Account: "0", Budget: "0", Type: "个人储蓄", UsageType: "日食",
        },
        Bill{
            PID: "1", Name: "测试流水2", Money: -200, Year: 2020, Month: 10, Day: 10,
            IsTrace: false, Account: "0", Budget: "0", Type: "个人储蓄", UsageType: "日食",
        },
    },
    IAccounts:   IAccounts{},
    Investments: Investments{},
    Budgets:     Budgets{},
    Waterfall:   chart.Waterfall{},

    AccountsReport:    "",
    BillsReport:       "",
    InvestmentsReport: "",
    IAccountsReport:   "",
    BudgetsReport:     "",
    WaterfallReport:   "",
    hasStatistic:      false,
}

func TestFullData_StatisticAll(t *testing.T) {
    // 测试数据
    testRawFullData.StatisticAll()
    // 比较结果
    if !testRawFullData.Compare(&teatResult) {
        t.Errorf("TestFullData_StatisticAll failed, got %v, want %v", testRawFullData, teatResult)
    }
}
