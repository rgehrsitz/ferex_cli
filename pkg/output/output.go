package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
	"rgehrsitz/ferex_cli/internal/models"
)

// Outputter handles various output formats
type Outputter struct {
	format     string
	outputFile string
	verbose    bool
}

// NewOutputter creates a new outputter
func NewOutputter(format, outputFile string, verbose bool) *Outputter {
	return &Outputter{
		format:     format,
		outputFile: outputFile,
		verbose:    verbose,
	}
}

// OutputResults outputs retirement calculation results
func (o *Outputter) OutputResults(results *models.RetirementResults) error {
	switch o.format {
	case "json":
		return o.outputJSON(results)
	case "csv":
		return o.outputCSV(results)
	case "yaml":
		return o.outputYAML(results)
	case "table":
		return o.outputTable(results)
	default:
		return fmt.Errorf("unsupported output format: %s", o.format)
	}
}

// OutputConfig outputs configuration as YAML
func (o *Outputter) OutputConfig(config *models.Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return o.writeOutput(string(data))
}

// OutputComparison outputs comparison results
func (o *Outputter) OutputComparison(comparison *models.ComparisonResults) error {
	switch o.format {
	case "json":
		data, err := json.MarshalIndent(comparison, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return o.writeOutput(string(data))
		
	case "csv":
		return o.outputComparisonCSV(comparison)
		
	case "table":
		return o.outputComparisonTable(comparison)
		
	default:
		return o.outputJSON(comparison)
	}
}

// outputJSON outputs results as JSON
func (o *Outputter) outputJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return o.writeOutput(string(jsonData))
}

// outputYAML outputs results as YAML
func (o *Outputter) outputYAML(data interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return o.writeOutput(string(yamlData))
}

// outputCSV outputs annual projections as CSV
func (o *Outputter) outputCSV(results *models.RetirementResults) error {
	var output string
	
	if o.outputFile != "" {
		file, err := os.Create(o.outputFile)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		return o.writeCSVData(writer, results)
	}

	// Output to stdout (convert to string format)
	headers := []string{
		"Year", "Age", "Pension Income", "FERS Supplement", "Social Security", 
		"TSP Withdrawal", "Gross Income", "Federal Tax", "State Tax", 
		"Total Deductions", "Net Income", "TSP Balance",
	}
	
	output = fmt.Sprintf("%s\n", joinStrings(headers, ","))
	
	for _, proj := range results.AnnualProjections {
		row := []string{
			strconv.Itoa(proj.Year),
			strconv.Itoa(proj.Age),
			fmt.Sprintf("%.2f", proj.PensionIncome),
			fmt.Sprintf("%.2f", proj.FERSSupplementIncome),
			fmt.Sprintf("%.2f", proj.SocialSecurityIncome),
			fmt.Sprintf("%.2f", proj.TSPWithdrawal),
			fmt.Sprintf("%.2f", proj.GrossIncome),
			fmt.Sprintf("%.2f", proj.FederalTax),
			fmt.Sprintf("%.2f", proj.StateTax),
			fmt.Sprintf("%.2f", proj.TotalDeductions),
			fmt.Sprintf("%.2f", proj.NetIncome),
			fmt.Sprintf("%.2f", proj.TSPEndBalance),
		}
		output += fmt.Sprintf("%s\n", joinStrings(row, ","))
	}

	return o.writeOutput(output)
}

