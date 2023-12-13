package main

// 还款方式
const (
	EqualLoanRepayment           = "1" // 等额本息
	EqualPrincipalRepayment      = "2" // 等额本金
	BothPrincipalAndInterest     = "3" // 息随本清
	BeforeInterestAfterPrincipal = "4" // 先息后本
	EqualPrincipalAndInterest    = "5" // 等本等息
)

const (
	DATE_DASH_FORMAT = "2006-01-02"
)

const (
	daysOfYear    = 360
	numberOfWeek  = 26
	numberOfMonth = 12
)

const (
	loanCycleDaily       = "01" // 日
	loanCycleFortnightly = "02" // 两周
	loanCycleMonthly     = "03" // 月
	loanCycleQuarterly   = "04" // 季
	loanCycleYearly      = "05" // 年
)

const (
	periodTypeYear  = "01"
	periodTypeMonth = "02"
)
