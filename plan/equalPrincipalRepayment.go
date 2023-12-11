package main

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/5 10:18
**/
func equalPrincipalRepayment(request *Request) (response *Response, err error) {
	/*// 1.init
	records := make([]RepayPlanRecord, 0)
	remainPrinciple := request.LoanAmount
	flag := false
	specialStartDate := ""
	hasRepayPrincipal := decimal.Decimal{}
	sumPlanRepayTotalAmount := decimal.Decimal{}
	sumPlanRepayTotalInterest := decimal.Decimal{}
	// 2.everyMonth need to repay principle
	planRepayPrinciple := request.LoanAmount.Div(decimal.NewFromInt(int64(request.PeriodNum)))
	planRepayPrinciple = planRepayPrinciple.RoundBank(2)
	// 3.make repay plan
	for i := 0; i < request.PeriodNum; i++ {

		// 3.1 get InterestCalculateStartDate,InterestCalculateEndDate,repaymentDate begin
		reCalStartDate := ""
		reCalEndDate := ""
		repaymentDate := ""
		if i == 0 {
			reCalStartDate = request.InterestCalculateStartDate
			repaymentDate = nextRepayDate.Format(constant.DATE_DASH_FORMAT)
		} else {
			if request.LoanCycleCode == "03" {
				if i == 1 {
					reCalStartDate = nextRepayDate.Format(constant.DATE_DASH_FORMAT)
				} else {
					reCalBeginTs := nextRepayDate.AddDate(0, 0, (i-1)*14)
					reCalStartDate = time.Unix(reCalBeginTs.Unix(), 0).Format(constant.DATE_DASH_FORMAT)
				}

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
		reCalEndDateTs, _ := time.Parse(DATE_DASH_FORMAT, repaymentDate)
		reCalEndDate = reCalEndDateTs.AddDate(0, 0, -1).Format(DATE_DASH_FORMAT)
		// 3.1 get InterestCalculateStartDate,InterestCalculateEndDate,repaymentDate end
		record := models.RepayPlanRecord{
			PeriodNum:                  i + 1,          // current period num 当前期次的期数
			InterestCalculateStartDate: reCalStartDate, // current InterestCalculateStartDate 当前期次的开始计息日
			InterestCalculateEndDate:   reCalEndDate,   // current InterestCalculateEndDate 当前期次的结束计息日
			PlanRepayDate:              repaymentDate,  // current repaymentDate 当前期次的还款日
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

		planRepayInterest := remainPrinciple.Mul(rateD).Mul(decimal.NewFromInt(days))
		planRepayInterest = planRepayInterest.RoundBank(2)

		// 3.2 get this period's interest end
		record.PlanRepayInterest = planRepayInterest
		// current interest amount 当前期次的利息金额
		record.DaysOfInterestCalculate = daysOfInterestCal

		// 3.3. calculate the cumulative repayment amount, remaining repayment amount begin
		if i == request.PeriodNum-1 {
			hasRepayPrincipal = hasRepayPrincipal.Add(remainPrinciple)           // 累积已还本金=累积已还本金+上一期总的剩余还款本金
			record.PlanRepayPrinciple = remainPrinciple                          // 当前期次还款本金=上一期总的剩余还款本金
			record.PlanRepayTotalAmount = planRepayInterest.Add(remainPrinciple) // 当前期次的总还款金额=当前期次的利息金额+当前期次还款本金
			record.MaintainPrinciple = request.LoanAmount.Sub(hasRepayPrincipal)
		} else {
			hasRepayPrincipal = hasRepayPrincipal.Add(planRepayPrinciple) // 累积已还本金
			record.PlanRepayPrinciple = planRepayPrinciple                // 当前期次还款本金
			// current total repay amount=current repay principle+current interest amount
			// 当前期次的总还款金额=当前期次还款本金+当前期次的利息金额
			record.PlanRepayTotalAmount = planRepayInterest.Add(planRepayPrinciple)
			remainPrinciple = remainPrinciple.Sub(planRepayPrinciple) // 剩余还款本金
			record.MaintainPrinciple = remainPrinciple                // 剩余还款本金
		}
		// 3.3. calculate the cumulative repayment amount, remaining repayment amount end

		sumPlanRepayTotalAmount = sumPlanRepayTotalAmount.Add(record.PlanRepayTotalAmount)
		sumPlanRepayTotalInterest = sumPlanRepayTotalInterest.Add(record.PlanRepayInterest)
		records = append(records, record)
	}
	response = &RepayPlanResponse{
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
