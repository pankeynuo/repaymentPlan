package main

import (
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

func CalculateRepaymentPlan(request *Request) (response *Response, err error) {
	if err = check(request); nil != err {
		return nil, err
	}
	return getRepaymentPlan(request)
}

// request 参数检查
func check(request *Request) error {

	if request.DaysOfYear == 0 {
		request.DaysOfYear = daysOfYear
	}
	if request.InterestRate.LessThanOrEqual(decimal.Zero) {
		return errors.New("interest Rate error")
	}
	if request.LoanAmount.LessThanOrEqual(decimal.Zero) {
		return errors.New("loan Amount error")
	}
	if e := checkLoanCycleCode(request.LoanCycleCode); nil != e {
		return e
	}
	if request.PeriodNum < 0 {
		return errors.New("period Num error")
	}
	if request.LoanStartDate == "" {
		return errors.New("interest Calculate Start Date can not be empty")
	}
	loanStartDate, e := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanStartDate, time.Local)
	if nil != e {
		return errors.New("interest Calculate Start Date error")
	}
	if request.LoanEndDate != "" {
		loanEndDate, e := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanEndDate, time.Local)
		if nil != e {
			return errors.New("interest Calculate End Date error")
		}
		if loanEndDate.Sub(loanStartDate) <= 0 {
			return errors.New("loan Start Date can not after or equal than loan end date")
		}
	}

	if request.RepayDay <= 0 || request.RepayDay >= 32 {
		return errors.New("repay Day error")
	}
	if e := checkPeriodType(request.PeriodType); nil != e {
		return e
	}
	switch request.RepayMethod {
	case EqualLoanRepayment, EqualPrincipalRepayment, BeforeInterestAfterPrincipal, EqualPrincipalAndInterest:
		if request.PeriodNum == 0 && request.LoanEndDate == "" {
			return errors.New("loanEndDate and periodNum can not be empty at the same time")
		}
	case BothPrincipalAndInterest:
		if request.LoanEndDate == "" {
			return errors.New("interest Calculate End Date can not be empty")
		}
	}
	return nil
}
func checkLoanCycleCode(loanCycleCode string) error {
	switch loanCycleCode {
	case loanCycleDaily, loanCycleFortnightly, loanCycleMonthly, loanCycleQuarterly, loanCycleYearly:
		return nil
	default:
		return errors.New("loan Cycle Code error")
	}
}
func checkPeriodType(periodType string) error {
	switch periodType {
	case periodTypeYear, periodTypeMonth:
		return nil
	default:
		return errors.New("period type error")
	}
}

func getRepaymentPlan(request *Request) (response *Response, err error) {

	var repayPlanRequest repayPlanRequest
	switch request.RepayMethod {
	case EqualLoanRepayment, EqualPrincipalRepayment, BeforeInterestAfterPrincipal, EqualPrincipalAndInterest:
		repayPlanRequest, response, err = prepareGetParameter(request)
	}

	switch request.RepayMethod {
	// 01-等额本息
	case EqualLoanRepayment:
		err = fixedInstallmentMethodPlan(repayPlanRequest, response)
	// 02-等额本金
	case EqualPrincipalRepayment:
		err = fixedPrincipalMethodPlan(repayPlanRequest, response)
	// 03-利随本清
	case BothPrincipalAndInterest:
		response, err = bothPrincipalAndInterest(request)
	// 04-先息后本
	case BeforeInterestAfterPrincipal:
		err = beforeInterestAfterPrincipalPlan(repayPlanRequest, response)
	// 05-等本等息
	case EqualPrincipalAndInterest:
		err = equalPrincipalAndInterestPlan(repayPlanRequest, response)
	default:
		return nil, errors.New("repay method error")
	}
	return response, err
}
func prepareGetParameter(request *Request) (repayPlanRequest, *Response, error) {
	loanStartDateParseLocal, err := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanStartDate, time.Local)
	if err != nil {
		return repayPlanRequest{}, nil, errors.New("loanStartDate date format error: " + err.Error())
	}

	firstRepayDate, err := getFirstRepayDate(request, loanStartDateParseLocal)
	if err != nil {
		return repayPlanRequest{}, nil, err
	}
	totalPeriodNum, err := getTotalPeriodNum(request, loanStartDateParseLocal)

	err = getLoanEndDate(request, firstRepayDate)

	loanEndDateParseLocal, err := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanEndDate, time.Local)
	if err != nil {
		return repayPlanRequest{}, nil, errors.New("loanStartDate date format error: " + err.Error())
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

	return repayPlanRequest{
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
		DaysInterestRate:        daysInterestRate,
	}, response, nil
}
