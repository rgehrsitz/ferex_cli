package calc

import (
	"math"
	"time"

	"rgehrsitz/ferex_cli/internal/models"
)

// generateAnnualProjections creates year-by-year projections
func (c *Calculator) generateAnnualProjections(pension models.PensionCalculation, ss models.SocialSecurityCalculation, fersup models.FERSSupplementCalculation) ([]models.AnnualProjection, error) {
	var projections []models.AnnualProjection
	
	startAge := c.calculateAgeAtRetirement()
	endAge := 95 // Project to age 95
	
	// Initialize TSP balance (traditional + roth)
	tspBalance := c.config.TSP.TraditionalBalance + c.config.TSP.RothBalance
	
	for age := startAge; age <= endAge; age++ {
		currentAge := time.Now().Year() - c.config.Personal.BirthDate.Year()
		year := time.Now().Year() + (age - currentAge)
		
		projection := models.AnnualProjection{
			Year:             year,
			Age:              age,
			TSPStartBalance:  tspBalance,
		}
		
		// Calculate income sources
		projection.PensionIncome = c.calculatePensionIncome(pension, age, startAge)
		projection.FERSSupplementIncome = c.calculateFERSSupplementIncome(fersup, age)
		projection.SocialSecurityIncome = c.calculateSSIncome(ss, age)
		
		// Calculate TSP withdrawal
		projection.TSPWithdrawal = c.calculateTSPWithdrawal(tspBalance, age)
		
		// Update TSP balance
		tspGrowth := tspBalance * c.config.TSP.GrowthRate
		tspBalance = tspBalance + tspGrowth - projection.TSPWithdrawal
		if tspBalance < 0 {
			tspBalance = 0
		}
		
		projection.TSPGrowth = tspGrowth
		projection.TSPEndBalance = tspBalance
		
		// Calculate gross income
		projection.GrossIncome = projection.PensionIncome + 
			projection.FERSSupplementIncome + 
			projection.SocialSecurityIncome + 
			projection.TSPWithdrawal
		
		// Calculate taxes and deductions
		projection.FederalTax = c.calculateFederalTax(projection, age)
		projection.StateTax = c.calculateStateTax(projection, age)
		projection.HealthInsurance = c.calculateHealthInsurance(age)
		projection.LifeInsurance = c.calculateLifeInsurance(age)
		
		projection.TotalDeductions = projection.FederalTax + 
			projection.StateTax + 
			projection.HealthInsurance + 
			projection.LifeInsurance
		
		projection.NetIncome = projection.GrossIncome - projection.TotalDeductions
		
		// Apply COLA
		projection.COLARate = c.calculateCOLA(age, startAge)
		projection.InflationRate = 0.025 // 2.5% default inflation
		
		projections = append(projections, projection)
	}
	
	return projections, nil
}

// calculatePensionIncome calculates annual pension income with COLA
func (c *Calculator) calculatePensionIncome(pension models.PensionCalculation, currentAge, startAge int) float64 {
	basePension := pension.FinalPension
	
	// Apply COLA adjustments
	yearsRetired := currentAge - startAge
	if yearsRetired < 0 {
		return 0 // Not yet retired
	}
	
	// First year of retirement - no COLA yet
	if yearsRetired == 0 {
		return basePension
	}
	
	// FERS COLA eligibility - most FERS retirees don't get COLA until 62
	if c.config.Personal.RetirementSystem == "FERS" && currentAge < 62 {
		return basePension
	}
	
	// Apply compound COLA for subsequent years
	colaRate := 0.025 // 2.5% average
	if c.config.Personal.RetirementSystem == "FERS" {
		colaRate = c.calculateFERSCOLA(colaRate)
	}
	
	return basePension * math.Pow(1+colaRate, float64(yearsRetired))
}

