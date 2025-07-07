package models

import "time"

// RetirementResults contains the complete retirement calculation results
type RetirementResults struct {
	Summary        RetirementSummary  `json:"summary"`
	AnnualProjections []AnnualProjection `json:"annual_projections"`
	Metadata       CalculationMetadata `json:"metadata"`
}

// RetirementSummary provides key summary metrics
type RetirementSummary struct {
	// Basic pension information
	MonthlyPension       float64 `json:"monthly_pension"`
	AnnualPension        float64 `json:"annual_pension"`
	PensionReductionPct  float64 `json:"pension_reduction_pct,omitempty"`
	
	// Survivor benefit impact
	SurvivorBenefitCost  float64 `json:"survivor_benefit_cost,omitempty"`
	NetMonthlyPension    float64 `json:"net_monthly_pension"`
	
	// FERS Supplement (if applicable)
	FERSSupplement       float64 `json:"fers_supplement,omitempty"`
	SupplementEndAge     int     `json:"supplement_end_age,omitempty"`
	
	// Social Security
	MonthlySocialSecurity float64 `json:"monthly_social_security"`
	SocialSecurityStartAge int    `json:"social_security_start_age"`
	
	// TSP projections
	TSPStartingBalance   float64 `json:"tsp_starting_balance"`
	TSPProjectedDepletion int    `json:"tsp_projected_depletion,omitempty"`
	
	// Overall financial picture
	FirstYearIncome      float64 `json:"first_year_income"`
	LifetimeIncome       float64 `json:"lifetime_income"`
	ReplacementRatio     float64 `json:"replacement_ratio"`
}

// AnnualProjection represents one year of retirement income and expenses
type AnnualProjection struct {
	Year        int     `json:"year"`
	Age         int     `json:"age"`
	
	// Income sources
	PensionIncome     float64 `json:"pension_income"`
	FERSSupplementIncome float64 `json:"fers_supplement_income"`
	SocialSecurityIncome float64 `json:"social_security_income"`
	TSPWithdrawal     float64 `json:"tsp_withdrawal"`
	OtherIncome       float64 `json:"other_income"`
	GrossIncome       float64 `json:"gross_income"`
	
	// Taxes and deductions
	FederalTax        float64 `json:"federal_tax"`
	StateTax          float64 `json:"state_tax"`
	HealthInsurance   float64 `json:"health_insurance"`
	LifeInsurance     float64 `json:"life_insurance"`
	TotalDeductions   float64 `json:"total_deductions"`
	NetIncome         float64 `json:"net_income"`
	
	// TSP account status
	TSPStartBalance   float64 `json:"tsp_start_balance"`
	TSPGrowth         float64 `json:"tsp_growth"`
	TSPEndBalance     float64 `json:"tsp_end_balance"`
	
	// COLA adjustments
	COLARate          float64 `json:"cola_rate"`
	InflationRate     float64 `json:"inflation_rate"`
}

// CalculationMetadata provides information about the calculation
type CalculationMetadata struct {
	CalculationDate   time.Time `json:"calculation_date"`
	ConfigVersion     string    `json:"config_version"`
	CalculationEngine string    `json:"calculation_engine"`
	Assumptions       CalculationAssumptions `json:"assumptions"`
	Warnings          []string  `json:"warnings,omitempty"`
}

// CalculationAssumptions documents the assumptions used
type CalculationAssumptions struct {
	InflationRate     float64 `json:"inflation_rate"`
	TSPGrowthRate     float64 `json:"tsp_growth_rate"`
	LifeExpectancy    int     `json:"life_expectancy"`
	FERSCOLARate      float64 `json:"fers_cola_rate"`
	SocialSecurityCOLA float64 `json:"social_security_cola"`
	TaxBracketYear    int     `json:"tax_bracket_year"`
}

// ComparisonResults contains comparison analysis
type ComparisonResults struct {
	Scenarios         []RetirementResults `json:"scenarios"`
	ComparisonMetrics ComparisonMetrics   `json:"comparison_metrics"`
}

// ComparisonMetrics provides comparison statistics
type ComparisonMetrics struct {
	ScenarioCount           int               `json:"scenario_count"`
	BestLifetimeIncome      RetirementSummary `json:"best_lifetime_income"`
	LifetimeIncomeSpread    float64           `json:"lifetime_income_spread"`
	ReplacementRatioSpread  float64           `json:"replacement_ratio_spread"`
}

// Intermediate calculation models
type PensionCalculation struct {
	BasePension      float64
	ReductionPercent float64
	AdjustedPension  float64
	SurvivorCost     float64
	FinalPension     float64
}

type SocialSecurityCalculation struct {
	PIA            float64
	ClaimingAge    int
	Adjustment     float64
	MonthlyBenefit float64
}

type FERSSupplementCalculation struct {
	Eligible      bool
	MonthlyAmount float64
	StartAge      int
	EndAge        int
	FERSYears     float64
	SSEstimate    float64
}