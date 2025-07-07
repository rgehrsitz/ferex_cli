package models

import (
	"time"
)

// Config represents the complete retirement planning configuration
type Config struct {
	Personal       PersonalInfo       `yaml:"personal" validate:"required"`
	Employment     EmploymentInfo     `yaml:"employment" validate:"required"`
	Retirement     RetirementInfo     `yaml:"retirement" validate:"required"`
	TSP            TSPInfo            `yaml:"tsp" validate:"required"`
	SocialSecurity SocialSecurityInfo `yaml:"social_security"`
	HealthInsurance HealthInsuranceInfo `yaml:"health_insurance,omitempty"`
	TaxInfo        TaxInfo            `yaml:"tax_info,omitempty"`
	Output         OutputOptions      `yaml:"output,omitempty"`
}

// PersonalInfo contains basic personal information
type PersonalInfo struct {
	Name           string    `yaml:"name" validate:"required"`
	BirthDate      time.Time `yaml:"birth_date" validate:"required"`
	CurrentAge     int       `yaml:"current_age" validate:"required,min=18,max=100"`
	RetirementSystem string  `yaml:"retirement_system" validate:"required,oneof=FERS CSRS"`
}

// EmploymentInfo contains federal employment details
type EmploymentInfo struct {
	HireDate        time.Time `yaml:"hire_date" validate:"required"`
	CurrentSalary   float64   `yaml:"current_salary" validate:"required,gt=0"`
	High3Salary     float64   `yaml:"high_3_salary,omitempty" validate:"omitempty,gt=0"`
	CreditableService CreditableService `yaml:"creditable_service" validate:"required"`
}

// CreditableService represents service time calculations
type CreditableService struct {
	TotalYears      float64           `yaml:"total_years" validate:"required,gt=0"`
	PartTimePeriods []PartTimePeriod  `yaml:"part_time_periods,omitempty"`
	MilitaryService *MilitaryService  `yaml:"military_service,omitempty"`
	UnusedSickLeave float64           `yaml:"unused_sick_leave,omitempty" validate:"omitempty,gte=0"`
}

// PartTimePeriod represents a period of part-time employment
type PartTimePeriod struct {
	StartDate time.Time `yaml:"start_date" validate:"required"`
	EndDate   time.Time `yaml:"end_date" validate:"required"`
	HoursPerWeek float64 `yaml:"hours_per_week" validate:"required,gt=0,lte=40"`
}

// MilitaryService represents military service buy-back
type MilitaryService struct {
	Years     float64 `yaml:"years" validate:"required,gt=0"`
	BoughtBack bool   `yaml:"bought_back"`
}

// RetirementInfo contains retirement planning details
type RetirementInfo struct {
	TargetAge       int    `yaml:"target_age" validate:"required,min=50,max=70"`
	SurvivorBenefit string `yaml:"survivor_benefit" validate:"required,oneof=full partial none"`
	EarlyRetirement *EarlyRetirementInfo `yaml:"early_retirement,omitempty"`
}

// EarlyRetirementInfo contains early retirement options
type EarlyRetirementInfo struct {
	Type         string `yaml:"type" validate:"required,oneof=MRA+10 VERA DSR"`
	PostponedStart bool `yaml:"postponed_start,omitempty"`
}

// TSPInfo contains Thrift Savings Plan information
type TSPInfo struct {
	TraditionalBalance  float64 `yaml:"traditional_balance" validate:"required,gte=0"`
	RothBalance         float64 `yaml:"roth_balance" validate:"required,gte=0"`
	WithdrawalStrategy  string  `yaml:"withdrawal_strategy" validate:"required,oneof=fixed_amount life_expectancy lump_sum percentage"`
	WithdrawalAmount    float64 `yaml:"withdrawal_amount,omitempty" validate:"omitempty,gt=0"`
	WithdrawalRate      float64 `yaml:"withdrawal_rate,omitempty" validate:"omitempty,gt=0,lte=0.20"`
	GrowthRate          float64 `yaml:"growth_rate,omitempty" validate:"omitempty,gte=0,lte=0.15"`
}

// SocialSecurityInfo contains Social Security benefit information
type SocialSecurityInfo struct {
	EstimatedPIA float64 `yaml:"estimated_pia" validate:"required,gt=0"`
	ClaimingAge  int     `yaml:"claiming_age" validate:"required,min=62,max=70"`
	SpouseBenefit *SpouseBenefit `yaml:"spouse_benefit,omitempty"`
	// Optional: Monthly estimates from SS statement at different ages
	MonthlyEstimates map[int]float64 `yaml:"monthly_estimates,omitempty"`
}

// SpouseBenefit represents spouse Social Security information
type SpouseBenefit struct {
	EstimatedPIA float64 `yaml:"estimated_pia" validate:"required,gt=0"`
	ClaimingAge  int     `yaml:"claiming_age" validate:"required,min=62,max=70"`
}

// HealthInsuranceInfo contains health insurance premium information
type HealthInsuranceInfo struct {
	CurrentPremium    float64 `yaml:"current_premium,omitempty" validate:"omitempty,gte=0"`
	RetirementPremium float64 `yaml:"retirement_premium,omitempty" validate:"omitempty,gte=0"`
	PremiumCOLA       float64 `yaml:"premium_cola,omitempty" validate:"omitempty,gte=0,lte=0.10"`
	Plan              string  `yaml:"plan,omitempty"`
}

// TaxInfo contains state and tax-related information
type TaxInfo struct {
	State              string  `yaml:"state,omitempty"`
	StateTaxRate       float64 `yaml:"state_tax_rate,omitempty" validate:"omitempty,gte=0,lte=0.15"`
	PensionTaxExempt   bool    `yaml:"pension_tax_exempt,omitempty"`
	SSTaxExempt        bool    `yaml:"ss_tax_exempt,omitempty"`
	FilingStatus       string  `yaml:"filing_status,omitempty" validate:"omitempty,oneof=single mfj mfs hoh"`
}

// OutputOptions controls output formatting
type OutputOptions struct {
	Format     string `yaml:"format" validate:"omitempty,oneof=json csv yaml table"`
	Verbose    bool   `yaml:"verbose,omitempty"`
	OutputFile string `yaml:"output_file,omitempty"`
	Monthly    bool   `yaml:"monthly,omitempty"`
}