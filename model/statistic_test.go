package model

import "testing"

func TestStatisticSpend(t *testing.T) {
	iAccounts := Accounts{
		Account{PID: "1", Name: "1", Money: 0, IMoney: 0, IEarning: 0, RMoney: 0, Type: "1"},
		Account{PID: "2", Name: "2", Money: 0, IMoney: 0, IEarning: 0, RMoney: 0, Type: "1"},
	}
	iBills := Bills{
		Bill{PID: "1", Name: "1", Money: 100, Year: 2022, Month: 10, Day: 1,
			Trace: false, Account: "1", Budget: "1", Type: "1", UsageType: "1"},
		Bill{PID: "2", Name: "2", Money: 200, Year: 2022, Month: 10, Day: 1,
			Trace: false, Account: "2", Budget: "2", Type: "2", UsageType: "2"},
	}
	wantAccounts := Accounts{
		Account{PID: "1", Name: "1", Money: 100, IMoney: 0, IEarning: 0, RMoney: 0, Type: "1"},
		Account{PID: "2", Name: "2", Money: 200, IMoney: 0, IEarning: 0, RMoney: 0, Type: "1"},
	}
	if oas := StatisticSpend(&iAccounts, iBills); !oas.Compare(&wantAccounts) {
		t.Errorf("StatisticSpend() = %v, want %v", oas, wantAccounts)
	}
}

func TestStatisticInvestment(t *testing.T) {
	iAccounts := IAccounts{
		IAccount{PID: "1", Name: "1", Money: 0, Earning: 10, RMoney: 0, RAID: "1"},
		IAccount{PID: "2", Name: "2", Money: 0, Earning: -20, RMoney: 0, RAID: "2"},
	}
	iInvestments := Investments{
		Investment{PID: "1", Name: "1", Money: 100, Earning: 0, Year: 2022, Month: 10, Day: 1, Account: "1", Type: "1"},
		Investment{PID: "2", Name: "2", Money: 200, Earning: 0, Year: 2022, Month: 10, Day: 1, Account: "2", Type: "2"},
	}
	wantAccounts := IAccounts{
		IAccount{PID: "1", Name: "1", Money: 100, Earning: 10, RMoney: 110, RAID: "1"},
		IAccount{PID: "2", Name: "2", Money: 200, Earning: -20, RMoney: 180, RAID: "2"},
	}
	if oas := StatisticInvestment(&iAccounts, &iInvestments); !oas.Compare(&wantAccounts) {
		t.Errorf("StatisticInvestment() = %v, want %v", oas, wantAccounts)
	}
}

func TestStatisticAccountWithIAccount(t *testing.T) {
	iAccounts := Accounts{
		Account{PID: "1", Name: "1", Money: 200, IMoney: 0, IEarning: 0, RMoney: 0, Type: "1"},
		Account{PID: "2", Name: "2", Money: 200, IMoney: 0, IEarning: 0, RMoney: 0, Type: "1"},
	}
	iIAccounts := IAccounts{
		IAccount{PID: "1", Name: "1", Money: 100, Earning: 10, RMoney: 110, RAID: "1"},
		IAccount{PID: "2", Name: "2", Money: 200, Earning: -20, RMoney: 180, RAID: "2"},
	}
	wantAccounts := Accounts{
		Account{PID: "1", Name: "1", Money: 210, IMoney: 110, IEarning: 10, RMoney: 100, Type: "1"},
		Account{PID: "2", Name: "2", Money: 180, IMoney: 180, IEarning: -20, RMoney: 0, Type: "1"},
	}
	if oas := StatisticAccountWithIAccount(&iAccounts, &iIAccounts); !oas.Compare(&wantAccounts) {
		t.Errorf("StatisticAccountWithIAccount() = %v, want %v", oas, wantAccounts)
	}
}
