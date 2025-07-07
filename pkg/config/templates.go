package config

import (
	"time"
	"rgehrsitz/ferex_cli/internal/models"
)

// generateBasicTemplate creates a basic FERS employee template
func generateBasicTemplate() *models.Config {
	return &models.Config{
		Personal: models.PersonalInfo{
			Name:             "John Doe",
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
			CurrentBalance:     500000,
			TraditionalBalance: 400000,
			RothBalance:        100000,
			WithdrawalStrategy: "life_expectancy",
			WithdrawalAmount:   0,
			GrowthRate:         0.07,
		},
		SocialSecurity: models.SocialSecurityInfo{
			EstimatedPIA: 2800,
			ClaimingAge:  67,
			SpouseBenefit: nil,
		},
		Output: models.OutputOptions{
			Format:     "table",
			Verbose:    false,
			OutputFile: "",
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
	
	earlyRetirement := &models.EarlyRetirementInfo{
		Type:           "MRA+10",
		PostponedStart: false,
	}
	
	spouseBenefit := &models.SpouseBenefit{
		EstimatedPIA: 2200,
		ClaimingAge:  67,
	}
	
	return &models.Config{
		Personal: models.PersonalInfo{
			Name:             "Jane Smith",
			BirthDate:        time.Date(1965, 7, 22, 0, 0, 0, 0, time.UTC),
			CurrentAge:       59,
			RetirementSystem: "FERS",
		},
		Employment: models.EmploymentInfo{
			HireDate:      time.Date(1995, 6, 1, 0, 0, 0, 0, time.UTC),
			CurrentSalary: 95000,
			High3Salary:   92000,
			CreditableService: models.CreditableService{
				TotalYears:      29,
				PartTimePeriods: partTimePeriods,
				MilitaryService: militaryService,
				UnusedSickLeave: 240, // 30 days
			},
		},
		Retirement: models.RetirementInfo{
			TargetAge:       56, // Early retirement at MRA
			SurvivorBenefit: "partial",
			EarlyRetirement: earlyRetirement,
		},
		TSP: models.TSPInfo{
			CurrentBalance:     750000,
			TraditionalBalance: 550000,
			RothBalance:        200000,
			WithdrawalStrategy: "fixed_amount",
			WithdrawalAmount:   30000,
			GrowthRate:         0.08,
		},
		SocialSecurity: models.SocialSecurityInfo{
			EstimatedPIA: 3200,
			ClaimingAge:  67,
			SpouseBenefit: spouseBenefit,
		},
		Output: models.OutputOptions{
			Format:     "csv",
			Verbose:    true,
			OutputFile: "retirement-analysis.csv",
		},
	}
}

// generateCSRSTemplate creates a CSRS employee template
func generateCSRSTemplate() *models.Config {
	return &models.Config{
		Personal: models.PersonalInfo{
			Name:             "Robert Johnson",
			BirthDate:        time.Date(1958, 11, 3, 0, 0, 0, 0, time.UTC),
			CurrentAge:       66,
			RetirementSystem: "CSRS",
		},
		Employment: models.EmploymentInfo{
			HireDate:      time.Date(1982, 9, 15, 0, 0, 0, 0, time.UTC),
			CurrentSalary: 105000,
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
			TargetAge:       66,
			SurvivorBenefit: "full",
			EarlyRetirement: nil,
		},
		TSP: models.TSPInfo{
			CurrentBalance:     300000, // CSRS employees typically have less TSP
			TraditionalBalance: 250000,
			RothBalance:        50000,
			WithdrawalStrategy: "life_expectancy",
			WithdrawalAmount:   0,
			GrowthRate:         0.06,
		},
		SocialSecurity: models.SocialSecurityInfo{
			EstimatedPIA: 1800, // Typically lower for CSRS due to limited SS-covered employment
			ClaimingAge:  67,
			SpouseBenefit: nil,
		},
		Output: models.OutputOptions{
			Format:     "table",
			Verbose:    false,
			OutputFile: "",
		},
	}
}