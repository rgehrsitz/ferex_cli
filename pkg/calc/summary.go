package calc

import (
	"strconv"
	"time"

	"rgehrsitz/ferex_cli/internal/models"
)

// createSummary creates a retirement summary from calculations
func (c *Calculator) createSummary(pension models.PensionCalculation, ss models.SocialSecurityCalculation, fersup models.FERSSupplementCalculation, projections []models.AnnualProjection) models.RetirementSummary {
	summary := models.RetirementSummary{
		MonthlyPension:        pension.FinalPension / 12,
		AnnualPension:         pension.FinalPension,
		PensionReductionPct:   pension.ReductionPercent,
		SurvivorBenefitCost:   pension.SurvivorCost,
		NetMonthlyPension:     pension.FinalPension / 12,
		MonthlySocialSecurity: ss.MonthlyBenefit,
		SocialSecurityStartAge: ss.ClaimingAge,
		TSPStartingBalance:    c.config.TSP.CurrentBalance,
	}

	// FERS Supplement info
	if fersup.Eligible {
		summary.FERSSupplement = fersup.MonthlyAmount
		summary.SupplementEndAge = fersup.EndAge
	}

	// Calculate first year income and lifetime totals
	if len(projections) > 0 {
		summary.FirstYearIncome = projections[0].NetIncome
		summary.LifetimeIncome = c.calculateLifetimeIncome(projections)
		summary.ReplacementRatio = c.calculateReplacementRatio(projections[0])
	}

	// Find TSP depletion age
	summary.TSPProjectedDepletion = c.findTSPDepletionAge(projections)

	return summary
}

// createMetadata creates calculation metadata
func (c *Calculator) createMetadata() models.CalculationMetadata {
	return models.CalculationMetadata{
		CalculationDate:   time.Now(),
		ConfigVersion:     "1.0",
		CalculationEngine: "ferex-cli-v1.0",
		Assumptions: models.CalculationAssumptions{
			InflationRate:      0.025,
			TSPGrowthRate:      c.config.TSP.GrowthRate,
			LifeExpectancy:     95,
			FERSCOLARate:       0.025,
			SocialSecurityCOLA: 0.025,
			TaxBracketYear:     2025,
		},
		Warnings: c.generateWarnings(),
	}
}

// calculateLifetimeIncome sums projected lifetime income
func (c *Calculator) calculateLifetimeIncome(projections []models.AnnualProjection) float64 {
	var total float64
	for _, p := range projections {
		total += p.NetIncome
	}
	return total
}

// calculateReplacementRatio calculates income replacement ratio
func (c *Calculator) calculateReplacementRatio(firstYear models.AnnualProjection) float64 {
	preRetirementIncome := c.config.Employment.CurrentSalary
	return firstYear.NetIncome / preRetirementIncome
}

// findTSPDepletionAge finds when TSP balance reaches zero
func (c *Calculator) findTSPDepletionAge(projections []models.AnnualProjection) int {
	for _, p := range projections {
		if p.TSPEndBalance <= 0 && p.TSPStartBalance > 0 {
			return p.Age
		}
	}
	return 0 // TSP doesn't deplete within projection period
}

// generateWarnings generates calculation warnings
func (c *Calculator) generateWarnings() []string {
	var warnings []string

	// Check eligibility
	if !c.checkRetirementEligibility() {
		warnings = append(warnings, "Retirement eligibility requirements may not be met")
	}

	// Check TSP balance consistency
	if c.config.TSP.CurrentBalance != c.config.TSP.TraditionalBalance+c.config.TSP.RothBalance {
		warnings = append(warnings, "TSP balance components don't add up to total balance")
	}

	// Check if High-3 seems low
	if c.config.Employment.High3Salary < c.config.Employment.CurrentSalary*0.9 {
		warnings = append(warnings, "High-3 salary appears significantly lower than current salary")
	}

	// Check early retirement
	if c.config.Retirement.TargetAge < 62 {
		warnings = append(warnings, "Early retirement will result in reduced pension benefits")
	}

	return warnings
}

// checkRetirementEligibility performs basic eligibility check
func (c *Calculator) checkRetirementEligibility() bool {
	age := c.config.Retirement.TargetAge
	service := c.config.Employment.CreditableService.TotalYears

	if c.config.Personal.RetirementSystem == "FERS" {
		// Basic FERS eligibility
		if age >= 62 && service >= 5 {
			return true
		}
		if age >= 60 && service >= 20 {
			return true
		}
		mra := c.calculateMRA()
		if age >= mra && service >= 30 {
			return true
		}
		if age >= mra && service >= 10 {
			return true
		}
	} else {
		// Basic CSRS eligibility
		if age >= 62 && service >= 5 {
			return true
		}
		if age >= 60 && service >= 20 {
			return true
		}
		if age >= 55 && service >= 30 {
			return true
		}
	}

	return false
}

// CompareRetirementAges compares multiple retirement ages
func CompareRetirementAges(baseConfig *models.Config, ageStrings []string) (*models.ComparisonResults, error) {
	var results []models.RetirementResults
	
	for _, ageStr := range ageStrings {
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			return nil, err
		}
		
		// Create a copy of the config with modified retirement age
		configCopy := *baseConfig
		configCopy.Retirement.TargetAge = age
		
		// Calculate results for this age
		calc := NewCalculator(&configCopy)
		result, err := calc.Calculate()
		if err != nil {
			return nil, err
		}
		
		results = append(results, *result)
	}
	
	// Create comparison
	comparison := &models.ComparisonResults{
		Scenarios:    results,
		ComparisonMetrics: calculateComparisonMetrics(results),
	}
	
	return comparison, nil
}

// calculateComparisonMetrics calculates comparison metrics
func calculateComparisonMetrics(results []models.RetirementResults) models.ComparisonMetrics {
	if len(results) == 0 {
		return models.ComparisonMetrics{}
	}
	
	metrics := models.ComparisonMetrics{
		ScenarioCount: len(results),
	}
	
	// Find best/worst scenarios
	var bestLifetimeIncome, worstLifetimeIncome float64
	var bestReplacementRatio, worstReplacementRatio float64
	
	for i, result := range results {
		lifetime := result.Summary.LifetimeIncome
		replacement := result.Summary.ReplacementRatio
		
		if i == 0 || lifetime > bestLifetimeIncome {
			bestLifetimeIncome = lifetime
			metrics.BestLifetimeIncome = result.Summary
		}
		if i == 0 || lifetime < worstLifetimeIncome {
			worstLifetimeIncome = lifetime
		}
		
		if i == 0 || replacement > bestReplacementRatio {
			bestReplacementRatio = replacement
		}
		if i == 0 || replacement < worstReplacementRatio {
			worstReplacementRatio = replacement
		}
	}
	
	metrics.LifetimeIncomeSpread = bestLifetimeIncome - worstLifetimeIncome
	metrics.ReplacementRatioSpread = bestReplacementRatio - worstReplacementRatio
	
	return metrics
}