package calc

import (
	"testing"
	"time"

	"rgehrsitz/ferex_cli/internal/models"
)

// createTestConfig creates a basic test configuration
func createTestConfig() *models.Config {
	return &models.Config{
		Personal: models.PersonalInfo{
			Name:             "Test User",
			BirthDate:        time.Date(1967, 3, 15, 0, 0, 0, 0, time.UTC),
			CurrentAge:       57,
			RetirementSystem: "FERS",
		},
		Employment: models.EmploymentInfo{
			HireDate:      time.Date(1999, 1, 15, 0, 0, 0, 0, time.UTC),
			CurrentSalary: 85000,
			High3Salary:   82000,
			CreditableService: models.CreditableService{
				TotalYears:      25,
				PartTimePeriods: []models.PartTimePeriod{},
				MilitaryService: nil,
				UnusedSickLeave: 0,
			},
		},
		Retirement: models.RetirementInfo{
			TargetAge:       62,
			SurvivorBenefit: "full",
			EarlyRetirement: nil,
		},
		TSP: models.TSPInfo{
			TraditionalBalance: 400000,
			RothBalance:        100000,
			WithdrawalStrategy: "life_expectancy",
			GrowthRate:         0.07,
		},
		SocialSecurity: models.SocialSecurityInfo{
			EstimatedPIA: 2800,
			ClaimingAge:  67,
			SpouseBenefit: nil,
		},
	}
}

func TestFERSPensionCalculation(t *testing.T) {
	config := createTestConfig()
	calc := NewCalculator(config)
	
	pension, err := calc.calculatePension()
	if err != nil {
		t.Fatalf("calculatePension failed: %v", err)
	}
	
	// Test basic FERS calculation: 25 years * 82000 * 1.1% (age 62 with 20+ years)
	expectedBase := 25.0 * 82000.0 * 0.011
	if pension.BasePension != expectedBase {
		t.Errorf("Expected base pension %.2f, got %.2f", expectedBase, pension.BasePension)
	}
	
	// Test survivor benefit cost (10% for full survivor benefit)
	expectedSurvivorCost := expectedBase * 0.10
	if pension.SurvivorCost != expectedSurvivorCost {
		t.Errorf("Expected survivor cost %.2f, got %.2f", expectedSurvivorCost, pension.SurvivorCost)
	}
	
	// Test final pension after survivor benefit
	expectedFinal := expectedBase - expectedSurvivorCost
	if pension.FinalPension != expectedFinal {
		t.Errorf("Expected final pension %.2f, got %.2f", expectedFinal, pension.FinalPension)
	}
}

func TestFERSEarlyRetirementReduction(t *testing.T) {
	config := createTestConfig()
	config.Retirement.TargetAge = 57 // Early retirement at MRA
	config.Employment.CreditableService.TotalYears = 15 // Only 15 years, triggering MRA+10 reduction
	
	calc := NewCalculator(config)
	pension, err := calc.calculatePension()
	if err != nil {
		t.Fatalf("calculatePension failed: %v", err)
	}
	
	// For MRA+10 (age 57 with 15 years), retiring 5 years before 62
	// Should have 25% reduction (5 years under 62 * 5% per year)
	expectedReduction := 25.0
	if pension.ReductionPercent != expectedReduction {
		t.Errorf("Expected reduction %.1f%%, got %.1f%%", expectedReduction, pension.ReductionPercent)
	}
}

func TestSocialSecurityCalculation(t *testing.T) {
	config := createTestConfig()
	calc := NewCalculator(config)
	
	ss := calc.calculateSocialSecurity()
	
	// Test claiming at FRA (67) - should be 100% of PIA
	if ss.ClaimingAge != 67 {
		t.Errorf("Expected claiming age 67, got %d", ss.ClaimingAge)
	}
	
	if ss.PIA != 2800 {
		t.Errorf("Expected PIA 2800, got %.2f", ss.PIA)
	}
	
	if ss.Adjustment != 1.0 {
		t.Errorf("Expected adjustment 1.0 (100%%), got %.2f", ss.Adjustment)
	}
	
	if ss.MonthlyBenefit != 2800 {
		t.Errorf("Expected monthly benefit 2800, got %.2f", ss.MonthlyBenefit)
	}
}

func TestSocialSecurityEarlyClaiming(t *testing.T) {
	config := createTestConfig()
	config.SocialSecurity.ClaimingAge = 62 // Early claiming
	
	calc := NewCalculator(config)
	ss := calc.calculateSocialSecurity()
	
	// Should be reduced for claiming 5 years early
	if ss.Adjustment >= 1.0 {
		t.Errorf("Expected adjustment < 1.0 for early claiming, got %.2f", ss.Adjustment)
	}
	
	if ss.MonthlyBenefit >= 2800 {
		t.Errorf("Expected reduced benefit < 2800, got %.2f", ss.MonthlyBenefit)
	}
}

