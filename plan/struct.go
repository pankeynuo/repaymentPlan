package main

import (
	"github.com/shopspring/decimal"
	"time"
)

type Request struct {
	LoanAmount    decimal.Decimal `json:"loanAmount" validate:"required"`    // 贷款金额
	LoanStartDate string          `json:"loanStartDate" validate:"required"` // 利息计算开始日期
	LoanEndDate   string          `json:"loanEndDate"`                       // 利息计算结束日期
	LoanCycleCode string          `json:"loanCycleCode"`                     // 还款周期频率 01-日 02-两周 03-月 04-季 05-年
	InterestRate  decimal.Decimal `json:"interestRate" validate:"required"`  // 年利率
	RepayMethod   string          `json:"repayMethod" validate:"required"`   // 还款方式:1-等额本息  2-等额本金  3-利随本清 4-先息后本 5-等本等息
	PeriodNum     int             `json:"periodNum"`                         // 期数
	PeriodType    string          `json:"periodType"`                        // 期数类型 01-年 02-月
	RepayDay      int             `json:"repayDay"`                          // 每一期还款日
	DaysOfYear    int             `json:"daysOfYear"`                        // 年天数 默认360
}

type Response struct {
	RepayMethod      string            `json:"repayMethod"`            // 还款方式:1-等额本息  2-等额本金  3-利随本清 4-先息后本 5-等本等息
	LoanStartDate    string            `json:"loanStartDate"`          // 利息计算开始日期
	LoanEndDate      string            `json:"loanEndDate"`            // 利息计算结束日期
	TotalPeriodNum   int               `json:"totalPeriodNum"`         // 期数
	TotalRepayAmount decimal.Decimal   `json:"totalRepayAmount"`       // 总还款金额
	LoanAmount       decimal.Decimal   `json:"loanAmount"`             // 贷款金额
	TotalInterest    decimal.Decimal   `json:"planRepayTotalInterest"` // 总还款利息
	InterestRate     decimal.Decimal   `json:"interestRate"`           // 年利率
	PlanRepayRecords []RepayPlanRecord `json:"planRepayRecords"`       // 还款计划
}
type RepayPlanRecord struct {
	PeriodNum              int             `json:"periodNum"`              // 期次
	PeriodStartDate        string          `json:"periodStartDate"`        // 本期开始日期
	PeriodEndDate          string          `json:"periodEndDate"`          // 本期结束日期
	DaysOfPeriod           int             `json:"daysOfPeriod"`           // 本期天数
	PeriodRepayDate        string          `json:"periodRepayDate"`        // 本期还款日期
	PeriodRepayTotalAmount decimal.Decimal `json:"periodRepayTotalAmount"` // 本期还款总金额
	PeriodRepayPrinciple   decimal.Decimal `json:"periodRepayPrinciple"`   // 本期还款本金
	PeriodRepayInterest    decimal.Decimal `json:"periodRepayInterest"`    // 本期还款利息
	MaintainPrinciple      decimal.Decimal `json:"maintainPrinciple"`      // 剩余还款金额
}

type repayPlanRequest struct {
	LoanAmount              decimal.Decimal // 贷款金额
	LoanStartDate           string          // 利息计算开始日期=开始贷款日期
	LoanEndDate             string          // 利息计算结束日期=最后一次还款日
	LoanCycleCode           string          // 还款周期频率 01-daily 日 02-fortnightly 两周 03-monthly 月 04-quarterly 季 05-yearly 年
	PeriodInterestRate      decimal.Decimal // 期利率
	TotalPeriodNum          int             // 总期数
	FirstRepayDate          time.Time       // 首个还款日
	LoanStartDateParseLocal time.Time
	LoanEndDateParseLocal   time.Time
	RepayDay                int
	DaysInterestRate        decimal.Decimal
}
