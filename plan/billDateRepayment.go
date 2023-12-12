package main

import (
	"errors"
	"git.multiverse.io/framework/log"
	"git.siriustech.internal/banking/backend/icredit/common/models"
	"time"
)

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/5 10:19
**/
func billDateRepayment(request *Request) (response *Response, err error) {

	loanStartDateParseLocal, err := time.ParseInLocation(DATE_DASH_FORMAT, request.LoanStartDate, time.Local)
	if err != nil {
		return nil, errors.New("loanStartDate date format error: " + err.Error())
	}

	nextBillDate, err := calcNextDateWithMonth(request)
	if err != nil {
		return nil, err
	}
	nextBillDateTs, e := time.Parse(DATE_DASH_FORMAT, nextBillDate)
	if e != nil {
		return nil, errors.New()
	}
	nextRepayDate, err := calcFirstRepayDateForBill(request)
	if err != nil {
		return nil, err
	}
	nextRepayDateTs, e := time.Parse(DATE_DASH_FORMAT, nextRepayDate)
	if e != nil {
		return nil, errors.New()
	}

	totalPeriodNum, planRecords := calcBillRepayPlan(request, nextBillDateTs, nextRepayDateTs, billDay, repayDay)

	response := &Response{
		RepayMethod:    request.RepayMethod,
		LoanStartDate:  request.LoanStartDate,
		LoanEndDate:    request.LoanEndDate,
		TotalPeriodNum: totalPeriodNum,
		LoanAmount:     request.LoanAmount,   // 好像没有
		InterestRate:   request.InterestRate, // 好像没有
		FirstBillDate:  nextBillDate,
		FirstRepayDate: nextRepayDate,
		Records:        planRecords,
	}
	return response, nil
}

// calc first repayment date when repayment method = 06
// 2）BD<RD，则BDn和RDn在同一个月内；
// 3）BD>RD，则RDn在BDn所在月的次月；
// 4）BD=RD，则页面提供选项，RDn是否与BDn在同一个月（选择否，则BDn在RDn所在月的次月）
func calcFirstRepayDateForBill(plan *Request) (nextRepayDate string, err errors) {
	billDay := plan.BillDay
	repayDay := plan.RepayDay
	if billDay < repayDay {
		nextRepayDate, err = calcNextDateWithMonth(plan, plan.RepayDay, 1)
	} else if billDay == repayDay {
		if plan.IsSameFlag {
			nextRepayDate, err = calcNextDateWithMonth(plan, plan.RepayDay, 1)
		} else {
			nextRepayDate, err = calcNextDateWithMonth(plan, plan.RepayDay, 2)
		}
	} else {
		nextRepayDate, err = calcNextDateWithMonth(plan, plan.RepayDay, 2)
	}
	if err != nil {
		return "", err
	}
	return nextRepayDate, nil
}