// calculateFERSSupplementIncome calculates FERS Supplement income
func (c *Calculator) calculateFERSSupplementIncome(fersup models.FERSSupplementCalculation, currentAge int) float64 {
	if !fersup.Eligible || currentAge < fersup.StartAge || currentAge >= fersup.EndAge {
		return 0
	}
	
	return fersup.MonthlyAmount * 12
}

// calculateSSIncome calculates Social Security income
func (c *Calculator) calculateSSIncome(ss models.SocialSecurityCalculation, currentAge int) float64 {
	if currentAge < ss.ClaimingAge {
		return 0
	}
	
	// Apply COLA adjustments
	yearsReceiving := currentAge - ss.ClaimingAge
	if yearsReceiving <= 0 {
		return ss.MonthlyBenefit * 12
	}
	
	// Apply compound COLA (typically similar to general inflation)
	colaRate := 0.025 // 2.5% average
	return ss.MonthlyBenefit * 12 * math.Pow(1+colaRate, float64(yearsReceiving))
}

// calculateTSPWithdrawal calculates TSP withdrawal amount
func (c *Calculator) calculateTSPWithdrawal(balance float64, age int) float64 {
	if balance <= 0 {
		return 0
	}
	
	switch c.config.TSP.WithdrawalStrategy {
	case "fixed_amount":
		if c.config.TSP.WithdrawalAmount > 0 {
			return math.Min(c.config.TSP.WithdrawalAmount, balance)
		}
		return 0
		
	case "life_expectancy":
		// Use IRS life expectancy table
		lifeExpectancy := c.calculateLifeExpectancy(age)
		return balance / lifeExpectancy
		
	case "percentage":
		// Percentage of balance (e.g., 4% rule)
		if c.config.TSP.WithdrawalRate > 0 {
			return balance * c.config.TSP.WithdrawalRate
		}
		return balance * 0.04 // Default 4% rule
		
	case "lump_sum":
		// Take everything at retirement
		if age == c.calculateAgeAtRetirement() {
			return balance
		}
		return 0
		
	default:
		return 0
	}
}

// calculateLifeExpectancy calculates remaining life expectancy for TSP calculations
func (c *Calculator) calculateLifeExpectancy(age int) float64 {
	// Simplified IRS Uniform Lifetime Table
	switch {
	case age < 70:
		return 27.4
	case age < 75:
		return 24.7
	case age < 80:
		return 21.8
	case age < 85:
		return 19.1
	case age < 90:
		return 16.9
	case age < 95:
		return 14.8
	default:
		return 12.7
	}
}

// calculateFederalTax calculates federal income tax
func (c *Calculator) calculateFederalTax(projection models.AnnualProjection, age int) float64 {
	// Simplified federal tax calculation
	taxableIncome := projection.PensionIncome + projection.TSPWithdrawal
	
	// Add taxable portion of Social Security
	taxableIncome += c.calculateTaxableSS(projection.SocialSecurityIncome, projection.GrossIncome)
	
	// Apply standard deduction
	standardDeduction := 14700.0 // 2025 single standard deduction
	if age >= 65 {
		standardDeduction += 1850.0 // Additional standard deduction for seniors
	}
	
	taxableIncome -= standardDeduction
	if taxableIncome <= 0 {
		return 0
	}
	
	// Apply tax brackets (simplified)
	return c.calculateTaxBrackets(taxableIncome)
}

// calculateTaxableSS calculates taxable portion of Social Security
func (c *Calculator) calculateTaxableSS(ssBenefit, grossIncome float64) float64 {
	if ssBenefit == 0 {
		return 0
	}
	
	// Simplified provisional income calculation
	provisionalIncome := grossIncome - ssBenefit + (ssBenefit * 0.5)
	
	// Apply thresholds (single filer)
	if provisionalIncome <= 25000 {
		return 0
	}
	if provisionalIncome <= 34000 {
		return math.Min(ssBenefit*0.5, (provisionalIncome-25000)*0.5)
	}
	
	// Up to 85% taxable
	return math.Min(ssBenefit*0.85, (provisionalIncome-34000)*0.85+4500)
}

