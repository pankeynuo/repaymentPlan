package main

import (
	"errors"
	"github.com/shopspring/decimal"
	"strconv"
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
		err = billDateRepayment(request)
	default:
		return nil, errors.New("repay method error")
	}
	return response, err
}

// 获取第一个还款日
func getFirstRepayDate(request *Request, loanStartDateParseLocal time.Time) (time.Time, error) {
	nextRepayDate := calculateFirstRepayDate(loanStartDateParseLocal, request.LoanCycleCode, request.RepayDay)

	// 如果下一还款日比到期日还大，则下一还款日就是到期日
	if request.LoanEndDate != "" {
		loanEndDateParseLocal, e := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanEndDate, time.Local)
		if nil != e {
			return time.Time{}, errors.New("interest Calculate End Date error")
		}
		if nextRepayDate.After(loanEndDateParseLocal) {
			return loanEndDateParseLocal, nil
		}
	}
	return nextRepayDate, nil
}

// 计算第一个还款日
func calculateFirstRepayDate(loanStartDateParseLocal time.Time, loanCycleCode string, repayDay int) time.Time {
	switch loanCycleCode {
	case loanCycleFortnightly:
		return getFirstRepayDateOfLoanCycleFortnightly(repayDay, loanStartDateParseLocal)
	case loanCycleMonthly:
		return getFirstRepayDateOfLoanCycleMonthly(repayDay, loanStartDateParseLocal)
	}

	return time.Time{}
}
func getFirstRepayDateOfLoanCycleFortnightly(repayDay int, loanStartDateParseLocal time.Time) time.Time {
	nextRepayDate := time.Time{}
	// startDate add 14 days and get his weekDay num
	// 获取 起息日+14天之后是周几
	startDateUnixTime := time.Unix(loanStartDateParseLocal.Unix(), 0)
	weekDay := weekDayToDay(startDateUnixTime.AddDate(0, 0, 14).Weekday())

	// If the week date  of the start date is different from the repayment date,
	// the corresponding day two weeks after the week of the start date is taken as the first repayment date
	// 首个还款日(起息日+14天)和 还款日不一样,则以起息日的下下周为首个还款日
	if repayDay != weekDay {
		nextRepayDate = loanStartDateParseLocal.AddDate(0, 0, 14+repayDay-weekDay)
	} else {
		nextRepayDate = loanStartDateParseLocal.AddDate(0, 0, 7)
	}
	return nextRepayDate
}
func getFirstRepayDateOfLoanCycleMonthly(repayDay int, loanStartDateParseLocal time.Time) time.Time {
	nextRepayDate := time.Time{}
	if repayDay > loanStartDateParseLocal.Day() {
		// 如果还没过还款日:这个月的还款日
		nextRepayDate = loanStartDateParseLocal.AddDate(0, 0, repayDay-loanStartDateParseLocal.Day())
	} else {
		// 如果过了还款日:下个月的还款日
		nextRepayDate = loanStartDateParseLocal.AddDate(0, 1, repayDay-loanStartDateParseLocal.Day())
	}
	if nextRepayDate.Day() != repayDay { // 日期加完和还款日不一致，月末问题,改为获取下一个月的月末
		nextRepayDate = getMonthLastDate(nextRepayDate)
	}

	if (nextRepayDate.Sub(loanStartDateParseLocal)).Hours() < 20*24 {
		// 小于20天 再取下一个月
		nextRepayDate = nextRepayDate.AddDate(0, 1, repayDay-loanStartDateParseLocal.Day())

		if nextRepayDate.Day() != repayDay { // 日期加完和还款日不一致，月末问题,改为获取下一个月的月末
			nextRepayDate = getMonthLastDate(nextRepayDate)
		}
	}
	return nextRepayDate
}
func getMonthLastDate(dateTime time.Time) time.Time {
	firstOfMonth := time.Date(dateTime.Year(), dateTime.Month(), 1, 0, 0, 0, 0, dateTime.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 0, -1)
	return lastOfMonth
}

func getTotalPeriodNum(request *Request, loanStartDateParseLocal time.Time) (totalPeriodNum int, err error) {
	if request.PeriodNum == 0 {
		totalPeriodNum, err = calculateTotalPeriodNum(request.LoanCycleCode, loanStartDateParseLocal, request.LoanEndDate, request.RepayDay)
		if nil != err {
			return 0, err
		}
	} else {
		if request.PeriodType == periodTypeYear { // 01-年
			totalPeriodNum = 12 * request.PeriodNum
		} else {
			totalPeriodNum = request.PeriodNum
		}
	}
	return totalPeriodNum, nil
}
func getLoanEndDate(request *Request, loanStartDateParseLocal time.Time) error {
	if request.LoanEndDate == "" {
		var loanEndDateParseLocal time.Time
		switch request.LoanCycleCode {
		case loanCycleFortnightly:
			loanEndDateParseLocal = calculateLoanEndDateWithLoanCycleFortnightly(loanStartDateParseLocal, request.PeriodNum)
		case loanCycleMonthly:
			loanEndDateParseLocal = calculateDateAndMonth(loanStartDateParseLocal, request.PeriodNum, request.RepayDay)
		}
		request.LoanEndDate = loanEndDateParseLocal.Format(DATE_DASH_FORMAT)
	}
	return nil
}
func calculateLoanEndDateWithLoanCycleFortnightly(loanStartDateParseLocal time.Time, periodNum int) time.Time {
	cycle := 14 * (periodNum)
	loanEndDate := loanStartDateParseLocal.AddDate(0, 0, cycle)
	return loanEndDate
}
func calculateDateAndMonth(loanStartDateParseLocal time.Time, periodNum, repayDay int) time.Time {
	// the first day of fist loan period Date's month
	firstDay := loanStartDateParseLocal.AddDate(0, 0, -loanStartDateParseLocal.Day()+1)

	lastPeriod := firstDay.AddDate(0, periodNum, 0)

	loanEndDate := lastPeriod.AddDate(0, 0, repayDay-1)

	if loanEndDate.Day() != repayDay {
		loanEndDate := lastPeriod.AddDate(0, 1, -lastPeriod.Day())
		return loanEndDate
	}
	return loanEndDate
}

