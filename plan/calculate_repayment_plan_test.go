package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"testing"
)

/**
  *@Description 等额本息 按月还款
  *@Author pauline
  *@Date 2023/12/4 11:17
**/
func Test_fixedInstallmentMethod(t *testing.T) {
	request := &Request{
		LoanAmount:    decimal.NewFromFloat(400000),
		LoanStartDate: "2022-01-01",
		LoanEndDate:   "",
		InterestRate:  decimal.NewFromFloat(6),
		PeriodNum:     5,
		RepayDay:      1,
		LoanCycleCode: "03", // 01-日 02-两周 03-月 04-季 05-年
		RepayMethod:   "1",  // 1-等额本息 2-等额本金
		PeriodType:    "02", // 01-年 02-月
		DaysOfYear:    360,
	}
	resp, err := CalculateRepaymentPlan(request)
	if err != nil {
		t.Error(err)
	} else {
		repayMethod := getRepayMethod(resp.RepayMethod)
		printStr := "还款方式:" + repayMethod + "\n" +
			"还款频率:" + getLoanCycleCode(request.LoanCycleCode) + "    " +
			"年利息：" + resp.InterestRate.String() + "    " +
			"总期数：" + strconv.Itoa(resp.TotalPeriodNum) + "\n" +
			"日期：" + resp.LoanStartDate + " 至 " + resp.LoanEndDate + "\n" +
			"贷款金额：" + resp.LoanAmount.String() + "    " +
			"利息：" + resp.TotalInterest.String() + "    " +
			"总还款金额：" + resp.TotalRepayAmount.String() + "\n" +
			"期次    起息日       结息日       还款日   天数  本期还款本金 " +
			" 本期还款利息  本期还款总金额  剩余还款金额 \n"
		fmt.Printf(printStr)
		for _, item := range resp.PlanRepayRecords {
			fmt.Printf("%2s  %s  %s  %2s  %2s %10s %12s %11s %11s\n",
				strconv.Itoa(item.PeriodNum), item.PeriodStartDate, item.PeriodEndDate,
				item.PeriodRepayDate, strconv.Itoa(item.DaysOfPeriod), item.PeriodRepayPrinciple.String(),
				item.PeriodRepayInterest.String(), item.PeriodRepayTotalAmount.String(), item.MaintainPrinciple.String())
		}

	}
}

/**
  *@Description 等额本金 按月还款
  *@Author pauline
  *@Date 2023/12/4 11:17
**/
func Test_FixedPrincipalMethod(t *testing.T) {
	request := &Request{
		LoanAmount:    decimal.NewFromFloat(400000),
		LoanStartDate: "2022-01-01",
		LoanEndDate:   "",
		InterestRate:  decimal.NewFromFloat(6),
		PeriodNum:     5,
		RepayDay:      1,
		LoanCycleCode: "03", // 01-日 02-两周 03-月 04-季 05-年
		RepayMethod:   "2",  // 1-等额本息 2-等额本金
		PeriodType:    "02", // 01-年 02-月
		DaysOfYear:    360,
	}
	resp, err := CalculateRepaymentPlan(request)
	if err != nil {
		t.Error(err)
	} else {
		repayMethod := getRepayMethod(resp.RepayMethod)
		printStr := "还款方式:" + repayMethod + "\n" +
			"还款频率:" + getLoanCycleCode(request.LoanCycleCode) + "    " +
			"年利息：" + resp.InterestRate.String() + "    " +
			"总期数：" + strconv.Itoa(resp.TotalPeriodNum) + "\n" +
			"日期：" + resp.LoanStartDate + " 至 " + resp.LoanEndDate + "\n" +
			"贷款金额：" + resp.LoanAmount.String() + "    " +
			"利息：" + resp.TotalInterest.String() + "    " +
			"总还款金额：" + resp.TotalRepayAmount.String() + "\n" +
			"期次    起息日       结息日       还款日   天数  本期还款本金 " +
			" 本期还款利息  本期还款总金额  剩余还款金额 \n"
		fmt.Printf(printStr)
		for _, item := range resp.PlanRepayRecords {
			fmt.Printf("%2s  %s  %s  %2s  %2s %10s %12s %11s %11s\n",
				strconv.Itoa(item.PeriodNum), item.PeriodStartDate, item.PeriodEndDate,
				item.PeriodRepayDate, strconv.Itoa(item.DaysOfPeriod), item.PeriodRepayPrinciple.String(),
				item.PeriodRepayInterest.String(), item.PeriodRepayTotalAmount.String(), item.MaintainPrinciple.String())
		}

	}
}
func getRepayMethod(repayMethod string) string {
	switch repayMethod {
	case EqualLoanRepayment:
		return "等额本息"
	case EqualPrincipalRepayment:
		return "等额本金"
	case BothPrincipalAndInterest:
		return "息随本清"
	case BeforeInterestAfterPrincipal:
		return "先息后本"
	case EqualPrincipalAndInterest:
		return "等本等息"
	}
	return ""
}
func getLoanCycleCode(loanCycleCode string) string {
	switch loanCycleCode {
	case loanCycleFortnightly:
		return "两周"
	case loanCycleMonthly:
		return "月"
	}
	return ""
}
