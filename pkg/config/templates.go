package config

import (
	"rgehrsitz/ferex_cli/internal/models"
	"time"
)

// generateBasicTemplate creates a basic FERS employee template
func generateBasicTemplate() *models.Config {
	return &models.Config{
		Personal: models.PersonalInfo{
			Name:             "John Doe",
			BirthDate:        time.Date(1967, 3, 15, 0, 0, 0, 0, time.UTC),
			RetirementSystem: "FERS",
		},
		Employment: models.EmploymentInfo{
			HireDate:      time.Date(1999, 1, 15, 0, 0, 0, 0, time.UTC),
			High3Salary:   82000,
			CreditableService: models.CreditableService{
				TotalYears:      25,
				PartTimePeriods: []models.PartTimePeriod{},
				MilitaryService: nil,
				UnusedSickLeave: 0,
			},
		},
		Retirement: models.RetirementInfo{
			TargetRetirementDate: time.Date(2029, 3, 15, 0, 0, 0, 0, time.UTC), // Age 62
			SurvivorBenefit: "full",
			EarlyRetirement: nil,
		},
		TSP: models.TSPInfo{
			TraditionalBalance: 400000,
			RothBalance:        100000,
			WithdrawalStrategy: "percentage",
			WithdrawalRate:     0.04,
			GrowthRate:         0.07,
		},
		SocialSecurity: models.SocialSecurityInfo{
			EstimatedPIA: 2800,
			ClaimingAge:  67,
			SpouseBenefit: nil,
			MonthlyEstimates: map[int]float64{
				62: 2240, // Example: reduced benefit at 62
				67: 2800, // Full benefit at FRA
				70: 3472, // Delayed retirement credit at 70
			},
		},
		HealthInsurance: models.HealthInsuranceInfo{
			RetirementPremium: 4800,
			PremiumCOLA:       0.03,
			Plan:              "Blue Cross Standard",
		},
		TaxInfo: models.TaxInfo{
			State:            "VA",
			StateTaxRate:     0.05,
			PensionTaxExempt: false,
			SSTaxExempt:      false,
			FilingStatus:     "mfj",
		},
		Output: models.OutputOptions{
			Format:     "table",
			Verbose:    false,
			OutputFile: "",
			Monthly:    false,
		},
	}
}

// generateAdvancedTemplate creates an advanced template with all options
func generateAdvancedTemplate() *models.Config {
	militaryService := &models.MilitaryService{
		Years:      4,
		BoughtBack: true,
	}

	partTimePeriods := []models.PartTimePeriod{
		{
			StartDate:    time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2012, 12, 31, 0, 0, 0, 0, time.UTC),
			HoursPerWeek: 32,
		},
	}

	// EarlyRetirement is optional; set to nil if not applicable
	var earlyRetirement *models.EarlyRetirementInfo = nil
	// Uncomment below to include early retirement options
	// earlyRetirement := &models.EarlyRetirementInfo{
	// 	Type:           "MRA+10",
	// 	PostponedStart: false,
	// }

	spouseBenefit := &models.SpouseBenefit{
		EstimatedPIA: 2200,
		ClaimingAge:  67,
	}

	return &models.Config{
		Personal: models.PersonalInfo{
			Name:             "Jane Smith",
			BirthDate:        time.Date(1965, 7, 22, 0, 0, 0, 0, time.UTC),
			RetirementSystem: "FERS",
		},
		Employment: models.EmploymentInfo{
			HireDate:      time.Date(1995, 6, 1, 0, 0, 0, 0, time.UTC),
			High3Salary:   92000,
			CreditableService: models.CreditableService{
				// TotalYears is derived/calculated; do not set here
				PartTimePeriods: partTimePeriods,
				MilitaryService: militaryService,
				UnusedSickLeave: 240, // 30 days
			},
		},
		Retirement: models.RetirementInfo{
			TargetRetirementDate: time.Date(2021, 7, 22, 0, 0, 0, 0, time.UTC), // Age 56 (early retirement)
			SurvivorBenefit: "partial",
			EarlyRetirement: earlyRetirement, // Optional; set to nil if not needed
		},
		TSP: models.TSPInfo{
			TraditionalBalance: 550000,
			RothBalance:        200000,
			WithdrawalStrategy: "fixed_amount", // options: fixed_amount, percentage, life_expectancy, lump_sum
			WithdrawalAmount:   30000,           // set if strategy is fixed_amount, else 0
			WithdrawalRate:     0,               // set if strategy is percentage, else 0
			GrowthRate:         0.08,
		},
		SocialSecurity: models.SocialSecurityInfo{
			EstimatedPIA: 3200,
			ClaimingAge:  67,
			SpouseBenefit: spouseBenefit,
			MonthlyEstimates: map[int]float64{
				62: 2560,
				67: 3200,
				70: 3968,
			},
		},
		HealthInsurance: models.HealthInsuranceInfo{
			RetirementPremium: 6000,
			PremiumCOLA:       0.035,
			Plan:              "Blue Cross High Option",
		},
		TaxInfo: models.TaxInfo{
			State:            "FL", // No state income tax
			StateTaxRate:     0.0,
			PensionTaxExempt: false,
			SSTaxExempt:      false,
			FilingStatus:     "mfj",
		},
		Output: models.OutputOptions{
			Format:     "csv",
			Verbose:    true,
			OutputFile: "retirement-analysis.csv",
			Monthly:    false,
		},
	}
}

// generateCSRSTemplate creates a CSRS employee template
func generateCSRSTemplate() *models.Config {
	return &models.Config{
		Personal: models.PersonalInfo{
			Name:             "Robert Johnson",
			BirthDate:        time.Date(1958, 11, 3, 0, 0, 0, 0, time.UTC),
			RetirementSystem: "CSRS",
		},
		Employment: models.EmploymentInfo{
			HireDate:      time.Date(1982, 9, 15, 0, 0, 0, 0, time.UTC),
			High3Salary:   102000,
			CreditableService: models.CreditableService{
				TotalYears:      42,
				PartTimePeriods: []models.PartTimePeriod{},
				MilitaryService: &models.MilitaryService{
					Years:      6,
					BoughtBack: true,
				},
				UnusedSickLeave: 180, // 22.5 days
			},
		},
		Retirement: models.RetirementInfo{
			TargetRetirementDate: time.Date(2024, 11, 3, 0, 0, 0, 0, time.UTC), // Age 66
			SurvivorBenefit: "full",
			EarlyRetirement: nil,
		},
		TSP: models.TSPInfo{
			TraditionalBalance: 250000, // CSRS employees typically have less TSP
			RothBalance:        50000,
			WithdrawalStrategy: "life_expectancy",
			GrowthRate:         0.06,
		},
		SocialSecurity: models.SocialSecurityInfo{
			EstimatedPIA: 1800, // Typically lower for CSRS due to limited SS-covered employment
			ClaimingAge:  67,
			SpouseBenefit: nil,
			MonthlyEstimates: map[int]float64{
				62: 1440,
				67: 1800,
				70: 2232,
			},
		},
		HealthInsurance: models.HealthInsuranceInfo{
			RetirementPremium: 5200,
			PremiumCOLA:       0.03,
			Plan:              "FEHB Standard",
		},
		TaxInfo: models.TaxInfo{
			State:            "MD",
			StateTaxRate:     0.04,
			PensionTaxExempt: true,  // MD exempts some pension income
			SSTaxExempt:      false,
			FilingStatus:     "mfj",
		},
		Output: models.OutputOptions{
			Format:     "table",
			Verbose:    false,
			OutputFile: "",
			Monthly:    false,
		},
	}
}