// writeCSVData writes CSV data using csv.Writer
func (o *Outputter) writeCSVData(writer *csv.Writer, results *models.RetirementResults) error {
	// Write headers
	headers := []string{
		"Year", "Age", "Pension Income", "FERS Supplement", "Social Security", 
		"TSP Withdrawal", "Gross Income", "Federal Tax", "State Tax", 
		"Total Deductions", "Net Income", "TSP Balance",
	}
	
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	// Write data rows
	for _, proj := range results.AnnualProjections {
		row := []string{
			strconv.Itoa(proj.Year),
			strconv.Itoa(proj.Age),
			fmt.Sprintf("%.2f", proj.PensionIncome),
			fmt.Sprintf("%.2f", proj.FERSSupplementIncome),
			fmt.Sprintf("%.2f", proj.SocialSecurityIncome),
			fmt.Sprintf("%.2f", proj.TSPWithdrawal),
			fmt.Sprintf("%.2f", proj.GrossIncome),
			fmt.Sprintf("%.2f", proj.FederalTax),
			fmt.Sprintf("%.2f", proj.StateTax),
			fmt.Sprintf("%.2f", proj.TotalDeductions),
			fmt.Sprintf("%.2f", proj.NetIncome),
			fmt.Sprintf("%.2f", proj.TSPEndBalance),
		}
		
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

// outputTable outputs results as formatted table
func (o *Outputter) outputTable(results *models.RetirementResults) error {
	output := o.formatSummaryTable(results.Summary)
	
	if o.verbose {
		output += "\n\nDetailed Annual Projections:\n"
		output += o.formatProjectionTable(results.AnnualProjections)
	}

	return o.writeOutput(output)
}

// formatSummaryTable formats the retirement summary as a table
func (o *Outputter) formatSummaryTable(summary models.RetirementSummary) string {
	output := "Retirement Planning Summary\n"
	output += "===========================\n\n"
	
	output += fmt.Sprintf("Monthly Pension:           $%.2f\n", summary.MonthlyPension)
	output += fmt.Sprintf("Annual Pension:            $%.2f\n", summary.AnnualPension)
	
	if summary.PensionReductionPct > 0 {
		output += fmt.Sprintf("Pension Reduction:         %.1f%%\n", summary.PensionReductionPct)
	}
	
	if summary.SurvivorBenefitCost > 0 {
		output += fmt.Sprintf("Survivor Benefit Cost:     $%.2f/year\n", summary.SurvivorBenefitCost*12)
	}
	
	if summary.FERSSupplement > 0 {
		output += fmt.Sprintf("FERS Supplement:           $%.2f/month (until age %d)\n", 
			summary.FERSSupplement, summary.SupplementEndAge)
	}
	
	output += fmt.Sprintf("Social Security:           $%.2f/month (starting age %d)\n", 
		summary.MonthlySocialSecurity, summary.SocialSecurityStartAge)
	
	output += fmt.Sprintf("TSP Starting Balance:      $%.2f\n", summary.TSPStartingBalance)
	
	if summary.TSPProjectedDepletion > 0 {
		output += fmt.Sprintf("TSP Depletion Age:         %d\n", summary.TSPProjectedDepletion)
	}
	
	output += fmt.Sprintf("\nFirst Year Income:         $%.2f\n", summary.FirstYearIncome)
	output += fmt.Sprintf("Lifetime Income:           $%.2f\n", summary.LifetimeIncome)
	output += fmt.Sprintf("Replacement Ratio:         %.1f%%\n", summary.ReplacementRatio*100)
	
	return output
}

// formatProjectionTable formats annual projections as a table
func (o *Outputter) formatProjectionTable(projections []models.AnnualProjection) string {
	output := fmt.Sprintf("%-6s %-4s %-12s %-12s %-12s %-12s %-12s %-12s\n",
		"Year", "Age", "Pension", "SS", "TSP Withdraw", "Gross", "Net", "TSP Balance")
	output += fmt.Sprintf("%s\n", "------------------------------------------------------------------------------------")
	
	for i, proj := range projections {
		if i > 20 && !o.verbose { // Limit output unless verbose
			output += fmt.Sprintf("... (use --verbose for complete projection)\n")
			break
		}
		
		output += fmt.Sprintf("%-6d %-4d $%-11.0f $%-11.0f $%-11.0f $%-11.0f $%-11.0f $%-11.0f\n",
			proj.Year, proj.Age, proj.PensionIncome, proj.SocialSecurityIncome,
			proj.TSPWithdrawal, proj.GrossIncome, proj.NetIncome, proj.TSPEndBalance)
	}
	
	return output
}

// outputComparisonCSV outputs comparison results as CSV
func (o *Outputter) outputComparisonCSV(comparison *models.ComparisonResults) error {
	output := "Scenario,Retirement Age,Monthly Pension,Annual Pension,First Year Income,Lifetime Income,Replacement Ratio,TSP Depletion Age\n"
	
	for i, scenario := range comparison.Scenarios {
		row := fmt.Sprintf("Scenario %d,%d,%.2f,%.2f,%.2f,%.2f,%.2f,%d\n",
			i+1, 
			scenario.AnnualProjections[0].Age, // Retirement age
			scenario.Summary.MonthlyPension,
			scenario.Summary.AnnualPension,
			scenario.Summary.FirstYearIncome,
			scenario.Summary.LifetimeIncome,
			scenario.Summary.ReplacementRatio*100,
			scenario.Summary.TSPProjectedDepletion)
		output += row
	}
	
	return o.writeOutput(output)
}

// outputComparisonTable outputs comparison results as a table
func (o *Outputter) outputComparisonTable(comparison *models.ComparisonResults) error {
	output := "Retirement Age Comparison\n"
	output += "=========================\n\n"
	
	output += fmt.Sprintf("%-10s %-15s %-15s %-15s %-15s %-15s %-15s\n",
		"Age", "Monthly Pension", "Annual Pension", "First Yr Income", "Lifetime Income", "Replace Ratio", "TSP Depletion")
	output += "--------------------------------------------------------------------------------------------------------\n"
	
	for _, scenario := range comparison.Scenarios {
		retirementAge := scenario.AnnualProjections[0].Age
		
		output += fmt.Sprintf("%-10d $%-14.0f $%-14.0f $%-14.0f $%-14.0f %-14.1f%% %-14d\n",
			retirementAge,
			scenario.Summary.MonthlyPension,
			scenario.Summary.AnnualPension,
			scenario.Summary.FirstYearIncome,
			scenario.Summary.LifetimeIncome,
			scenario.Summary.ReplacementRatio*100,
			scenario.Summary.TSPProjectedDepletion)
	}
	
	output += "\nComparison Metrics:\n"
	output += fmt.Sprintf("Scenarios compared:        %d\n", comparison.ComparisonMetrics.ScenarioCount)
	output += fmt.Sprintf("Lifetime income spread:    $%.2f\n", comparison.ComparisonMetrics.LifetimeIncomeSpread)
	output += fmt.Sprintf("Replacement ratio spread:  %.1f%%\n", comparison.ComparisonMetrics.ReplacementRatioSpread*100)
	
	return o.writeOutput(output)
}

// writeOutput writes output to file or stdout
func (o *Outputter) writeOutput(content string) error {
	if o.outputFile != "" {
		return os.WriteFile(o.outputFile, []byte(content), 0644)
	}
	
	fmt.Print(content)
	return nil
}

// joinStrings joins a slice of strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	
	return result
}