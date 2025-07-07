package config

import (
	"fmt"
	"os"
	"time"

	"rgehrsitz/ferex_cli/internal/models"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// LoadConfig loads and validates a configuration file
func LoadConfig(filename string) (*models.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Fill in calculated fields if missing
	if err := fillCalculatedFields(&config); err != nil {
		return nil, fmt.Errorf("failed to calculate derived fields: %w", err)
	}

	return &config, nil
}

// ValidateConfig validates a configuration struct
func ValidateConfig(config *models.Config) error {
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Custom validation logic
	if err := validateBusinessRules(config); err != nil {
		return fmt.Errorf("business rule validation failed: %w", err)
	}

	return nil
}

// ValidateConfigFile validates a configuration file
func ValidateConfigFile(filename string, fixInteractive bool) error {
	config, err := LoadConfig(filename)
	if err != nil {
		return err
	}

	if err := ValidateConfig(config); err != nil {
		if fixInteractive {
			return interactiveValidationFix(config, filename, err)
		}
		return err
	}

	fmt.Printf("✓ Configuration file %s is valid\n", filename)
	return nil
}

// GenerateTemplate generates a configuration template
func GenerateTemplate(templateType string) (*models.Config, error) {
	switch templateType {
	case "basic":
		return generateBasicTemplate(), nil
	case "advanced":
		return generateAdvancedTemplate(), nil
	case "csrs":
		return generateCSRSTemplate(), nil
	default:
		return nil, fmt.Errorf("unknown template type: %s", templateType)
	}
}

// fillCalculatedFields fills in calculated fields that may be missing
func fillCalculatedFields(config *models.Config) error {
	// Always calculate total years of service from hire date to target retirement date
	serviceYears := calculateServiceYears(config.Employment.HireDate, config.Retirement.TargetRetirementDate)
	config.Employment.CreditableService.TotalYears = serviceYears

	// Set default TSP growth rate if not provided
	if config.TSP.GrowthRate == 0 {
		config.TSP.GrowthRate = 0.07 // 7% default
	}
	
	// Set default withdrawal rate for percentage strategy
	if config.TSP.WithdrawalStrategy == "percentage" && config.TSP.WithdrawalRate == 0 {
		config.TSP.WithdrawalRate = 0.04 // 4% default
	}
	
	// Set default health insurance COLA
	if config.HealthInsurance.PremiumCOLA == 0 && config.HealthInsurance.RetirementPremium > 0 {
		config.HealthInsurance.PremiumCOLA = 0.03 // 3% default
	}

	return nil
}

// validateBusinessRules validates business logic rules
// Optional fields (like early_retirement) may be omitted from the config YAML.
func validateBusinessRules(config *models.Config) error {
	// Check retirement age eligibility
	if config.Personal.RetirementSystem == "FERS" {
		if err := validateFERSEligibility(config); err != nil {
			return err
		}
	}

	// Validate TSP withdrawal strategy configuration
	switch config.TSP.WithdrawalStrategy {
	case "fixed_amount":
		if config.TSP.WithdrawalAmount <= 0 {
			return fmt.Errorf("fixed_amount strategy requires withdrawal_amount > 0")
		}
		if config.TSP.WithdrawalRate > 0 {
			return fmt.Errorf("withdrawal_rate should be zero for fixed_amount strategy")
		}
	case "percentage":
		if config.TSP.WithdrawalRate <= 0 || config.TSP.WithdrawalRate > 0.20 {
			return fmt.Errorf("percentage strategy requires withdrawal_rate between 0 and 0.20 (20%%)")
		}
		if config.TSP.WithdrawalAmount > 0 {
			return fmt.Errorf("withdrawal_amount should be zero for percentage strategy")
		}
	}

	// Check dates are logical
	if config.Employment.HireDate.After(time.Now()) {
		return fmt.Errorf("hire date cannot be in the future")
	}

	if config.Personal.BirthDate.After(config.Employment.HireDate) {
		return fmt.Errorf("birth date must be before hire date")
	}
	
	if config.Retirement.TargetRetirementDate.Before(config.Employment.HireDate) {
		return fmt.Errorf("retirement date must be after hire date")
	}

	return nil
}

// validateFERSEligibility validates FERS retirement eligibility
func validateFERSEligibility(config *models.Config) error {
	age := calculateAgeAtDate(config.Personal.BirthDate, config.Retirement.TargetRetirementDate)
	service := config.Employment.CreditableService.TotalYears

	// Check basic eligibility scenarios
	if age >= 62 && service >= 5 {
		return nil // Age 62 with 5+ years
	}
	if age >= 60 && service >= 20 {
		return nil // Age 60 with 20+ years
	}
	if service >= 30 {
		// MRA + 30 years (MRA varies by birth year)
		mra := calculateMRA(config.Personal.BirthDate)
		if age >= mra {
			return nil
		}
	}
	if service >= 10 {
		// MRA + 10 years (with reduction)
		mra := calculateMRA(config.Personal.BirthDate)
		if age >= mra {
			return nil
		}
	}

	return fmt.Errorf("FERS eligibility not met: age %d with %.1f years of service", age, service)
}

// calculateMRA calculates Minimum Retirement Age based on birth year
func calculateMRA(birthDate time.Time) int {
	birthYear := birthDate.Year()
	
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

// calculateAge calculates current age from birth date
func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	
	// Adjust if birthday hasn't occurred this year
	if now.Month() < birthDate.Month() || 
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}
	
	return age
}

// calculateServiceYears calculates years of service between hire and retirement dates
func calculateServiceYears(hireDate, retirementDate time.Time) float64 {
	duration := retirementDate.Sub(hireDate)
	years := duration.Hours() / (24 * 365.25) // Account for leap years
	return years
}

// calculateAgeAtDate calculates age at a specific date
func calculateAgeAtDate(birthDate, targetDate time.Time) int {
	years := targetDate.Year() - birthDate.Year()
	
	// Adjust if birthday hasn't occurred by target date
	if targetDate.Month() < birthDate.Month() || 
		(targetDate.Month() == birthDate.Month() && targetDate.Day() < birthDate.Day()) {
		years--
	}
	
	return years
}

// interactiveValidationFix attempts to fix validation issues interactively
func interactiveValidationFix(config *models.Config, filename string, validationErr error) error {
	fmt.Printf("Validation errors found in %s:\n", filename)
	fmt.Printf("Error: %v\n", validationErr)
	fmt.Printf("\nWould you like to try automatic fixes? (y/n): ")
	
	var response string
	fmt.Scanln(&response)
	
	if response == "y" || response == "Y" {
		// Try to fill missing fields
		if err := fillCalculatedFields(config); err != nil {
			return fmt.Errorf("failed to apply fixes: %w", err)
		}
		
		// Re-validate
		if err := ValidateConfig(config); err != nil {
			return fmt.Errorf("validation still fails after fixes: %w", err)
		}
		
		// Save the fixed config
		data, err := yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal fixed config: %w", err)
		}
		
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return fmt.Errorf("failed to write fixed config: %w", err)
		}
		
		fmt.Printf("✓ Configuration fixed and saved to %s\n", filename)
		return nil
	}
	
	return validationErr
}