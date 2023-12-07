package main

import (
	"errors"
	"github.com/shopspring/decimal"
	"log"
	"strconv"
	"time"
)

func CalculateRepaymentPlan(request *RepayPlanRequest) (response *RepayPlanResponse, err error) {
	if e := check(request); nil != e {
		return nil, e
	}
	return getRepaymentPlan(request)
}

// request 参数检查（不涉及业务）
func check(request *RepayPlanRequest) error {
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
func getRepaymentPlan(request *RepayPlanRequest) (response *RepayPlanResponse, err error) {
	response = &RepayPlanResponse{}
	switch request.RepayMethod {
	// 01-等额本息
	case EqualLoanRepayment:
		err = equalLoanRepayment(request, response)
	// 02-等额本金
	case EqualPrincipalRepayment:
		err = equalPrincipalRepayment(request, response)
	// 03-到期（一次性）还本付息（息随本清)(到期一次性还本还息）
	case BothPrincipalAndInterest:
		err = bothPrincipalAndInterest(request, response)
	// 04-到期还本周期还息(分期付息到期还本（先息后本))
	case BeforeInterestAfterPrincipal:
		err = beforeInterestAfterPrincipal(request, response)
	// 05-等本等息（每期还本还息还款额都相等，每期计息的本金为贷款总本金）
	case EqualPrincipalAndInterest:
		err = equalPrincipalAndInterest(request, response)
	// 06 账单。。。。
	case BillDateRepayment:
		err = billDateRepayment(request, response)
	default:
		return nil, errors.New("repay method error")
	}
	return response, err
}

// 获取第一个还款日
func getFirstRepayDate(loanStartDate, loanEndDate, loanCycleCode string, repayDay int) (string, error) {
	nextRepayDate, err := calculateFirstRepayDate(loanStartDate, loanCycleCode, repayDay)
	if nil != err {
		return nextRepayDate, err
	}
	// 如果下一还款日比到期日还大，则下一还款日就是到期日
	if loanEndDate != "" && nextRepayDate > loanEndDate {
		return loanEndDate, err
	}
	return nextRepayDate, err
}

// 计算第一个还款日
func calculateFirstRepayDate(loanStartDate, loanCycleCode string, repayDay int) (nextRepayDate string, err error) {
	loanStartDateParseLocal, err := time.ParseInLocation(DATE_DASH_FORMAT, loanStartDate, time.Local)
	if err != nil {
		return "", errors.New("loanStartDate date format error: " + err.Error())
	}

	switch loanCycleCode {
	case loanCycleFortnightly:
		return getFirstRepayDateOfLoanCycleFortnightly(repayDay, loanStartDateParseLocal), nil
	case loanCycleMonthly:
		return getFirstRepayDateOfLoanCycleMonthly(repayDay, loanStartDateParseLocal), nil
	}

	return "", nil
}
func getFirstRepayDateOfLoanCycleFortnightly(repayDay int, loanStartDateParseLocal time.Time) string {
	dateTime := time.Time{}
	// startDate add 14 days and get his weekDay num
	// 获取 起息日+14天之后是周几
	startDateUnixTime := time.Unix(loanStartDateParseLocal.Unix(), 0)
	weekDay := weekDayToDay(startDateUnixTime.AddDate(0, 0, 14).Weekday())

	// If the week date  of the start date is different from the repayment date,
	// the corresponding day two weeks after the week of the start date is taken as the first repayment date
	// 首个还款日(起息日+14天)和 还款日不一样,则以起息日的下下周为首个还款日
	if repayDay != weekDay {
		dateTime = loanStartDateParseLocal.AddDate(0, 0, 14+repayDay-weekDay)
	} else {
		dateTime = loanStartDateParseLocal.AddDate(0, 0, 7)
	}
	nextRepayDate := dateTime.Format(DATE_DASH_FORMAT)
	return nextRepayDate
}
func getFirstRepayDateOfLoanCycleMonthly(repayDay int, loanStartDateParseLocal time.Time) string {
	dateTime := time.Time{}
	if repayDay > loanStartDateParseLocal.Day() {
		// 如果还没过还款日:这个月的还款日
		dateTime = loanStartDateParseLocal.AddDate(0, 0, repayDay-loanStartDateParseLocal.Day())
	} else {
		// 如果过了还款日:下个月的还款日
		dateTime = loanStartDateParseLocal.AddDate(0, 1, repayDay-loanStartDateParseLocal.Day())
	}
	log.Printf("temp:%s", dateTime)
	if dateTime.Day() != repayDay { // 日期加完和还款日不一致，月末问题,改为获取下一个月的月末
		dateTime = getMonthLastDate(dateTime)
	}

	if (dateTime.Sub(loanStartDateParseLocal)).Hours() < 20*24 {
		// 小于20天 再取下一个月
		dateTime = dateTime.AddDate(0, 1, repayDay-loanStartDateParseLocal.Day())

		if dateTime.Day() != repayDay { // 日期加完和还款日不一致，月末问题,改为获取下一个月的月末
			dateTime = getMonthLastDate(dateTime)
		}
	}
	nextRepayDate := dateTime.Format(DATE_DASH_FORMAT)
	return nextRepayDate
}
func getMonthLastDate(dateTime time.Time) time.Time {
	firstOfMonth := time.Date(dateTime.Year(), dateTime.Month(), 1, 0, 0, 0, 0, dateTime.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 0, -1)
	return lastOfMonth
}

func getLoanEndDateAndTotalPeriodNum(periodNum int, periodType, loanEndDate, loanStartDate, loanCycleCode string) (totalPeriodNum int, _ string, err error) {
	if periodNum == 0 && loanEndDate == "" {
		return 0, "", errors.New("loanEndDate and periodNum can not be empty at the same time")
	}
	if periodNum == 0 {
		// 总期数为空，最后还款日不为空：计算总期数
		// TODO totalPeriodNum
		totalPeriodNum = calculateTotalPeriodNum()
	} else {
		//总期次不为空，最后一个还款日为空：计算最后一个还款日
		if periodType == "01" { // 01-年
			totalPeriodNum = 12 * periodNum
		} else {
			totalPeriodNum = periodNum
		}
		loanEndDate = calculateLoanEndDate(loanStartDate, periodNum, loanCycleCode)
	}
	return totalPeriodNum, loanEndDate, nil
}
func calculateLoanEndDate(loanStartDateStr string, period int, loanCycleCode string) string {
	loanStartDate, _ := time.ParseInLocation(DATE_DASH_FORMAT, loanStartDateStr, time.Local)
	var matureDate time.Time
	switch loanCycleCode {
	case loanCycleFortnightly:
		cycle := 14 * (period)
		matureDate = loanStartDate.AddDate(0, 0, cycle)
	case loanCycleMonthly:
		matureDate = loanStartDate.AddDate(0, period, 0)
	}
	return matureDate.Format(DATE_DASH_FORMAT)
}
func calculateTotalPeriodNum() int {
	return 0
}
func calculatePeriodNum(nextRepayDate time.Time, loanCycleCode, InterestCalculateEndDate, RepayDay string) (int, error) {
	// TODO
	return 0, nil
}
func calculatePeriodNumAndNewNextRepayDate(nextRepayDate time.Time, loanCycleCode, interestCalculateStartDate, InterestCalculateEndDate, RepayDay string) (time.Time, int, error) {
	// TODO
	return time.Time{}, 0, nil
}
func calculateMatureDate(nextRepayDate time.Time, loanCycleCode, RepayDay string, periodNum int) (string, error) {
	// TODO
	return "", nil
}
func calculateDailyInterestRate(interestRate decimal.Decimal, daysOfYear int) decimal.Decimal {
	return interestRate.Div(decimal.NewFromInt(int64(daysOfYear))).Div(decimal.NewFromFloat(100))
}

func stringToInt(num string) (int, error) {
	intNum, err := strconv.Atoi(num)
	if err != nil {
		return 0, err
	}
	return intNum, nil
}
func weekDayToDay(weekday time.Weekday) int {
	switch weekday {
	case time.Monday:
		return 1
	case time.Tuesday:
		return 2
	case time.Wednesday:
		return 3
	case time.Thursday:
		return 4
	case time.Friday:
		return 5
	case time.Saturday:
		return 6
	case time.Sunday:
		return 7
	default:
		return 0
	}
}

//获取日期月份第一天、最后一天
func getFirstLastDateOfMonth(t string) (string, string) {
	d, err := time.ParseInLocation(DATE_DASH_FORMAT, t, time.Local)
	if err != nil {
		d = time.Now()
	}
	firstDate := d.AddDate(0, 0, -d.Day()+1)
	lastDate := d.AddDate(0, 1, -d.Day())
	return firstDate.Format(DATE_DASH_FORMAT), lastDate.Format(DATE_DASH_FORMAT)
}
