package calc

import (
	"fmt"
	"math"

	"rgehrsitz/ferex_cli/internal/models"
)

// Calculator handles retirement calculations
type Calculator struct {
	config *models.Config
}

// NewCalculator creates a new calculator instance
func NewCalculator(config *models.Config) *Calculator {
	return &Calculator{config: config}
}

// Calculate performs the complete retirement calculation
func (c *Calculator) Calculate() (*models.RetirementResults, error) {
	// Calculate basic pension
	pension, err := c.calculatePension()
	if err != nil {
		return nil, fmt.Errorf("pension calculation failed: %w", err)
	}

	// Calculate Social Security
	socialSecurity := c.calculateSocialSecurity()

	// Calculate FERS Supplement if applicable
	ferssupplement := c.calculateFERSSupplement()

	// Generate annual projections
	projections, err := c.generateAnnualProjections(pension, socialSecurity, ferssupplement)
	if err != nil {
		return nil, fmt.Errorf("projection generation failed: %w", err)
	}

	// Create summary
	summary := c.createSummary(pension, socialSecurity, ferssupplement, projections)

	// Create metadata
	metadata := c.createMetadata()

	return &models.RetirementResults{
		Summary:           summary,
		AnnualProjections: projections,
		Metadata:          metadata,
	}, nil
}

// calculatePension calculates the basic FERS/CSRS pension
func (c *Calculator) calculatePension() (models.PensionCalculation, error) {
	service := c.config.Employment.CreditableService.TotalYears
	high3 := c.config.Employment.High3Salary
	age := c.config.Retirement.TargetAge

	var basePension float64
	var reductionPct float64

	if c.config.Personal.RetirementSystem == "FERS" {
		basePension = c.calculateFERSPension(service, high3, age)
		reductionPct = c.calculateFERSReduction(age, service)
	} else {
		basePension = c.calculateCSRSPension(service, high3)
		reductionPct = c.calculateCSRSReduction(age, service)
	}

	// Apply reduction
	adjustedPension := basePension * (1 - reductionPct/100)

	// Apply survivor benefit reduction
	survivorCost := c.calculateSurvivorBenefitCost(adjustedPension)
	finalPension := adjustedPension - survivorCost

	return models.PensionCalculation{
		BasePension:      basePension,
		ReductionPercent: reductionPct,
		AdjustedPension:  adjustedPension,
		SurvivorCost:     survivorCost,
		FinalPension:     finalPension,
	}, nil
}

// calculateFERSPension calculates basic FERS pension
func (c *Calculator) calculateFERSPension(service, high3 float64, age int) float64 {
	var multiplier float64
	
	// Determine multiplier based on age and service
	if age >= 62 && service >= 20 {
		multiplier = 0.011 // 1.1% for age 62+ with 20+ years
	} else {
		multiplier = 0.01  // 1.0% for all other cases
	}
	
	return high3 * multiplier * service
}

// calculateFERSReduction calculates early retirement reduction for FERS
func (c *Calculator) calculateFERSReduction(age int, service float64) float64 {
	// No reduction for unreduced retirement
	if age >= 62 && service >= 5 {
		return 0 // Age 62 with 5+ years
	}
	if age >= 60 && service >= 20 {
		return 0 // Age 60 with 20+ years
	}
	
	// MRA + 30 has no reduction
	mra := c.calculateMRA()
	if age >= mra && service >= 30 {
		return 0
	}
	
	// MRA + 10 has reduction if starting before age 62
	if age >= mra && service >= 10 {
		if age < 62 {
			yearsUnder62 := 62 - age
			return float64(yearsUnder62) * 5.0 // 5% per year under 62
		}
		return 0
	}
	
	return 0 // Should not reach here for eligible retirees
}

// calculateCSRSPension calculates basic CSRS pension
func (c *Calculator) calculateCSRSPension(service, high3 float64) float64 {
	// CSRS has a tiered calculation
	var pension float64
	
	// First 5 years: 1.5%
	first5 := math.Min(service, 5) * 0.015 * high3
	pension += first5
	
	// Next 5 years (6-10): 1.75%
	if service > 5 {
		next5 := math.Min(service-5, 5) * 0.0175 * high3
		pension += next5
	}
	
	// Remaining years: 2.0%
	if service > 10 {
		remaining := (service - 10) * 0.02 * high3
		pension += remaining
	}
	
	return pension
}

// calculateCSRSReduction calculates early retirement reduction for CSRS
func (c *Calculator) calculateCSRSReduction(age int, service float64) float64 {
	// CSRS reductions are more complex - simplified here
	if age >= 62 {
		return 0 // No reduction at 62+
	}
	if age >= 60 && service >= 20 {
		return 0 // No reduction for 60+20
	}
	if age >= 55 && service >= 30 {
		return 0 // No reduction for 55+30
	}
	
	// Early retirement reduction (simplified)
	if age >= 55 && service >= 20 {
		yearsUnder62 := 62 - age
		return math.Min(float64(yearsUnder62)*2.0, 25.0) // 2% per year, max 25%
	}
	
	return 0
}

