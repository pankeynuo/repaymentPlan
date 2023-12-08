package main

import (
	"github.com/shopspring/decimal"
	"time"
)

type Request struct {
	LoanAmount    decimal.Decimal `json:"loanAmount" validate:"required"`     // 贷款金额
	LoanStartDate string          `json:"loanStartDate" validate:"required"`  // 利息计算开始日期=开始贷款日期
	LoanEndDate   string          `json:"loanEndDate"`                        // 利息计算结束日期=最后一次还款日
	LoanCycleCode string          `json:"loanCycleCode"  validate:"required"` // 还款周期频率 01-daily 日 02-fortnightly 两周 03-monthly 月 04-quarterly 季 05-yearly 年
	InterestRate  decimal.Decimal `json:"interestRate" validate:"required"`   // 年利率
	RepayMethod   string          `json:"repayMethod" validate:"required"`    // 还款方式:01-等额本息  02-等额本金  03-到期（一次性）还本付息（息随本清)(到期一次性还本还息）  04-到期还本周期还息(分期付息到期还本（先息后本)) 05-等本等息（每期还本还息还款额都相等，每期计息的本金为贷款总本金）
	PeriodNum     int             `json:"periodNum"`                          // 期数
	PeriodType    string          `json:"periodType"`                         // 期数类型 01-年 02-月
	RepayDay      int             `json:"repayDay" validate:"required"`       // 每一期还款日
	DaysOfYear    int             `json:"daysOfYear"`                         // 年天数 默认360
	DaysOfMonth   int             `json:"daysOfMonth"`                        // 月天数 默认30

	// NextRepayDate string        `json:"nextRepayDate"` // 下个还款日
	// BillDay       string        `json:"billDay"`       // 账单日
	// EditFlag      string        `json:"editFlag" description:"edit repayment flag Y-Yes "`
	// IsSameFlag    bool          `json:"isSameFlag"`
	// SkipReduce    SkipAndReduce `json:"skipReduce"`
	// ListFee       []FeeOnTop    `json:"listFee"`
}
type SkipAndReduce struct {
	CurrentPeriodBeginDate string          `json:"currentPeriodBeginDate" description:"current period begin date"`
	CurrentPeriod          int             `json:"currentPeriod" description:"current period"`
	Installment            decimal.Decimal `json:"installment"`
	ReducePercentage       decimal.Decimal `json:"reducePercentage"`
	ReducePeriodCount      int             `json:"reducePeriodCount"`
	SkipPeriodCount        int             `json:"skipPeriodCount"`
}
type Response struct {
	RepayMethod      string            `json:"repayMethod"`            // 还款方式:01-等额本息  02-等额本金  03-到期（一次性）还本付息（息随本清)(到期一次性还本还息）  04-到期还本周期还息(分期付息到期还本（先息后本)) 05-等本等息（每期还本还息还款额都相等，每期计息的本金为贷款总本金）
	LoanStartDate    string            `json:"loanStartDate"`          // 利息计算开始日期=开始贷款日期
	LoanEndDate      string            `json:"loanEndDate"`            // 利息计算结束日期=最后一次还款日
	TotalPeriodNum   int               `json:"totalPeriodNum"`         // 期数
	TotalRepayAmount decimal.Decimal   `json:"totalRepayAmount"`       // 总还款金额
	LoanAmount       decimal.Decimal   `json:"loanAmount"`             // 贷款金额
	TotalInterest    decimal.Decimal   `json:"planRepayTotalInterest"` // 总还款利息
	InterestRate     decimal.Decimal   `json:"interestRate"`           // 年利率
	PlanRepayRecords []RepayPlanRecord `json:"planRepayRecords"`       // 还款计划
	// ListOnTopRecords []OnTopRecords    `json:"listOnTopRecords"`
}
type RepayPlanRecord struct {
	PeriodNum              int             `json:"periodNum"`              // 期次
	PeriodStartDate        string          `json:"periodStartDate"`        // 本期开始日期
	PeriodEndDate          string          `json:"periodEndDate"`          // 本期结束日期
	DaysOfPeriod           int             `json:"daysOfPeriod"`           // 本期天数
	RepayDate              string          `json:"repayDate"`              // 本期还款日期
	PeriodRepayTotalAmount decimal.Decimal `json:"periodRepayTotalAmount"` // 本期还款总金额
	PeriodRepayPrinciple   decimal.Decimal `json:"periodRepayPrinciple"`   // 本期还款本金
	PeriodRepayInterest    decimal.Decimal `json:"periodRepayInterest"`    // 本期还款利息
	MaintainPrinciple      decimal.Decimal `json:"maintainPrinciple"`      // 剩余还款金额
	// DaysOfInterestCalculate    int             `json:"daysOfInterestCalculate"`
}

type FeeOnTop struct {
	FeeType     string          `json:"feeType" description:"fee type"`
	FeeAmount   decimal.Decimal `json:"feeAmount" description:"fee amount"`
	OnTop       string          `json:"onTop" description:"on top"`
	OnTopPeriod int             `json:"onTopPeriod" description:"on top period"`
}
type OnTopRecords struct {
	FeeType              string              `json:"feeType" description:"fee type"`
	PlanRepayTotalAmount decimal.Decimal     `json:"planRepayTotalAmount" description:"plan repay total amount"`
	OnTopRepayPlans      []FeeOnTopRepayPlan `json:"onTopRepayPlans" description:"on top repay plans"`
}

type FeeOnTopRepayPlan struct {
	PeriodNum                  int             `json:"periodNum"`
	InterestCalculateStartDate string          `json:"interestCalculateStartDate" description:"interest calculate start date"`
	InterestCalculateEndDate   string          `json:"interestCalculateEndDate" description:"interest calculate end date"`
	PlanRepayDate              string          `json:"planRepayDate" description:"plan repay date"`
	PlanRepayTotalAmount       decimal.Decimal `json:"planRepayTotalAmount"`
	PlanRepayPrinciple         decimal.Decimal `json:"planRepayPrinciple" description:"plan repay principle"`
	PlanRepayInterest          decimal.Decimal `json:"planRepayInterest"`
	MaintainPrinciple          decimal.Decimal `json:"maintainPrinciple" description:"maintain principle"`
}
type repayPlanRequest struct {
	LoanAmount              decimal.Decimal // 贷款金额
	LoanStartDate           string          // 利息计算开始日期=开始贷款日期
	LoanEndDate             string          // 利息计算结束日期=最后一次还款日
	LoanCycleCode           string          // 还款周期频率 01-daily 日 02-fortnightly 两周 03-monthly 月 04-quarterly 季 05-yearly 年
	PeriodInterestRate      decimal.Decimal // 期利率
	TotalPeriodNum          int             // 总期数
	FirstRepayDate          string          // 首个还款日
	LoanStartDateParseLocal time.Time
	LoanEndDateParseLocal   time.Time
}
