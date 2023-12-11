package main

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/5 10:18
**/
func bothPrincipalAndInterest(request *Request) (response *Response, err error) {
	/*if request.InterestCalculateEndDate == "" {
		return nil, errors.New(constant.SYSPANIC, "MatureDate can not be null")
	}
	if request.InterestCalculateEndDate <= request.InterestCalculateStartDate {
		return nil, errors.New(constant.SYSPANIC, "The Drawdown Date cannot be less than or equal to the Maturity Date")
	}
	reCalEndDateTs, _ := time.Parse(constant.DATE_DASH_FORMAT, request.InterestCalculateEndDate)
	reCalEndDate := reCalEndDateTs.AddDate(0, 0, -1).Format(constant.DATE_DASH_FORMAT)

	// 1.Daily interest rate 日利率
	rateD := request.InterestRate.Div(decimal.NewFromInt(int64(request.DaysOfYear))).Div(decimal.NewFromFloat(100)) //日利率
	log.Debugsf("rateD: %+v", rateD)
	// 2.get this period's interest begin 计息天数
	days, e := GetTimeInterval(request.InterestCalculateStartDate, reCalEndDate)
	if e != nil {
		return nil, errors.New(constant.SYSPANIC, e)
	}
	daysOfInterestCal, e := conversionInt(days)
	if e != nil {
		return nil, errors.New(constant.SYSPANIC, e)
	}
	// 3.calculate total interest amount
	planRepayInterest := request.LoanAmount.Mul(rateD).Mul(decimal.NewFromInt(days)).Round(2)
	log.Debugsf("planRepayInterest = %s", planRepayInterest)
	// 4.response only one period's plan
	totalAmount := planRepayInterest.Add(request.LoanAmount)
	record := models.RepayPlanRecord{
		PeriodNum:                  1,
		InterestCalculateStartDate: request.InterestCalculateStartDate,
		InterestCalculateEndDate:   reCalEndDate,
		PlanRepayDate:              request.InterestCalculateEndDate,
		PlanRepayPrinciple:         request.LoanAmount,
		MaintainPrinciple:          decimal.NewFromFloat(0),
		DaysOfInterestCalculate:    daysOfInterestCal,
		PlanRepayInterest:          planRepayInterest,
		PlanRepayTotalAmount:       totalAmount,
	}
	response = &models.RepayPlanResponse{
		PeriodNum:                  1,
		RepayMethod:                request.RepayMethod,
		InterestCalculateStartDate: request.InterestCalculateStartDate,
		InterestCalculateEndDate:   request.InterestCalculateEndDate,
		LoanAmount:                 request.LoanAmount,
		InterestRate:               request.InterestRate,
		PlanRepayTotalAmount:       totalAmount,
		PlanRepayTotalInterest:     planRepayInterest,
		Records:                    []models.RepayPlanRecord{record},
	}*/
	return nil, nil
}
