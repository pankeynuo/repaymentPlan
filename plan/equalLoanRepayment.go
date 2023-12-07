package main

/**
  *@Description TODO
  *@Author pauline
  *@Date 2023/12/5 10:09
**/
func equalLoanRepayment(request *RepayPlanRequest, response *RepayPlanResponse) (err error) {
	nextRepayDate, err := getFirstRepayDate(request.LoanStartDate, request.LoanEndDate, request.LoanCycleCode, request.RepayDay)
	totalPeriodNum, loanEndDate, err := getLoanEndDateAndTotalPeriodNum(request.PeriodNum, request.PeriodType, request.LoanEndDate, request.LoanStartDate, request.LoanCycleCode)
	dailyInterestRate := calculateDailyInterestRate(request.InterestRate, request.DaysOfYear)
	equalLoanPlan()
	return nil
}
func equalLoanPlan() {

}
