package main

import (
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

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
		nextRepayDate = calculateDateAddMonth(loanStartDateParseLocal, 0, repayDay)
	} else {
		// 如果过了还款日:下个月的还款日
		nextRepayDate = calculateDateAddMonth(loanStartDateParseLocal, 1, repayDay)
	}

	if (nextRepayDate.Sub(loanStartDateParseLocal)).Hours() < 20*24 {
		// 小于20天 取下一个月
		nextRepayDate = calculateDateAddMonth(nextRepayDate, 1, repayDay)

	}
	return nextRepayDate
}

func getTotalPeriodNum(request *Request, loanStartDateParseLocal time.Time) (totalPeriodNum int, err error) {
	if request.PeriodNum == 0 {
		totalPeriodNum, err = calculateTotalPeriodNum(request.LoanCycleCode, loanStartDateParseLocal, request.LoanEndDate, request.RepayDay)
		if nil != err {
			return 0, err
		}
	} else {
		if request.PeriodType == periodTypeYear {
			totalPeriodNum = 12 * request.PeriodNum
		} else {
			totalPeriodNum = request.PeriodNum
		}
	}
	return totalPeriodNum, nil
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
			repayDate := calculateDateAddMonth(loanStartDateParseLocal, period, repayDay)
			if repayDate.After(loanEndDateParseLocal) || repayDate.Equal(loanEndDateParseLocal) {
				break
			}
			period = period + 1
		}
	}
	return period, nil
}

func getLoanEndDate(request *Request, firstRepayDate time.Time) error {
	if request.LoanEndDate == "" || (request.PeriodNum != 0 && request.LoanEndDate != "") {
		var loanEndDateParseLocal time.Time
		switch request.LoanCycleCode {
		case loanCycleFortnightly:
			loanEndDateParseLocal = calculateLoanEndDateWithLoanCycleFortnightly(firstRepayDate, request.PeriodNum-1)
		case loanCycleMonthly:
			loanEndDateParseLocal = calculateDateAddMonth(firstRepayDate, request.PeriodNum-1, request.RepayDay)
		}
		request.LoanEndDate = loanEndDateParseLocal.Format(DATE_DASH_FORMAT)
	}
	return nil
}
func calculateLoanEndDateWithLoanCycleFortnightly(firstRepayDate time.Time, periodNum int) time.Time {
	cycle := 14 * (periodNum)
	loanEndDate := firstRepayDate.AddDate(0, 0, cycle)
	return loanEndDate
}

// 日期加n个月后的指定天
func calculateDateAddMonth(date time.Time, monthsNum, day int) time.Time {
	// the first day of fist loan period Date's month
	firstDay := date.AddDate(0, 0, -date.Day()+1)

	firstDayAddMonth := firstDay.AddDate(0, monthsNum, 0)
	firstDayAddMonthAddDay := firstDayAddMonth.AddDate(0, 0, day-1)
	if firstDayAddMonthAddDay.Day() != day {
		firstDayAndMonthAddDay := firstDayAddMonth.AddDate(0, 1, -firstDayAddMonth.Day())
		return firstDayAndMonthAddDay
	}
	return firstDayAddMonthAddDay
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

func getDaysBetweenDate(startDate, endDate time.Time) int64 {
	startTime := startDate.Unix()
	endTime := endDate.Unix()
	// 求相差天数
	date := (endTime-startTime)/86400 + 1
	return date
}

// calculate every period startDate,endDate,repayDate
func calculatePeriodDate(request repayPlanRequest) map[int][]time.Time {
	dateMap := make(map[int][]time.Time)
	for i := 0; i < request.TotalPeriodNum; i++ {
		periodStartDate := time.Time{} // 计息开始日
		periodRepayDate := time.Time{} // 还款日

		if i == 0 { // 第一期
			periodStartDate = request.LoanStartDateParseLocal
			periodRepayDate = request.FirstRepayDate
		} else {
			if request.LoanCycleCode == loanCycleFortnightly {
				periodStartDate = dateMap[i-1][2]
				repaymentTs := request.FirstRepayDate.AddDate(0, 0, i*14)
				periodRepayDate = time.Unix(repaymentTs.Unix(), 0)
			}
			if request.LoanCycleCode == loanCycleMonthly {
				periodRepayDateTemp := calculateDateAddMonth(request.FirstRepayDate, i, request.RepayDay)
				periodRepayDate = periodRepayDateTemp
				periodStartDate = dateMap[i-1][2]
			}

			if i == request.TotalPeriodNum-1 { // 最后一期
				periodRepayDate = request.LoanEndDateParseLocal
			}
		}
		// 计息结束日=还款日的前一天
		periodEndDate := periodRepayDate.AddDate(0, 0, -1)
		dateMap[i] = []time.Time{periodStartDate, periodEndDate, periodRepayDate}
	}
	return dateMap
}
