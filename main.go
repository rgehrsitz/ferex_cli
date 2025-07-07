package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"rgehrsitz/ferex_cli/pkg/config"
	"rgehrsitz/ferex_cli/pkg/calc"
	"rgehrsitz/ferex_cli/pkg/output"
)

var (
	cfgFile string
	verbose bool
	format  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ferex",
	Short: "Federal Retirement Explorer - CLI retirement planning tool",
	Long: `Federal Retirement Explorer (ferex) is a command-line tool for 
federal employees to calculate and analyze retirement scenarios.

Features:
- FERS and CSRS pension calculations
- TSP withdrawal modeling
- Social Security integration
- Tax calculations
- Excel/CSV export for analysis

Examples:
  ferex init --template basic > my-retirement.yaml
  ferex calc my-retirement.yaml
  ferex calc my-retirement.yaml --output results.csv
  ferex validate my-retirement.yaml`,
}

// calcCmd represents the calculate command
var calcCmd = &cobra.Command{
	Use:   "calc [config-file]",
	Short: "Calculate retirement projections",
	Long: `Calculate retirement projections based on a configuration file.

The config file should be in YAML format with all required fields.
Use 'ferex init' to generate a template configuration file.

Examples:
  ferex calc retirement-plan.yaml
  ferex calc plan.yaml --output results.csv --format csv
  ferex calc plan.yaml --verbose`,
	Args: cobra.ExactArgs(1),
	RunE: runCalc,
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a configuration file template",
	Long: `Generate a configuration file template with default values.
	
Templates available:
- basic: Basic FERS employee template
- advanced: Advanced template with all options
- csrs: CSRS employee template

Examples:
  ferex init > retirement-plan.yaml
  ferex init --template advanced > advanced-plan.yaml
  ferex init --template csrs > csrs-plan.yaml`,
	RunE: runInit,
}

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate [config-file]",
	Short: "Validate a configuration file",
	Long: `Validate a configuration file for required fields and correct values.

Can optionally fix common issues interactively.

Examples:
  ferex validate retirement-plan.yaml
  ferex validate plan.yaml --fix-interactive`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare [config-file]",
	Short: "Compare different retirement scenarios",
	Long: `Compare different retirement scenarios by varying key parameters.

Examples:
  ferex compare plan.yaml --ages 57,62
  ferex compare plan.yaml --ages 57,60,62 --output comparison.csv`,
	Args: cobra.ExactArgs(1),
	RunE: runCompare,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ferex.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "table", "output format (table, json, csv, yaml)")

	// Add subcommands
	rootCmd.AddCommand(calcCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(compareCmd)

	// calcCmd flags
	calcCmd.Flags().StringP("output", "o", "", "output file (default: stdout)")
	
	// initCmd flags
	initCmd.Flags().StringP("template", "t", "basic", "template type (basic, advanced, csrs)")
	
	// validateCmd flags
	validateCmd.Flags().Bool("fix-interactive", false, "interactively fix validation issues")
	
	// compareCmd flags
	compareCmd.Flags().StringSlice("ages", []string{"57", "62"}, "retirement ages to compare")
	compareCmd.Flags().StringP("output", "o", "", "output file (default: stdout)")
}

func runCalc(cmd *cobra.Command, args []string) error {
	configFile := args[0]
	
	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	// Validate configuration
	if err := config.ValidateConfig(cfg); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	
	// Run calculations
	calculator := calc.NewCalculator(cfg)
	results, err := calculator.Calculate()
	if err != nil {
		return fmt.Errorf("calculation failed: %w", err)
	}
	
	// Output results
	outputFile, _ := cmd.Flags().GetString("output")
	outputter := output.NewOutputter(format, outputFile, verbose)
	
	return outputter.OutputResults(results)
}

func runInit(cmd *cobra.Command, args []string) error {
	template, _ := cmd.Flags().GetString("template")
	
	cfg, err := config.GenerateTemplate(template)
	if err != nil {
		return fmt.Errorf("failed to generate template: %w", err)
	}
	
	outputter := output.NewOutputter("yaml", "", false)
	return outputter.OutputConfig(cfg)
}

func runValidate(cmd *cobra.Command, args []string) error {
	configFile := args[0]
	fixInteractive, _ := cmd.Flags().GetBool("fix-interactive")
	
	return config.ValidateConfigFile(configFile, fixInteractive)
}

func runCompare(cmd *cobra.Command, args []string) error {
	configFile := args[0]
	ages, _ := cmd.Flags().GetStringSlice("ages")
	outputFile, _ := cmd.Flags().GetString("output")
	
	// Load base configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	// Run comparison
	comparison, err := calc.CompareRetirementAges(cfg, ages)
	if err != nil {
		return fmt.Errorf("comparison failed: %w", err)
	}
	
	// Output results
	outputter := output.NewOutputter(format, outputFile, verbose)
	return outputter.OutputComparison(comparison)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}