package main

import (
	"github.com/shopspring/decimal"
)

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/5 10:19
**/
func beforeInterestAfterPrincipal(request *Request) (response *Response, err error) {
	var sumTotalInterest, sumTotalRepayAmount decimal.Decimal

	records := make([]RepayPlanRecord, 0)
	dateMap := calculatePeriodDate(request)

	// 2.make repay plan
	for i := 0; i < request.PeriodNum; i++ {
		periodStartDate := dateMap[i][0]
		periodEndDate := dateMap[i][1]
		periodRepayDate := dateMap[i][2]
		// reCalStartDate := ""
		// reCalEndDate := ""
		// repaymentDate := ""
		// 当前期次的计息天数
		daysOfPeriod := getDaysBetweenDate(periodStartDate, periodEndDate)

		// 当前期次的利息金额=贷款本金*计息天数*日利息
		periodRepayInterest := request.LoanAmount.Mul(rateD).Mul(decimal.NewFromInt(daysOfPeriod)).Round(2)

		record := RepayPlanRecord{
			PeriodNum:           i + 1,
			PeriodStartDate:     periodStartDate.Format(DATE_DASH_FORMAT), // 当前期次的开始计息日
			PeriodEndDate:       periodEndDate.Format(DATE_DASH_FORMAT),   // 当前期次的结束计息日
			RepayDate:           periodRepayDate.Format(DATE_DASH_FORMAT), // 当前期次的还款日
			DaysOfPeriod:        int(daysOfPeriod),                        // 当前期次的计息天数
			PeriodRepayInterest: periodRepayInterest,                      // 当前期次的利息
		}

		if i == request.PeriodNum-1 {
			record.PeriodRepayPrinciple = request.LoanAmount
			record.PeriodRepayTotalAmount = periodRepayInterest.Add(request.LoanAmount)
			record.MaintainPrinciple = decimal.NewFromInt(0)
		} else {
			record.PeriodRepayPrinciple = decimal.NewFromFloat(0)
			record.PeriodRepayTotalAmount = periodRepayInterest
			record.MaintainPrinciple = request.LoanAmount
		}
		// 累积还款总金额
		sumTotalRepayAmount = sumTotalRepayAmount.Add(record.PeriodRepayTotalAmount)
		// 累积还款总利息
		sumTotalInterest = sumTotalInterest.Add(record.PeriodRepayInterest)

		records = append(records, record)
	}
	response = &Response{
		PlanRepayRecords: records,
		RepayMethod:      request.RepayMethod,
		LoanStartDate:    request.LoanStartDate,
		LoanEndDate:      request.LoanEndDate,
		TotalPeriodNum:   request.PeriodNum,
		LoanAmount:       request.LoanAmount,
		InterestRate:     request.InterestRate,
		TotalRepayAmount: sumTotalRepayAmount,
		TotalInterest:    sumTotalInterest,
	}
	return nil, nil
}
