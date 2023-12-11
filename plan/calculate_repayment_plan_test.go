package main

import (
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/4 11:17
**/
func Test_FirstRepayDate(t *testing.T) {
	a := time.Now()
	for i := 0; i < 20; i++ {
		_, _ = CalculateRepaymentPlan(&Request{
			LoanAmount:    decimal.NewFromFloat(400000),
			LoanStartDate: "2022-01-01",
			LoanEndDate:   "",
			InterestRate:  decimal.NewFromFloat(6),
			PeriodNum:     5,
			RepayDay:      10,
			LoanCycleCode: "03",
			RepayMethod:   "01", // 01-等额本息
			PeriodType:    "02", // 01-年 02-月
			DaysOfYear:    360,
			DaysOfMonth:   30,
		})
	}
	b := time.Now()
	t.Log(b.Sub(a))
}