// calculateSurvivorBenefitCost calculates the cost of survivor benefits
func (c *Calculator) calculateSurvivorBenefitCost(pension float64) float64 {
	switch c.config.Retirement.SurvivorBenefit {
	case "full":
		if c.config.Personal.RetirementSystem == "FERS" {
			return pension * 0.10 // 10% reduction for full survivor benefit
		} else {
			// CSRS has more complex calculation
			return c.calculateCSRSSurvivorCost(pension)
		}
	case "partial":
		if c.config.Personal.RetirementSystem == "FERS" {
			return pension * 0.05 // 5% reduction for partial survivor benefit
		} else {
			// CSRS partial survivor benefit calculation
			return c.calculateCSRSSurvivorCost(pension) * 0.5
		}
	default:
		return 0 // No survivor benefit
	}
}

// calculateCSRSSurvivorCost calculates CSRS survivor benefit cost
func (c *Calculator) calculateCSRSSurvivorCost(pension float64) float64 {
	// CSRS: 2.5% of first $3600 + 10% of remainder
	if pension <= 3600 {
		return pension * 0.025
	}
	return 3600*0.025 + (pension-3600)*0.10
}

// calculateMRA calculates Minimum Retirement Age based on birth year
func (c *Calculator) calculateMRA() int {
	birthYear := c.config.Personal.BirthDate.Year()
	
	switch {
	case birthYear < 1948:
		return 55
	case birthYear < 1953:
		// 1948-1952: increases from 55 to 56 gradually, simplified to 56 for 1950+
		if birthYear < 1950 {
			return 55
		}
		return 56
	case birthYear < 1965:
		return 56
	case birthYear < 1970:
		return 57
	default:
		return 57
	}
}

// calculateSocialSecurity calculates Social Security benefits
func (c *Calculator) calculateSocialSecurity() models.SocialSecurityCalculation {
	pia := c.config.SocialSecurity.EstimatedPIA
	claimingAge := c.config.SocialSecurity.ClaimingAge
	
	var monthlyBenefit float64
	var adjustment float64
	
	// Use monthly estimates if available
	if c.config.SocialSecurity.MonthlyEstimates != nil {
		if estimate, exists := c.config.SocialSecurity.MonthlyEstimates[claimingAge]; exists {
			monthlyBenefit = estimate
			adjustment = estimate / pia // Calculate effective adjustment
		} else {
			// Fall back to calculated adjustment
			adjustment = c.calculateSSClaimingAdjustment(claimingAge)
			monthlyBenefit = pia * adjustment
		}
	} else {
		// Use calculated adjustment
		adjustment = c.calculateSSClaimingAdjustment(claimingAge)
		monthlyBenefit = pia * adjustment
	}
	
	return models.SocialSecurityCalculation{
		PIA:            pia,
		ClaimingAge:    claimingAge,
		Adjustment:     adjustment,
		MonthlyBenefit: monthlyBenefit,
	}
}

// calculateSSClaimingAdjustment calculates Social Security claiming age adjustment
func (c *Calculator) calculateSSClaimingAdjustment(claimingAge int) float64 {
	// Simplified - assumes FRA of 67
	fra := 67
	
	if claimingAge == fra {
		return 1.0 // 100% at FRA
	}
	if claimingAge < fra {
		// Reduction for early claiming
		monthsEarly := (fra - claimingAge) * 12
		if monthsEarly <= 36 {
			return 1.0 - (float64(monthsEarly) * 0.00555) // 5/9 of 1% per month
		}
		// Additional reduction for claiming more than 36 months early
		return 1.0 - (36*0.00555 + float64(monthsEarly-36)*0.00416) // 5/12 of 1% per month
	}
	
	// Delayed retirement credits
	monthsLate := (claimingAge - fra) * 12
	return 1.0 + (float64(monthsLate) * 0.00666) // 2/3 of 1% per month
}

// calculateFERSSupplement calculates FERS Supplement if applicable
func (c *Calculator) calculateFERSSupplement() models.FERSSupplementCalculation {
	// Only for FERS retirees under 62
	if c.config.Personal.RetirementSystem != "FERS" || c.config.Retirement.TargetAge >= 62 {
		return models.FERSSupplementCalculation{
			Eligible: false,
		}
	}
	
	// Check eligibility (simplified)
	service := c.config.Employment.CreditableService.TotalYears
	age := c.config.Retirement.TargetAge
	mra := c.calculateMRA()
	
	eligible := false
	if age >= mra && service >= 30 {
		eligible = true // MRA + 30
	}
	if age >= 60 && service >= 20 {
		eligible = true // Age 60 + 20
	}
	
	if !eligible {
		return models.FERSSupplementCalculation{
			Eligible: false,
		}
	}
	
	// Calculate supplement (simplified formula)
	ssEstimate := c.config.SocialSecurity.EstimatedPIA
	fersYears := service // Simplified - assumes all service is FERS
	supplement := (ssEstimate / 40) * fersYears
	
	return models.FERSSupplementCalculation{
		Eligible:        true,
		MonthlyAmount:   supplement,
		StartAge:        age,
		EndAge:          62,
		FERSYears:       fersYears,
		SSEstimate:      ssEstimate,
	}
}