// calculateTaxBrackets applies federal tax brackets
func (c *Calculator) calculateTaxBrackets(income float64) float64 {
	// 2025 tax brackets (single filer)
	brackets := []struct {
		min  float64
		max  float64
		rate float64
	}{
		{0, 11000, 0.10},
		{11000, 44725, 0.12},
		{44725, 95375, 0.22},
		{95375, 182050, 0.24},
		{182050, 231250, 0.32},
		{231250, 578125, 0.35},
		{578125, math.Inf(1), 0.37},
	}
	
	var tax float64
	for _, bracket := range brackets {
		if income <= bracket.min {
			break
		}
		
		taxableInBracket := math.Min(income, bracket.max) - bracket.min
		tax += taxableInBracket * bracket.rate
	}
	
	return tax
}

// calculateStateTax calculates state income tax
func (c *Calculator) calculateStateTax(projection models.AnnualProjection, age int) float64 {
	// Use configured state tax rate if available
	if c.config.TaxInfo.StateTaxRate > 0 {
		taxableIncome := projection.GrossIncome
		
		// Apply exemptions for pension if configured
		if c.config.TaxInfo.PensionTaxExempt {
			taxableIncome -= projection.PensionIncome
		}
		
		// Apply exemptions for Social Security if configured
		if c.config.TaxInfo.SSTaxExempt {
			taxableIncome -= projection.SocialSecurityIncome
		}
		
		if taxableIncome <= 0 {
			return 0
		}
		
		return taxableIncome * c.config.TaxInfo.StateTaxRate
	}
	
	// Default state tax estimate based on known state patterns
	stateName := c.config.TaxInfo.State
	switch stateName {
	case "FL", "TX", "NV", "AK", "SD", "WY", "WA", "TN", "NH":
		return 0 // No state income tax
	case "PA":
		// PA taxes TSP but not pension
		return projection.TSPWithdrawal * 0.0307
	case "IL":
		// IL has flat 4.95% tax but exempts retirement income over 65
		if age >= 65 {
			return projection.TSPWithdrawal * 0.0495
		}
		return projection.GrossIncome * 0.0495
	default:
		// Default 5% state tax rate for unknown states
		return projection.GrossIncome * 0.05
	}
}

// calculateHealthInsurance calculates health insurance premiums
func (c *Calculator) calculateHealthInsurance(age int) float64 {
	startAge := c.calculateAgeAtRetirement()
	yearsRetired := age - startAge
	
	// Use configured premiums if available
	if c.config.HealthInsurance.RetirementPremium > 0 {
		basePremium := c.config.HealthInsurance.RetirementPremium
		
		// Apply COLA if specified
		if c.config.HealthInsurance.PremiumCOLA > 0 && yearsRetired > 0 {
			colaRate := c.config.HealthInsurance.PremiumCOLA
			return basePremium * math.Pow(1+colaRate, float64(yearsRetired))
		}
		
		return basePremium
	}
	
	// Default FEHB premium estimate
	basePremium := 4800.0 // $400/month
	
	// Apply default 3% annual increase
	if yearsRetired > 0 {
		return basePremium * math.Pow(1.03, float64(yearsRetired))
	}
	
	return basePremium
}

// calculateLifeInsurance calculates life insurance premiums
func (c *Calculator) calculateLifeInsurance(_ int) float64 {
	// Simplified FEGLI premium estimate
	return 600.0 // $50/month
}

// calculateCOLA calculates Cost of Living Adjustment
func (c *Calculator) calculateCOLA(_, _ int) float64 {
	// Simplified COLA calculation
	return 0.025 // 2.5% average
}

// calculateFERSCOLA applies FERS COLA rules
func (c *Calculator) calculateFERSCOLA(baseRate float64) float64 {
	// FERS COLA caps
	if baseRate <= 0.02 {
		return baseRate
	}
	if baseRate <= 0.03 {
		return 0.02
	}
	return baseRate - 0.01
}