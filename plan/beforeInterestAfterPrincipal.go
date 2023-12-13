package main

import (
	"github.com/shopspring/decimal"
)

/**
  *@Description 先息后本：分期付息：除最后一期之外的前几期，每期仅需归还当期利息。到期还本：贷款到期日（末期）全额还本及相应的利息
  *@Author pauline
  *@Date 2023/12/5 10:19
**/
func beforeInterestAfterPrincipalPlan(request repayPlanRequest, response *Response) error {
	var sumTotalInterest, sumTotalRepayAmount decimal.Decimal

	records := make([]RepayPlanRecord, 0)
	dateMap := calculatePeriodDate(request)

	// 2.make repay plan
	for i := 0; i < request.TotalPeriodNum; i++ {
		periodStartDate := dateMap[i][0]
		periodEndDate := dateMap[i][1]
		periodRepayDate := dateMap[i][2]

		// 当前期次的计息天数
		daysOfPeriod := getDaysBetweenDate(periodStartDate, periodEndDate)

		// 当前期次的利息金额=贷款本金*计息天数*日利息
		periodRepayInterest := request.LoanAmount.Mul(request.DaysInterestRate).Mul(decimal.NewFromInt(daysOfPeriod)).Round(2)

		record := RepayPlanRecord{
			PeriodNum:           i + 1,
			PeriodStartDate:     periodStartDate.Format(DATE_DASH_FORMAT), // 当前期次的开始计息日
			PeriodEndDate:       periodEndDate.Format(DATE_DASH_FORMAT),   // 当前期次的结束计息日
			PeriodRepayDate:     periodRepayDate.Format(DATE_DASH_FORMAT), // 当前期次的还款日
			DaysOfPeriod:        int(daysOfPeriod),                        // 当前期次的计息天数
			PeriodRepayInterest: periodRepayInterest,                      // 当前期次的利息
		}

		if i == request.TotalPeriodNum-1 {
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
	response.PlanRepayRecords = records
	response.TotalRepayAmount = sumTotalRepayAmount
	response.TotalInterest = sumTotalInterest

	return nil
}