// calc the repayment date , plan bill date for billing repayment plan
func calcBillRepayPlan(plan *Request, nextBillDateTs, nextRepayDateTs time.Time,
	billDay, repayDay int) (totalPeriod int, records []models.RepayPlanRecord) {

	endDateTs, e := time.Parse(DATE_DASH_FORMAT, plan.InterestCalculateEndDate)
	if e != nil {
		panic(e)
	}

	//If the next billing date calculated for the first time is equal to the due date,
	//then it means that there is no bill for this data
	//return total period equal zero , and no bill records
	//如果首次计算出来的下一个账单日等于到期日，那么代表该笔数据是没有账单的
	if nextBillDateTs.After(endDateTs) || nextBillDateTs.Equal(endDateTs) {
		return 0, nil
	}

	//calc bill
	planRecords := []models.RepayPlanRecord{}
	period := 0
	for {
		record := models.RepayPlanRecord{}
		if period == 0 {
			record.InterestCalculateStartDate = plan.InterestCalculateStartDate
			period += 1
			log.Debugsf("period = %d", period)
			log.Debugsf("nextBillDateTs = %s", nextBillDateTs)
			record.PeriodNum = period
			record.PlanBillDate = nextBillDateTs.Format(constant.DATE_DASH_FORMAT)
			record.InterestCalculateEndDate = nextBillDateTs.Format(constant.DATE_DASH_FORMAT)
			record.PlanRepayDate = nextRepayDateTs.Format(constant.DATE_DASH_FORMAT)
			dayOfCalc, _ := GetTimeArr(record.InterestCalculateStartDate, record.InterestCalculateEndDate)
			record.DaysOfInterestCalculate = int(dayOfCalc)
			planRecords = append(planRecords, record)
			continue
		}
		if nextBillDateTs.After(endDateTs) || nextBillDateTs.Equal(endDateTs) {
			break
		}
		record.InterestCalculateStartDate = nextBillDateTs.Format(constant.DATE_DASH_FORMAT)

		nextBillDateTs = calcNextBillDateOrRepaymentDate(nextBillDateTs, billDay)
		if nextBillDateTs.After(endDateTs) || nextBillDateTs.Equal(endDateTs) {
			break
		}
		nextRepayDateTs = calcNextBillDateOrRepaymentDate(nextRepayDateTs, repayDay)

		period += 1
		log.Debugsf("period = %d", period)
		log.Debugsf("nextBillDateTs = %s", nextBillDateTs)
		record.PeriodNum = period

		interestCalculateEndDate := nextBillDateTs.Format(constant.DATE_DASH_FORMAT)
		if nextRepayDateTs.After(endDateTs) || nextRepayDateTs.Equal(endDateTs) {
			interestCalculateEndDate = endDateTs.Format(constant.DATE_DASH_FORMAT)
		}

		record.PlanBillDate = nextBillDateTs.Format(constant.DATE_DASH_FORMAT)
		record.InterestCalculateEndDate = interestCalculateEndDate
		record.PlanRepayDate = nextRepayDateTs.Format(constant.DATE_DASH_FORMAT)
		dayOfCalc, _ := GetTimeArr(record.InterestCalculateStartDate, record.InterestCalculateEndDate)
		record.DaysOfInterestCalculate = int(dayOfCalc)
		planRecords = append(planRecords, record)

	}
	log.Debugsf("total period = %d", period)

	return period, planRecords
}

func calcNextBillDateOrRepaymentDate(nextDateTs time.Time, requestDay int) (newNextDateTs time.Time) {
	log.Debugsf("nextDateTs = %s", nextDateTs)
	if requestDay == 29 || requestDay == 30 {
		firstMonDate := GetFirstDateOfMonthTime(nextDateTs)
		firstAddMonthTs := firstMonDate.AddDate(0, 1, 0)
		monthStr := firstAddMonthTs.Month()
		year := firstAddMonthTs.Year()
		log.Debugsf("year: %v", year)
		if monthStr == 02 {
			secondMonthDays := 0
			if !IsALeapYear(year) {
				secondMonthDays = 28
			} else {
				secondMonthDays = 29
			}
			newNextDateTs = firstAddMonthTs.AddDate(0, 0, secondMonthDays-1)

		} else {
			newNextDateTs = firstAddMonthTs.AddDate(0, 0, requestDay-1)
		}
	} else if requestDay == 31 {
		firstMonDate := GetFirstDateOfMonthTime(nextDateTs)
		firstAddMonthTs := firstMonDate.AddDate(0, 1, 0)
		newNextDateTs = GetLastDateOfMonth(firstAddMonthTs)
	} else {
		newNextDateTs = nextDateTs.AddDate(0, 1, 0)
	}
	log.Debugsf("end newNextDateTs = %s", newNextDateTs)
	return newNextDateTs
}

