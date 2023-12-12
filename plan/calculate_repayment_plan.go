package main

import (
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

func CalculateRepaymentPlan(request *Request) (response *Response, err error) {
	return getRepaymentPlan(request)
}

// request 参数检查
func check(request *Request) error {

	if request.DaysOfYear == 0 {
		request.DaysOfYear = daysOfYear
	}
	if request.DaysOfMonth == 0 {
		request.DaysOfMonth = daysOfMonth
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
	if request.PeriodNum == 0 && request.LoanEndDate == "" {
		return errors.New("loanEndDate and periodNum can not be empty at the same time")
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
	switch request.RepayMethod {
	// 01-等额本息
	case EqualLoanRepayment:
		if err = check(request); nil != err {
			return nil, err
		}
		response, err = fixedInstallmentMethod(request)
	// 02-等额本金
	case EqualPrincipalRepayment:
		if err = check(request); nil != err {
			return nil, err
		}
		response, err = fixedPrincipalMethod(request)
	// 03-利随本清
	case BothPrincipalAndInterest:
		loanStartDate, loanEndDate, err := check2(request)
		if nil != err {
			return nil, err
		}
		response, err = bothPrincipalAndInterest(request, loanStartDate, loanEndDate)
	// 04-先息后本
	case BeforeInterestAfterPrincipal:
		if err = check(request); nil != err {
			return nil, err
		}
		response, err = beforeInterestAfterPrincipal(request)
	// 05-等本等息
	case EqualPrincipalAndInterest:
		if err = check(request); nil != err {
			return nil, err
		}
		response, err = equalPrincipalAndInterest(request)
	// 06 账单。。。。
	case BillDateRepayment:
		if err = check3(request);nil!=err {
			return nil,err
		}
		response, err = billDateRepayment(request)
	default:
		return nil, errors.New("repay method error")
	}
	return response, err
}
func check2(request *Request) (time.Time, time.Time, error) {
	if request.LoanStartDate == "" {
		return time.Time{}, time.Time{}, errors.New("interest Calculate Start Date can not be empty")
	}
	if request.LoanEndDate == "" {
		return time.Time{}, time.Time{}, errors.New("interest Calculate End Date can not be empty")
	}
	loanStartDate, e := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanStartDate, time.Local)
	if nil != e {
		return time.Time{}, time.Time{}, errors.New("interest Calculate Start Date error")
	}
	loanEndDate, e := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanEndDate, time.Local)
	if nil != e {
		return time.Time{}, time.Time{}, errors.New("interest Calculate End Date error")
	}
	if loanEndDate.Sub(loanStartDate) <= 0 {
		return time.Time{}, time.Time{}, errors.New("loan Start Date can not after or equal than loan end date")
	}
	if request.DaysOfYear == 0 {
		request.DaysOfYear = daysOfYear
	}
	if request.InterestRate.LessThanOrEqual(decimal.Zero) {
		return time.Time{}, time.Time{}, errors.New("interest Rate error")
	}
	return loanStartDate, loanEndDate, nil
}
func check3(request *Request) error {
	if request.LoanStartDate == "" {
		return errors.New("interest Calculate Start Date can not be empty")
	}
	if request.LoanEndDate == "" {
		return errors.New("interest Calculate End Date can not be empty")
	}
	loanStartDate, e := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanStartDate, time.Local)
	if nil != e {
		return errors.New("interest Calculate Start Date error")
	}
	loanEndDate, e := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanEndDate, time.Local)
	if nil != e {
		return errors.New("interest Calculate End Date error")
	}
	if loanEndDate.Sub(loanStartDate) <= 0 {
		return errors.New("loan Start Date can not after or equal than loan end date")
	}
	if request.LoanCycleCode != loanCycleMonthly {

	}
	if request.RepayDay <= 0 || request.RepayDay >= 32 {
		return errors.New("repay Day error")
	}
	if request.BillDay <= 0 || request.BillDay >= 32 {
		return errors.New("repay Day error")
	}
	return nil
}
