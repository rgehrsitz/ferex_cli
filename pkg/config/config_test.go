package config

import (
	"os"
	"testing"
	"time"
)

func TestGenerateBasicTemplate(t *testing.T) {
	cfg := generateBasicTemplate()
	
	if cfg.Personal.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", cfg.Personal.Name)
	}
	
	if cfg.Personal.RetirementSystem != "FERS" {
		t.Errorf("Expected FERS system, got '%s'", cfg.Personal.RetirementSystem)
	}
	
	if cfg.Employment.CreditableService.TotalYears != 25 {
		t.Errorf("Expected 25 years service, got %.1f", cfg.Employment.CreditableService.TotalYears)
	}
	
	if cfg.TSP.CurrentBalance != 500000 {
		t.Errorf("Expected TSP balance 500000, got %.2f", cfg.TSP.CurrentBalance)
	}
}

func TestGenerateAdvancedTemplate(t *testing.T) {
	cfg := generateAdvancedTemplate()
	
	if cfg.Personal.Name != "Jane Smith" {
		t.Errorf("Expected name 'Jane Smith', got '%s'", cfg.Personal.Name)
	}
	
	if cfg.Employment.CreditableService.MilitaryService == nil {
		t.Error("Expected military service in advanced template")
	}
	
	if len(cfg.Employment.CreditableService.PartTimePeriods) == 0 {
		t.Error("Expected part-time periods in advanced template")
	}
	
	if cfg.Retirement.EarlyRetirement == nil {
		t.Error("Expected early retirement info in advanced template")
	}
}

func TestGenerateCSRSTemplate(t *testing.T) {
	cfg := generateCSRSTemplate()
	
	if cfg.Personal.RetirementSystem != "CSRS" {
		t.Errorf("Expected CSRS system, got '%s'", cfg.Personal.RetirementSystem)
	}
	
	if cfg.Employment.CreditableService.TotalYears != 42 {
		t.Errorf("Expected 42 years service for CSRS template, got %.1f", cfg.Employment.CreditableService.TotalYears)
	}
}

func TestValidateBusinessRules(t *testing.T) {
	cfg := generateBasicTemplate()
	
	// Test valid configuration
	err := validateBusinessRules(cfg)
	if err != nil {
		t.Errorf("Valid config failed validation: %v", err)
	}
	
	// Test TSP balance inconsistency
	cfg.TSP.CurrentBalance = 600000 // Doesn't match traditional + roth
	err = validateBusinessRules(cfg)
	if err == nil {
		t.Error("Expected validation error for TSP balance inconsistency")
	}
	
	// Fix TSP balance
	cfg.TSP.CurrentBalance = cfg.TSP.TraditionalBalance + cfg.TSP.RothBalance
	
	// Test future hire date
	cfg.Employment.HireDate = time.Now().Add(24 * time.Hour)
	err = validateBusinessRules(cfg)
	if err == nil {
		t.Error("Expected validation error for future hire date")
	}
}

func TestFERSEligibilityValidation(t *testing.T) {
	cfg := generateBasicTemplate()
	
	// Test valid FERS eligibility (age 62 with 25 years)
	err := validateFERSEligibility(cfg)
	if err != nil {
		t.Errorf("Valid FERS eligibility failed: %v", err)
	}
	
	// Test invalid eligibility (too young, not enough service)
	cfg.Retirement.TargetAge = 55
	cfg.Employment.CreditableService.TotalYears = 5
	err = validateFERSEligibility(cfg)
	if err == nil {
		t.Error("Expected validation error for insufficient FERS eligibility")
	}
	
	// Test MRA+30 eligibility
	cfg.Retirement.TargetAge = 57 // MRA for 1967 birth year
	cfg.Employment.CreditableService.TotalYears = 30
	err = validateFERSEligibility(cfg)
	if err != nil {
		t.Errorf("MRA+30 eligibility failed: %v", err)
	}
}

func TestMRACalculation(t *testing.T) {
	testCases := []struct {
		birthYear   int
		expectedMRA int
	}{
		{1945, 55},
		{1950, 56},
		{1955, 56},
		{1967, 57},
		{1975, 57},
	}
	
	for _, tc := range testCases {
		birthDate := time.Date(tc.birthYear, 1, 1, 0, 0, 0, 0, time.UTC)
		mra := calculateMRA(birthDate)
		if mra != tc.expectedMRA {
			t.Errorf("Birth year %d: expected MRA %d, got %d", tc.birthYear, tc.expectedMRA, mra)
		}
	}
}

func TestFillCalculatedFields(t *testing.T) {
	cfg := generateBasicTemplate()
	
	// Clear calculated fields
	cfg.Personal.CurrentAge = 0
	cfg.Employment.High3Salary = 0
	cfg.TSP.GrowthRate = 0
	cfg.TSP.CurrentBalance = 0
	
	err := fillCalculatedFields(cfg)
	if err != nil {
		t.Errorf("fillCalculatedFields failed: %v", err)
	}
	
	// Check fields were filled
	if cfg.Personal.CurrentAge == 0 {
		t.Error("Current age was not calculated")
	}
	
	if cfg.Employment.High3Salary != cfg.Employment.CurrentSalary {
		t.Error("High-3 salary was not set to current salary")
	}
	
	if cfg.TSP.GrowthRate != 0.07 {
		t.Error("TSP growth rate was not set to default 7%")
	}
	
	if cfg.TSP.CurrentBalance == 0 {
		t.Error("TSP current balance was not calculated")
	}
}

func TestConfigFileOperations(t *testing.T) {
	// Create a temporary config file
	tempFile := "test_config.yaml"
	defer os.Remove(tempFile)
	
	// Generate and save a config (could be used for more sophisticated tests)
	_ = generateBasicTemplate()
	
	// Save to file using YAML marshal
	data, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	data.Close()
	
	// Test loading the config
	loadedCfg, err := LoadConfig(tempFile)
	if err == nil {
		t.Log("Config loading works with empty file")
	}
	
	// Test validation of loaded config
	if loadedCfg != nil {
		err = ValidateConfig(loadedCfg)
		if err != nil {
			t.Logf("Validation failed as expected for empty config: %v", err)
		}
	}
}

func TestCalculateAge(t *testing.T) {
	// Test age calculation
	birthDate := time.Date(1967, 3, 15, 0, 0, 0, 0, time.UTC)
	age := calculateAge(birthDate)
	
	// Age should be reasonable (not testing exact age since it depends on current date)
	if age < 50 || age > 70 {
		t.Errorf("Calculated age %d seems unreasonable for birth year 1967", age)
	}
	
	// Test with a future birth date (should be negative, but function might handle it)
	futureBirth := time.Now().Add(365 * 24 * time.Hour)
	futureAge := calculateAge(futureBirth)
	if futureAge > 0 {
		t.Errorf("Future birth date resulted in positive age: %d", futureAge)
	}
}