package main

import (
	"github.com/shopspring/decimal"
)

/**
  *@Description 等本等息：在还款计划中，各期次应还款额相等，且各期次应还本金相等，各期次应还利息相等。
  *@Author pauline
  *@Date 2023/12/5 10:19
**/
func equalPrincipalAndInterestPlan(request repayPlanRequest, response *Response) error {
	var sumTotalInterest, hasRepayPrincipal, sumTotalRepayAmount decimal.Decimal

	records := make([]RepayPlanRecord, 0)

	planRepayPrinciplePeriod := request.LoanAmount.Div(decimal.NewFromInt(int64(request.TotalPeriodNum))).Round(2)

	totalCalcDays := getDaysBetweenDate(request.LoanStartDateParseLocal, request.LoanEndDateParseLocal)

	// 2.3 everyMonth need to repay interest amount = LoanAmount*daysRate*totalDays/periodNum
	planRepayInterestPeriod := request.LoanAmount.Mul(request.DaysInterestRate).
		Mul(decimal.NewFromInt(totalCalcDays)).Div(decimal.NewFromInt(int64(request.TotalPeriodNum))).Round(2)

	dateMap := calculatePeriodDate(request)

	for i := 0; i < request.TotalPeriodNum; i++ {
		periodStartDate := dateMap[i][0]
		periodEndDate := dateMap[i][1]
		periodRepayDate := dateMap[i][2]

		daysOfPeriod := getDaysBetweenDate(periodStartDate, periodEndDate)

		record := RepayPlanRecord{
			PeriodNum:           i + 1,                                    // 当前期次的期数
			PeriodStartDate:     periodStartDate.Format(DATE_DASH_FORMAT), // 当前期次的开始计息日
			PeriodEndDate:       periodEndDate.Format(DATE_DASH_FORMAT),   // 当前期次的结束计息日
			PeriodRepayDate:     periodRepayDate.Format(DATE_DASH_FORMAT), // 当前期次的还款日
			DaysOfPeriod:        int(daysOfPeriod),
			PeriodRepayInterest: planRepayInterestPeriod,
		}

		if i == request.TotalPeriodNum-1 {
			remainPrinciple := request.LoanAmount.Sub(hasRepayPrincipal)
			hasRepayPrincipal = hasRepayPrincipal.Add(remainPrinciple)
			record.PeriodRepayPrinciple = remainPrinciple
			record.PeriodRepayTotalAmount = planRepayInterestPeriod.Add(remainPrinciple)
			record.MaintainPrinciple = request.LoanAmount.Sub(hasRepayPrincipal)
		} else {
			hasRepayPrincipal = hasRepayPrincipal.Add(planRepayPrinciplePeriod)
			record.PeriodRepayPrinciple = planRepayPrinciplePeriod
			record.PeriodRepayTotalAmount = planRepayInterestPeriod.Add(planRepayPrinciplePeriod)
		}

		record.MaintainPrinciple = request.LoanAmount.Sub(hasRepayPrincipal) // 剩余还款本金

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
