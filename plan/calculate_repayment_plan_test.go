package main

import (
	"encoding/json"
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
		LoanAmount:    decimal.NewFromFloat(400000),
		LoanStartDate: "2022-01-01",
		LoanEndDate:   "",
		InterestRate:  decimal.NewFromFloat(6),
		PeriodNum:     5,
		RepayDay:      10,
		LoanCycleCode: "03",
		RepayMethod:   "02", // 01-等额本息 02-等额本金
		PeriodType:    "02", // 01-年 02-月
		DaysOfYear:    360,
		DaysOfMonth:   30,
	})
	resp2, _ := json.Marshal(resp)
	t.Log(string(resp2))
	t.Log(err)
}