func calcNextDateWithMonth(request *Request) (nextDate string, err error) {
	requestDay := request.BillDay
	month := 1
	interestCalculateStartDate, e := time.Parse(DATE_DASH_FORMAT, request.InterestCalculateStartDate)
	if e != nil {
		return "", errors.New(constant.SYSPANIC, e)
	}
	startDay := interestCalculateStartDate.Day()
	log.Debugsf("startDay: %+v", startDay)

	requestDayInput, e := stringToInt(requestDay)
	if e != nil {
		return "", errors.New(constant.SYSPANIC, e)
	}

	log.Debugsf("StartDate startDay: %d", startDay)
	log.Debugsf("requestDayInput %d", requestDayInput)
	if startDay > requestDayInput {
		if requestDayInput == 29 || requestDayInput == 30 || requestDayInput == 31 {
			firstDayTs := GetFirstDateOfMonthTime(interestCalculateStartDate)
			nextMonthFirstDayTs := firstDayTs.AddDate(0, month, 0)
			firstAddDt := nextMonthFirstDayTs.Format(constant.DATE_DASH_FORMAT)
			monthStr := MonthToInt(nextMonthFirstDayTs.Month())
			year := nextMonthFirstDayTs.Year()
			log.Debugsf("year: %v", year)
			if monthStr == 2 {
				if !IsALeapYear(year) {
					nextDate = nextMonthFirstDayTs.AddDate(0, 0, 28-1).Format(constant.DATE_DASH_FORMAT)
					//nextDate = firstAddDt[0:len(firstAddDt)-2] + "28"
				} else {
					nextDate = nextMonthFirstDayTs.AddDate(0, 0, 29-1).Format(constant.DATE_DASH_FORMAT)
					//nextDate = firstAddDt[0:len(firstAddDt)-2] + "29"
				}
			} else {
				if requestDayInput == 31 {
					lastDateOfMonthTs := GetLastDateOfMonth(nextMonthFirstDayTs)
					lastDay := lastDateOfMonthTs.Day()
					if lastDay == 30 {
						nextDate = nextMonthFirstDayTs.AddDate(0, 0, 30-1).Format(constant.DATE_DASH_FORMAT)
						//nextDate = firstAddDt[0:len(firstAddDt)-2] + "30"
					} else {
						nextDate = nextMonthFirstDayTs.AddDate(0, 0, requestDayInput-1).Format(constant.DATE_DASH_FORMAT)
						//nextDate = firstAddDt[0:len(firstAddDt)-2] + requestDay
					}
				} else {
					nextDate = firstAddDt[0:len(firstAddDt)-2] + requestDay
				}
			}
		} else {
			firstDayTs := GetFirstDateOfMonthTime(interestCalculateStartDate)
			addDayTs := firstDayTs.AddDate(0, month, requestDayInput-1)
			nextDate = addDayTs.Format(constant.DATE_DASH_FORMAT)
		}
	} else {
		if requestDayInput == 29 || requestDayInput == 30 || requestDayInput == 31 {

			firstDayTs := GetFirstDateOfMonthTime(interestCalculateStartDate)
			nextMonthFirstDayTs := firstDayTs.AddDate(0, month, 0)
			firstAddDt := nextMonthFirstDayTs.Format(constant.DATE_DASH_FORMAT)
			monthStr := MonthToInt(nextMonthFirstDayTs.Month())
			year := nextMonthFirstDayTs.Year()
			log.Debugsf("year: %v", year)
			if monthStr == 2 {
				if !IsALeapYear(year) {
					nextDate = nextMonthFirstDayTs.AddDate(0, 0, 28-1).Format(constant.DATE_DASH_FORMAT)
					//nextDate = firstAddDt[0:len(firstAddDt)-2] + "28"
				} else {
					nextDate = nextMonthFirstDayTs.AddDate(0, 0, 29-1).Format(constant.DATE_DASH_FORMAT)
					//nextDate = firstAddDt[0:len(firstAddDt)-2] + "29"
				}
			} else {
				if requestDayInput == 31 {
					lastDateOfMonthTs := GetLastDateOfMonth(nextMonthFirstDayTs)
					lastDay := lastDateOfMonthTs.Day()
					if lastDay == 30 {
						//nextDate = firstAddDt[0:len(firstAddDt)-2] + "30"
						nextDate = nextMonthFirstDayTs.AddDate(0, 0, 30-1).Format(constant.DATE_DASH_FORMAT)
					} else {
						//nextDate = firstAddDt[0:len(firstAddDt)-2] + requestDay
						nextDate = nextMonthFirstDayTs.AddDate(0, 0, requestDayInput-1).Format(constant.DATE_DASH_FORMAT)
					}
				} else {
					nextDate = firstAddDt[0:len(firstAddDt)-2] + requestDay
				}
			}
		} else {
			firstDayTs := GetFirstDateOfMonthTime(interestCalculateStartDate)
			addDayTs := firstDayTs.AddDate(0, month, requestDayInput-1)
			nextDate = addDayTs.Format(constant.DATE_DASH_FORMAT)
		}
	}

	//如果计算出来的日期大于等于输入的到期日（结束日），nextDate 取到期日（结束日）
	//If the calculated date is greater than or equal to the input due date (InterestCalculateEndDate),
	//nextDate takes the due date (InterestCalculateEndDate)
	if nextDate >= request.InterestCalculateEndDate {
		nextDate = request.InterestCalculateEndDate
	}

	log.Debugsf("-- nextDate end: %v", nextDate)
	return nextDate, nil
}
