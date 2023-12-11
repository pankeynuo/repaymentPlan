package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/4 11:17
**/
func Test_FirstRepayDate(t *testing.T) {
	resp, err := CalculateRepaymentPlan(&Request{
		LoanAmount:    decimal.NewFromFloat(10000),
		LoanStartDate: "2023-01-01",
		LoanEndDate:   "",
		LoanCycleCode: "03",
		InterestRate:  decimal.NewFromFloat(2),
		RepayMethod:   "01", // 01-等额本息
		PeriodNum:     10,
		PeriodType:    "02", // 01-年 02-月
		RepayDay:      30,
		DaysOfYear:    360,
		DaysOfMonth:   30,
	})
	fmt.Printf("resp:%v,error:%v", resp, err)
}

// 2023-01-20 还款日10 2023-02-10
// 2023-01-10 还款日30 2023-01-30
// 2023-01-31 还款日30 2023-02-28/29
// 2023-01-30 还款日30 2023-02-28/29
// 2023-02-28 还款日30 2023-03-30   下个月不是2月，每个月的30， 下个月是2月，2月的最后一天
// 2023-03-30 还款日31 2023-04-30   月底，下个月的最后一天

// 月底30 但是本月已经最后一天 28/29/30，下个月的30
// 月底31 但是本月已经最后一天 28/29/30，下个月的最后一天

// 还款日大于今天
// 拼接这个月的还款日 日期不存在 为最后一天
// 还款日小于等于今天
// 拼接下个月的还款日。如果日期不存在，为最后一天

// 如果小于20天 在加一个月