func calculateTotalPeriodNum(loanCycleCode string, loanStartDateParseLocal time.Time, loanEndDate string, repayDay int) (int, error) {
	period := 1
	loanEndDateParseLocal, err := time.ParseInLocation(DATE_DASH_FORMAT, loanEndDate, time.Local)
	if err != nil {
		return 0, errors.New("loanEndDate date format error: " + err.Error())
	}
	switch loanCycleCode {
	case loanCycleFortnightly:
		cycle := 14
		for {
			// accumulate period 14 days 累加period个14天
			repayDateAddCycle := loanStartDateParseLocal.AddDate(0, 0, period*cycle)
			// when after repayDate add cycle less or equal maturityDate,finish
			// 累加后的日期小于maturityDate,结束循环，返回 period
			if repayDateAddCycle.After(loanEndDateParseLocal) || repayDateAddCycle.Equal(loanEndDateParseLocal) {
				break
			}
			period = period + 1
		}
	case loanCycleMonthly:
		for {
			repayDate := calculateDateAndMonth(loanStartDateParseLocal, period, repayDay)
			if repayDate.After(loanEndDateParseLocal) || repayDate.Equal(loanEndDateParseLocal) {
				break
			}
			period = period + 1
		}
	}
	return period, nil
}

func calculatePeriodInterestRate(interestRate decimal.Decimal, loanCycleCode string) decimal.Decimal {
	switch loanCycleCode {
	case loanCycleFortnightly:
		return interestRate.Div(decimal.NewFromFloat(float64(numberOfWeek))).Div(decimal.NewFromFloat(100))
	case loanCycleMonthly:
		return interestRate.Div(decimal.NewFromInt(int64(numberOfMonth))).Div(decimal.NewFromFloat(100))
	}
	return decimal.Decimal{}
}
func calculateDaysInterestRate(interestRate decimal.Decimal, daysOfYear int) decimal.Decimal {
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

// 获取日期月份第一天、最后一天
func getFirstLastDateOfMonth(t string) (string, string) {
	d, err := time.ParseInLocation(DATE_DASH_FORMAT, t, time.Local)
	if err != nil {
		d = time.Now()
	}
	firstDate := d.AddDate(0, 0, -d.Day()+1)
	lastDate := d.AddDate(0, 1, -d.Day())
	return firstDate.Format(DATE_DASH_FORMAT), lastDate.Format(DATE_DASH_FORMAT)
}
func getDaysBetweenDate(startDate, endDate time.Time) int64 {
	startTime := startDate.Unix()
	endTime := endDate.Unix()
	// 求相差天数
	date := (endTime-startTime)/86400 + 1
	return date
}

// calculate every period startDate,endDate,repayDate
func calculatePeriodDate(request repayPlanRequest) map[int][]string {
	dateMap := make(map[int][]string)
	for i := 0; i < request.TotalPeriodNum; i++ {
		periodStartDate := "" // 计息开始日
		periodRepayDate := "" // 还款日

		if i == 0 { // 第一期
			periodStartDate = request.LoanStartDate
			periodRepayDate = request.FirstRepayDate.Format(DATE_DASH_FORMAT)
		} else {
			if request.LoanCycleCode == loanCycleFortnightly {

				// recalBgnTs := request.FirstRepayDate.AddDate(0, 0, (i-1)*14)
				// periodStartDate = time.Unix(recalBgnTs.Unix(), 0).Format(DATE_DASH_FORMAT)
				periodStartDate = dateMap[i-1][2]
				repaymentTs := request.FirstRepayDate.AddDate(0, 0, i*14)
				periodRepayDate = time.Unix(repaymentTs.Unix(), 0).Format(DATE_DASH_FORMAT)
			}
			if request.LoanCycleCode == loanCycleMonthly {
				periodRepayDateTemp := calculateDateAndMonth(request.FirstRepayDate, i, request.RepayDay)
				periodRepayDate = periodRepayDateTemp.Format(DATE_DASH_FORMAT)
				periodStartDate = dateMap[i-1][2]
			}

			if i == request.TotalPeriodNum-1 { // 最后一期
				periodRepayDate = request.LoanEndDate
			}
		}
		// 计息结束日=还款日的前一天
		periodEndDateTs, _ := time.Parse(DATE_DASH_FORMAT, periodRepayDate)
		periodEndDate := periodEndDateTs.AddDate(0, 0, -1).Format(DATE_DASH_FORMAT)
		dateMap[i] = []string{periodStartDate, periodEndDate, periodRepayDate}
	}
	return dateMap
}
