package main

import (
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

func CalculateRepaymentPlan(request *Request) (response *Response, err error) {
	if e := check(request); nil != e {
		return nil, e
	}
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
		response, err = equalLoanRepayment(request)
	// 02-等额本金
	case EqualPrincipalRepayment:
		response, err = equalPrincipalRepayment(request)
	// 03-到期（一次性）还本付息（息随本清)(到期一次性还本还息）
	case BothPrincipalAndInterest:
		response, err = bothPrincipalAndInterest(request)
	// 04-到期还本周期还息(分期付息到期还本（先息后本))
	case BeforeInterestAfterPrincipal:
		response, err = beforeInterestAfterPrincipal(request)
	// 05-等本等息（每期还本还息还款额都相等，每期计息的本金为贷款总本金）
	case EqualPrincipalAndInterest:
		response, err = equalPrincipalAndInterest(request)
	// 06 账单。。。。
	case BillDateRepayment:
		response, err = billDateRepayment(request)
	default:
		return nil, errors.New("repay method error")
	}
	return response, err
}
