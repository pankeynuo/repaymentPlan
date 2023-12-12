package main

import (
	"errors"
	"github.com/shopspring/decimal"
	"math"
	"strconv"
	"time"
)

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/5 10:09
**/
func fixedInstallmentMethod(request *Request) (*Response, error) {
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
	err = fixedInstallmentMethodPlan(response, repayPlanRequest{
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

func fixedInstallmentMethodPlan(response *Response, request repayPlanRequest, daysInterestRate decimal.Decimal) error {

	var sumTotalInterest, hasRepayPrincipal, sumTotalRepayAmount decimal.Decimal
	records := make([]RepayPlanRecord, 0)

	// 每期总还款金额(本金+利息)
	everyPeriodRepayAmount, err := calculateFixedInstallmentMethod(request.LoanAmount, request.PeriodInterestRate, request.TotalPeriodNum)
	if err != nil {
		return err
	}

	dateMap := calculatePeriodDate(request)

	for i := 0; i < request.TotalPeriodNum; i++ {
		periodStartDate := dateMap[i][0]
		periodEndDate := dateMap[i][1]
		periodRepayDate := dateMap[i][2]

		// 当前期次的计息天数
		daysOfPeriod := getDaysBetweenDate(periodStartDate, periodEndDate)

		// 当前期次的利息=当前剩余本金*计息天数*日利息
		periodRepayInterest := (request.LoanAmount.Sub(hasRepayPrincipal)).Mul(daysInterestRate).Mul(decimal.NewFromInt(daysOfPeriod)).RoundBank(2)

		record := RepayPlanRecord{
			PeriodNum:           i + 1,                                    // 当前期次的期数
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
			periodRepayPrinciple := everyPeriodRepayAmount.Sub(periodRepayInterest) // 当前期次还款本金
			hasRepayPrincipal = hasRepayPrincipal.Add(periodRepayPrinciple)         // 累积已还本金
			record.PeriodRepayPrinciple = periodRepayPrinciple.Round(2)             // 当前期次还款本金
			record.PeriodRepayTotalAmount = everyPeriodRepayAmount                  // 当前期次的总还款金额(除了最后一期，其他期次一样的金额)
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

// calculate the repayable amount of Fixed Installment Method
func calculateFixedInstallmentMethod(loanAmount, periodInterestRate decimal.Decimal, totalPeriodNum int) (decimal.Decimal, error) {
	planRepayAmount := decimal.Zero
	//如果利率为0
	if periodInterestRate.Equal(decimal.Zero) {
		planRepayAmount = loanAmount.Div(decimal.NewFromInt(int64(totalPeriodNum))).Round(2)
		return planRepayAmount, nil
	} else {
		// calculate every month need to repay total amount
		// (1+期利率）
		periodRateCal, err := strconv.ParseFloat(periodInterestRate.Add(decimal.NewFromInt(1)).String(), 64)
		if err != nil {
			return decimal.Zero, errors.New("int to float error:" + err.Error())
		}
		// (1+期利率)^期数
		pow := math.Pow(periodRateCal, float64(totalPeriodNum))
		// 每月还款金额=贷款本金*期利率*(1+期利率)^期数/((1+期利率)^期数-1) 保留两位小数
		planRepayAmount = loanAmount.Mul(periodInterestRate).Mul(decimal.NewFromFloat(pow)).
			Div(decimal.NewFromFloat(pow).Sub(decimal.NewFromFloat(1))).Round(2)
	}
	return planRepayAmount, nil
}
