package main

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/5 10:19
**/
func equalPrincipalAndInterest(request *Request) (response *Response, err error) {
	/*// 1.init
	sumPlanRepayTotalAmount := decimal.Decimal{}
	sumPlanRepayTotalInterest := decimal.Decimal{}
	records := make([]models.RepayPlanRecord, 0)
	remainPrinciple := request.LoanAmount
	flag := false
	specialStartDate := ""
	hasRepayPrincipal := decimal.Decimal{}

	// 2.calculate principal and interest amount
	// 2.1 get everyMonth need to repay principal amount A=loanAmount/periodNum
	planRepayPrinciplePeriod := request.LoanAmount.Div(decimal.NewFromInt(int64(request.PeriodNum))).Round(2)
	// 2.2 total days
	totalCalcDays, e := GetTimeArr(request.InterestCalculateStartDate, request.InterestCalculateEndDate)
	if e != nil {
		return nil, errors.New(constant.SYSPANIC, e)
	}
	// 2.3 everyMonth need to repay interest amount = LoanAmount*daysRate*totalDays/periodNum
	planRepayInterestPeriod := request.LoanAmount.Mul(rateD).
		Mul(decimal.NewFromInt(totalCalcDays)).Div(decimal.NewFromInt(int64(request.PeriodNum))).Round(2)

	// 3.make repay plan
	for i := 0; i < request.PeriodNum; i++ {
		log.Debugsf("period = %d", i)
		// 3.1 get InterestCalculateStartDate,InterestCalculateEndDate,repaymentDate begin
		reCalStartDate := ""
		reCalEndDate := ""
		repaymentDate := ""
		if i == 0 {
			reCalStartDate = request.InterestCalculateStartDate
			repaymentDate = nextRepayDate.Format(constant.DATE_DASH_FORMAT)
		} else {
			if request.LoanCycleCode == "03" {
				reCalBeginTs := nextRepayDate.AddDate(0, 0, (i-1)*14)
				reCalStartDate = time.Unix(reCalBeginTs.Unix(), 0).Format(constant.DATE_DASH_FORMAT)
				repaymentTs := nextRepayDate.AddDate(0, 0, i*14)
				repaymentDate = time.Unix(repaymentTs.Unix(), 0).Format(constant.DATE_DASH_FORMAT)
			}
			if request.LoanCycleCode == "04" {
				reCalStartDate, repaymentDate, flag, specialStartDate, err =
					calcPeriodStartDateAndEndDateByMonthly(i, nextRepayDate, flag, specialStartDate, request.RepayDay)
				if err != nil {
					return nil, errors.New(constant.SYSPANIC, err)
				}
			}
			if i == request.PeriodNum-1 {
				repaymentDate = request.InterestCalculateEndDate
				remainPrinciple = request.LoanAmount.Sub(hasRepayPrincipal)
			}
		}
		// 计息结束日=还款日的前一天
		reCalEndDateTs, _ := time.Parse(constant.DATE_DASH_FORMAT, repaymentDate)
		reCalEndDate = reCalEndDateTs.AddDate(0, 0, -1).Format(constant.DATE_DASH_FORMAT)
		// 3.1 get InterestCalculateStartDate,InterestCalculateEndDate,repaymentDate end

		record := models.RepayPlanRecord{
			PeriodNum:                  i + 1,          // current period num 当前期次的期数
			InterestCalculateStartDate: reCalStartDate, // current InterestCalculateStartDate 当前期次的开始计息日
			InterestCalculateEndDate:   reCalEndDate,   // current InterestCalculateEndDate 当前期次的结束计息日
			PlanRepayDate:              repaymentDate,  // current repaymentDate 当前期次的还款日
			PlanRepayInterest:          planRepayInterestPeriod,
		}
		// 3.2 get this period's interest begin 当前期次的计息天数
		days, err := GetTimeInterval(record.InterestCalculateStartDate, record.InterestCalculateEndDate)
		if err != nil {
			return nil, errors.New(constant.SYSPANIC, err)
		}
		daysOfInterestCal, err := conversionInt(days)
		if err != nil {
			return nil, errors.New(constant.SYSPANIC, err)
		}
		// 3.2 get this period's interest end

		record.DaysOfInterestCalculate = daysOfInterestCal

		if i == request.PeriodNum-1 {
			hasRepayPrincipal = hasRepayPrincipal.Add(remainPrinciple)
			record.PlanRepayPrinciple = remainPrinciple
			record.PlanRepayTotalAmount = planRepayInterestPeriod.Add(remainPrinciple)
			record.MaintainPrinciple = request.LoanAmount.Sub(hasRepayPrincipal)
		} else {
			hasRepayPrincipal = hasRepayPrincipal.Add(planRepayPrinciplePeriod)
			record.PlanRepayPrinciple = planRepayPrinciplePeriod
			record.PlanRepayTotalAmount = planRepayInterestPeriod.Add(planRepayPrinciplePeriod)
			remainPrinciple = remainPrinciple.Sub(planRepayPrinciplePeriod)
			record.MaintainPrinciple = remainPrinciple
		}
		log.Debugsf("hasRepayPrincipal =%s", hasRepayPrincipal)
		log.Debugsf("MaintainPrinciple =%s", record.MaintainPrinciple)
		sumPlanRepayTotalAmount = sumPlanRepayTotalAmount.Add(record.PlanRepayTotalAmount)
		sumPlanRepayTotalInterest = sumPlanRepayTotalInterest.Add(record.PlanRepayInterest)

		records = append(records, record)
	}
	response = &models.RepayPlanResponse{
		Records:                    records,
		RepayMethod:                request.RepayMethod,
		InterestCalculateStartDate: request.InterestCalculateStartDate,
		InterestCalculateEndDate:   request.InterestCalculateEndDate,
		PeriodNum:                  request.PeriodNum,
		LoanAmount:                 request.LoanAmount,
		InterestRate:               request.InterestRate,
		PlanRepayTotalAmount:       sumPlanRepayTotalAmount,
		PlanRepayTotalInterest:     sumPlanRepayTotalInterest,
	}*/
	return nil, nil
}
