package main

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/5 10:19
**/
func beforeInterestAfterPrincipal(request *Request) (response *Response, err error) {
	/*// 1.init
	sumPlanRepayTotalAmount := decimal.Decimal{}
	sumPlanRepayTotalInterest := decimal.Decimal{}
	records := make([]models.RepayPlanRecord, 0)
	flag := false
	specialStartDate := ""
	// 2.make repay plan
	for i := 0; i < request.PeriodNum; i++ {
		log.Debugsf("period = %d", i)
		// 2.1 get InterestCalculateStartDate,InterestCalculateEndDate,repaymentDate begin

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
			}
		}
		// 计息结束日=还款日的前一天
		reCalEndDateTs, _ := time.Parse(constant.DATE_DASH_FORMAT, repaymentDate)
		reCalEndDate = reCalEndDateTs.AddDate(0, 0, -1).Format(constant.DATE_DASH_FORMAT)
		// 2.1 get InterestCalculateStartDate,InterestCalculateEndDate,repaymentDate end

		record := models.RepayPlanRecord{
			PeriodNum:                  i + 1,
			InterestCalculateStartDate: reCalStartDate,
			InterestCalculateEndDate:   reCalEndDate,
			PlanRepayDate:              repaymentDate,
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

		// 当前期次的利息金额=贷款本金*计息天数*日利息
		planRepayInterest := request.LoanAmount.Mul(rateD).Mul(decimal.NewFromInt(days)).Round(2)
		log.Debugsf("planRepayInterest = %s", planRepayInterest)
		// 2.2 get this period's interest end
		record.DaysOfInterestCalculate = daysOfInterestCal
		record.PlanRepayInterest = planRepayInterest

		// 2.3. calculate the cumulative repayment amount, remaining repayment amount begin
		if i == request.PeriodNum-1 {
			record.PlanRepayPrinciple = request.LoanAmount
			record.PlanRepayTotalAmount = planRepayInterest.Add(request.LoanAmount)
			record.MaintainPrinciple = decimal.NewFromInt(0)
		} else {
			record.PlanRepayPrinciple = decimal.NewFromFloat(0)
			record.PlanRepayTotalAmount = planRepayInterest
			record.MaintainPrinciple = request.LoanAmount
		}
		// 2.3 alculate the cumulative repayment amount, remaining repayment amount end
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
