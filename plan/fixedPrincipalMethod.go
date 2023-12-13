package main

import (
	"github.com/shopspring/decimal"
)

/**
  *@Description 等额本金：贷款利率保持不变的前提条件下，每期还贷本金均等，每期还贷利息随着贷款本金的逐步减少而递减。
  *@Author pauline
  *@Date 2023/12/5 10:18
**/
func fixedPrincipalMethodPlan(request repayPlanRequest, response *Response) error {
	var sumTotalInterest, hasRepayPrincipal, sumTotalRepayAmount decimal.Decimal
	records := make([]RepayPlanRecord, 0)

	// 每期还款本金
	periodRepayPrinciple := calculateFixedPrincipalMethod(request.LoanAmount, request.TotalPeriodNum)

	dateMap := calculatePeriodDate(request)
	for i := 0; i < request.TotalPeriodNum; i++ {
		periodStartDate := dateMap[i][0]
		periodEndDate := dateMap[i][1]
		periodRepayDate := dateMap[i][2]

		// 当前期次的计息天数
		daysOfPeriod := getDaysBetweenDate(periodStartDate, periodEndDate)

		periodRepayInterest := (request.LoanAmount.Sub(hasRepayPrincipal)).Mul(request.DaysInterestRate).Mul(decimal.NewFromInt(daysOfPeriod)).RoundBank(2)

		record := RepayPlanRecord{
			PeriodNum:           i + 1,                                    // current period num 当前期次的期数
			PeriodStartDate:     periodStartDate.Format(DATE_DASH_FORMAT), // 当前期次的开始计息日
			PeriodEndDate:       periodEndDate.Format(DATE_DASH_FORMAT),   // 当前期次的结束计息日
			PeriodRepayDate:     periodRepayDate.Format(DATE_DASH_FORMAT), // 当前期次的还款日
			DaysOfPeriod:        int(daysOfPeriod),                        // 当前期次的计息天数
			PeriodRepayInterest: periodRepayInterest,                      // 当前期次的利息
		}

		// if this is the last period 如果是最后一期
		if i == request.TotalPeriodNum-1 {
			remainPrinciple := request.LoanAmount.Sub(hasRepayPrincipal)
			record.PeriodRepayPrinciple = remainPrinciple.RoundBank(2)               // 当前期次还款本金=上一期总的剩余还款本金
			record.PeriodRepayTotalAmount = periodRepayInterest.Add(remainPrinciple) // 当前期次的总还款金额=当前期次的利息金额+当前期次还款本金
			hasRepayPrincipal = hasRepayPrincipal.Add(remainPrinciple)               // 累积已还本金=累积已还本金+上一期总的剩余还款本金
		} else {
			// if this not the last period 非最后一期
			hasRepayPrincipal = hasRepayPrincipal.Add(periodRepayPrinciple) // 累积已还本金
			record.PeriodRepayPrinciple = periodRepayPrinciple              // 当前期次还款本金
			// 当前期次的总还款金额=当前期次还款本金+当前期次的利息金额
			record.PeriodRepayTotalAmount = periodRepayInterest.Add(periodRepayPrinciple)
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
func calculateFixedPrincipalMethod(loanAmount decimal.Decimal, totalPeriodNum int) decimal.Decimal {
	return loanAmount.Div(decimal.NewFromInt(int64(totalPeriodNum))).RoundBank(2)
}