func TestPensionIncomeFirstYear(t *testing.T) {
	config := createTestConfig()
	calc := NewCalculator(config)
	
	pension, err := calc.calculatePension()
	if err != nil {
		t.Fatalf("calculatePension failed: %v", err)
	}
	
	// Test first year pension income (should equal final pension, no COLA yet)
	firstYearIncome := calc.calculatePensionIncome(pension, 62, 62) // Same age = first year
	if firstYearIncome != pension.FinalPension {
		t.Errorf("Expected first year income %.2f, got %.2f", pension.FinalPension, firstYearIncome)
	}
	
	// Test income before retirement (should be 0)
	preRetirement := calc.calculatePensionIncome(pension, 61, 62) // Age 61, retiring at 62
	if preRetirement != 0 {
		t.Errorf("Expected pre-retirement income 0, got %.2f", preRetirement)
	}
}

func TestTSPWithdrawals(t *testing.T) {
	config := createTestConfig()
	calc := NewCalculator(config)
	
	// Test life expectancy withdrawal at age 62
	withdrawal := calc.calculateTSPWithdrawal(500000, 62)
	
	// Should be based on life expectancy (approximately 27.4 years at age 62)
	expectedWithdrawal := 500000.0 / 27.4
	tolerance := 1000.0 // Allow some tolerance for life expectancy table differences
	
	if withdrawal < expectedWithdrawal-tolerance || withdrawal > expectedWithdrawal+tolerance {
		t.Errorf("Expected TSP withdrawal around %.2f, got %.2f", expectedWithdrawal, withdrawal)
	}
}

func TestFullCalculationFlow(t *testing.T) {
	config := createTestConfig()
	calc := NewCalculator(config)
	
	results, err := calc.Calculate()
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}
	
	// Test that we have projections
	if len(results.AnnualProjections) == 0 {
		t.Fatal("No annual projections generated")
	}
	
	// Test first year projection
	firstYear := results.AnnualProjections[0]
	
	// Should have pension income in first year
	if firstYear.PensionIncome <= 0 {
		t.Errorf("Expected pension income > 0 in first year, got %.2f", firstYear.PensionIncome)
	}
	
	// Should not have Social Security until claiming age
	if firstYear.SocialSecurityIncome != 0 {
		t.Errorf("Expected no Social Security in first year (age %d, claiming age %d), got %.2f", 
			firstYear.Age, config.SocialSecurity.ClaimingAge, firstYear.SocialSecurityIncome)
	}
	
	// Should have TSP withdrawal
	if firstYear.TSPWithdrawal <= 0 {
		t.Errorf("Expected TSP withdrawal > 0 in first year, got %.2f", firstYear.TSPWithdrawal)
	}
	
	// Net income should be positive
	if firstYear.NetIncome <= 0 {
		t.Errorf("Expected positive net income in first year, got %.2f", firstYear.NetIncome)
	}
}

func TestCSRSCalculation(t *testing.T) {
	config := createTestConfig()
	config.Personal.RetirementSystem = "CSRS"
	
	calc := NewCalculator(config)
	pension, err := calc.calculatePension()
	if err != nil {
		t.Fatalf("CSRS calculatePension failed: %v", err)
	}
	
	// CSRS has tiered calculation
	// First 5 years: 1.5%, next 5 years: 1.75%, remaining: 2.0%
	high3 := 82000.0
	expected := (5 * 0.015 * high3) + (5 * 0.0175 * high3) + (15 * 0.02 * high3)
	
	if pension.BasePension != expected {
		t.Errorf("Expected CSRS base pension %.2f, got %.2f", expected, pension.BasePension)
	}
}

func TestMRACalculation(t *testing.T) {
	config := createTestConfig()
	calc := NewCalculator(config)
	
	// Test birth year 1967 should have MRA of 57
	mra := calc.calculateMRA()
	if mra != 57 {
		t.Errorf("Expected MRA 57 for birth year 1967, got %d", mra)
	}
	
	// Test different birth year
	config.Personal.BirthDate = time.Date(1955, 1, 1, 0, 0, 0, 0, time.UTC)
	mra = calc.calculateMRA()
	if mra != 56 {
		t.Errorf("Expected MRA 56 for birth year 1955, got %d", mra)
	}
}