package main

import (
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/5 10:18
**/
func fixedPrincipalMethod(request *Request) (*Response, error) {
	loanStartDateParseLocal, err := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanStartDate, time.Local)
	if err != nil {
		return nil, errors.New("loanStartDate date format error: " + err.Error())
	}

	firstRepayDate, err := getFirstRepayDate(request, loanStartDateParseLocal)
	if err != nil {
		return nil, err
	}
	totalPeriodNum, err := getTotalPeriodNum(request, loanStartDateParseLocal)

	err = getLoanEndDate(request, firstRepayDate)

	loanEndDateParseLocal, err := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanEndDate, time.Local)
	if err != nil {
		return nil, errors.New("loanStartDate date format error: " + err.Error())
	}

	periodInterestRate := calculatePeriodInterestRate(request.InterestRate, request.LoanCycleCode)

	daysInterestRate := calculateDaysInterestRate(request.InterestRate, request.DaysOfYear)

	response := &Response{
		RepayMethod:    request.RepayMethod,
		LoanStartDate:  request.LoanStartDate,
		LoanEndDate:    request.LoanEndDate,
		TotalPeriodNum: totalPeriodNum,
		LoanAmount:     request.LoanAmount,
		InterestRate:   request.InterestRate,
	}
	err = fixedPrincipalMethodPlan(response, repayPlanRequest{
		LoanAmount:              request.LoanAmount,
		LoanStartDate:           request.LoanStartDate,
		LoanEndDate:             request.LoanEndDate,
		LoanCycleCode:           request.LoanCycleCode,
		PeriodInterestRate:      periodInterestRate,
		TotalPeriodNum:          totalPeriodNum,
		FirstRepayDate:          firstRepayDate,
		LoanStartDateParseLocal: loanStartDateParseLocal,
		LoanEndDateParseLocal:   loanEndDateParseLocal,
		RepayDay:                request.RepayDay,
	}, daysInterestRate)
	return response, nil
}

func fixedPrincipalMethodPlan(response *Response, request repayPlanRequest, daysInterestRate decimal.Decimal) error {
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

		periodRepayInterest := (request.LoanAmount.Sub(hasRepayPrincipal)).Mul(daysInterestRate).Mul(decimal.NewFromInt(daysOfPeriod)).RoundBank(2)

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
