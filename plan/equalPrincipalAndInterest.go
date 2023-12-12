package main

import (
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/5 10:19
**/
func equalPrincipalAndInterest(request *Request) (*Response, error) {
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
	err = equalPrincipalAndInterestPlan(response, repayPlanRequest{
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
func equalPrincipalAndInterestPlan(response *Response, request repayPlanRequest, daysInterestRate decimal.Decimal) error {
	var sumTotalInterest, hasRepayPrincipal, sumTotalRepayAmount decimal.Decimal

	records := make([]RepayPlanRecord, 0)

	planRepayPrinciplePeriod := request.LoanAmount.Div(decimal.NewFromInt(int64(request.TotalPeriodNum))).Round(2)

	totalCalcDays := getDaysBetweenDate(request.LoanStartDateParseLocal, request.LoanEndDateParseLocal)

	// 2.3 everyMonth need to repay interest amount = LoanAmount*daysRate*totalDays/periodNum
	planRepayInterestPeriod := request.LoanAmount.Mul(daysInterestRate).
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
