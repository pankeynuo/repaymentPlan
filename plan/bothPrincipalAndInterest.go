package main

import (
	"github.com/shopspring/decimal"
	"time"
)

/**
  *@Description 利(息)随本清：在贷款的到期日，归还贷款全额本金及全部利息。无需分期归还贷款本息
  *@Author pauline
  *@Date 2023/12/5 10:18
**/
func bothPrincipalAndInterest(request *Request) (response *Response, err error) {
	loanStartDateParseLocal, _ := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanStartDate, time.Local)
	loanEndDateParseLocal, _ := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanEndDate, time.Local)

	// 1.Daily interest rate 日利率
	daysInterestRate := calculateDaysInterestRate(request.InterestRate, request.DaysOfYear)

	reCalEndDate := loanEndDateParseLocal.AddDate(0, 0, -1)

	// 2.get this period's interest begin 计息天数
	daysOfPeriod := getDaysBetweenDate(loanStartDateParseLocal, reCalEndDate)

	// 3.calculate total interest amount
	totalInterest := request.LoanAmount.Mul(daysInterestRate).Mul(decimal.NewFromInt(daysOfPeriod)).Round(2)

	// 4.response only one period's plan
	totalAmount := totalInterest.Add(request.LoanAmount)
	record := RepayPlanRecord{
		PeriodNum:              1,
		PeriodStartDate:        request.LoanStartDate,
		PeriodEndDate:          reCalEndDate.Format(DATE_DASH_FORMAT),
		PeriodRepayDate:        request.LoanEndDate,
		PeriodRepayPrinciple:   request.LoanAmount,
		MaintainPrinciple:      decimal.NewFromFloat(0),
		DaysOfPeriod:           int(daysOfPeriod),
		PeriodRepayInterest:    totalInterest,
		PeriodRepayTotalAmount: totalAmount,
	}
	response = &Response{
		TotalPeriodNum:   1,
		RepayMethod:      request.RepayMethod,
		LoanStartDate:    request.LoanStartDate,
		LoanEndDate:      request.LoanEndDate,
		LoanAmount:       request.LoanAmount,
		InterestRate:     request.InterestRate,
		TotalRepayAmount: totalAmount,
		TotalInterest:    totalInterest,
		PlanRepayRecords: []RepayPlanRecord{record},
	}
	return response, nil
}
