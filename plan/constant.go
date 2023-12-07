package main

// 还款方式
const (
	EqualLoanRepayment           = "01" // 等额本息
	EqualPrincipalRepayment      = "02" // 等额本金
	BothPrincipalAndInterest     = "03" // 息随本清:到期(一次性)还本付息(息随本清)
	BeforeInterestAfterPrincipal = "04" // 先息后本
	EqualPrincipalAndInterest    = "05" // 等本等息
	BillDateRepayment            = "06" // 账单日还款
)

const (
	DATE_DASH_FORMAT = "2006-01-02"
)

// 还款周期
const (
	DoubleWeek = "03"
	Monthly    = "04"
)

const (
	daysOfYear  = 360
	daysOfMonth = 30
)

const (
	loanCycleDaily       = "01" // 日
	loanCycleFortnightly = "02" // 两周
	loanCycleMonthly     = "03" // 月
	loanCycleQuarterly   = "04" //季
	loanCycleYearly      = "05" //年
)

const (
	periodTypeYear  = "01"
	periodTypeMonth = "02"
)
