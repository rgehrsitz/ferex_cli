# Ferex CLI - Federal Retirement Explorer

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A command-line retirement planning tool for U.S. federal employees. Ferex provides accurate calculations for FERS and CSRS pensions, TSP withdrawals, Social Security benefits, and tax projections.

## Features

- **Comprehensive FERS/CSRS Calculations**
  - Basic annuity computation with 1.0% and 1.1% multipliers
  - Early retirement reductions (MRA+10)
  - Survivor benefit cost analysis
  - FERS Annuity Supplement calculations

- **Social Security Integration**
  - Post-WEP/GPO repeal calculations (2025+)
  - Claiming age adjustments
  - Taxable benefit calculations

- **TSP Modeling**
  - Multiple withdrawal strategies (fixed amount, life expectancy)
  - Traditional vs. Roth tax implications
  - Required Minimum Distributions (RMDs)
  - Growth projections

- **Tax Calculations**
  - Federal income tax with 2025 brackets
  - State tax estimates
  - IRS Simplified Method for pension taxation
  - Social Security benefit taxation

- **Multiple Output Formats**
  - Table: Human-readable summary reports
  - CSV: Excel/Google Sheets integration
  - JSON: Data processing and APIs
  - YAML: Configuration management

## Installation

### From Source
```bash
git clone https://github.com/yourusername/ferex_cli.git
cd ferex_cli
go build -o ferex
```

### Using Go Install
```bash
go install github.com/yourusername/ferex_cli@latest
```

## Quick Start

1. **Generate a configuration template:**
   ```bash
   ferex init > my-retirement.yaml
   ```

2. **Edit the configuration file** with your personal information:
   ```yaml
   personal:
     name: "Your Name"
     birth_date: "1967-03-15T00:00:00Z"
     current_age: 57
     retirement_system: "FERS"
   
   employment:
     hire_date: "1999-01-15T00:00:00Z"
     current_salary: 85000
     high_3_salary: 82000
     creditable_service:
       total_years: 25
   
   retirement:
     target_age: 62
     survivor_benefit: "full"
   
   tsp:
     current_balance: 500000
     traditional_balance: 400000
     roth_balance: 100000
     withdrawal_strategy: "life_expectancy"
     growth_rate: 0.07
   
   social_security:
     estimated_pia: 2800
     claiming_age: 67
   ```

3. **Validate your configuration:**
   ```bash
   ferex validate my-retirement.yaml
   ```

4. **Calculate your retirement projection:**
   ```bash
   ferex calc my-retirement.yaml
   ```

## Usage Examples

### Basic Calculation
```bash
# Display results in table format
ferex calc my-retirement.yaml

# Verbose output with detailed projections
ferex calc my-retirement.yaml --verbose
```

### Export to Excel/Sheets
```bash
# Generate CSV for spreadsheet analysis
ferex calc my-retirement.yaml --format csv --output analysis.csv

# JSON for data processing
ferex calc my-retirement.yaml --format json --output data.json
```

### Compare Retirement Ages
```bash
# Compare retiring at different ages
ferex compare my-retirement.yaml --ages 57,60,62

# Export comparison to CSV
ferex compare my-retirement.yaml --ages 57,60,62 --format csv --output comparison.csv
```

### Template Generation
```bash
# Basic FERS employee template
ferex init --template basic > fers-plan.yaml

# Advanced template with military service and part-time periods
ferex init --template advanced > advanced-plan.yaml

# CSRS employee template
ferex init --template csrs > csrs-plan.yaml
```

## Configuration Templates

Ferex provides three configuration templates:

### Basic Template
- Simple FERS employee
- Full-time service
- Standard retirement planning

### Advanced Template
- Complex scenarios (military service, part-time periods)
- Early retirement options
- Spouse Social Security benefits

### CSRS Template
- Civil Service Retirement System
- Long federal career
- Different calculation rules

## Key Calculations

### FERS Pension
- **Standard**: 1.0% × High-3 × Years of Service
- **Age 62+ with 20+ years**: 1.1% × High-3 × Years of Service
- **Early retirement reduction**: 5% per year under age 62 (MRA+10)

### Survivor Benefits
- **Full survivor benefit**: 10% reduction to pension (FERS)
- **Partial survivor benefit**: 5% reduction to pension (FERS)
- **CSRS**: Complex calculation (2.5% of first $3,600 + 10% of remainder)

### Social Security
- Uses your estimated PIA from SSA statements
- Applies claiming age adjustments
- Calculates taxable portion based on provisional income

### TSP Withdrawals
- **Life expectancy**: Uses IRS Uniform Lifetime Table
- **Fixed amount**: User-specified annual withdrawal
- **Lump sum**: One-time withdrawal at retirement

## Output Formats

### Table Format
```
Retirement Planning Summary
===========================

Monthly Pension:           $1,691.25
Annual Pension:            $20,295.00
Social Security:           $2,800.00/month (starting age 67)
TSP Starting Balance:      $500,000.00
First Year Income:         $28,574.84
Lifetime Income:           $2,873,524.80
Replacement Ratio:         33.6%
```

### CSV Format
Perfect for Excel analysis with columns:
- Year, Age, Pension Income, Social Security, TSP Withdrawal
- Gross Income, Federal Tax, State Tax, Net Income
- TSP Balance progression

### JSON Format
Structured data for integration with other tools and APIs.

## Validation and Error Handling

Ferex includes comprehensive validation:

- **Retirement eligibility** checking
- **TSP balance consistency** verification
- **Date logic** validation
- **Interactive fixes** for common issues

```bash
# Validate with interactive fixes
ferex validate my-retirement.yaml --fix-interactive
```

## Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./pkg/calc -v
```

## Architecture

```
ferex_cli/
├── cmd/                    # Command-line interface
├── pkg/
│   ├── calc/              # Core calculation engine
│   ├── config/            # Configuration management
│   └── output/            # Output formatting
├── internal/
│   └── models/            # Data structures
└── main.go                # Application entry point
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## Accuracy and Disclaimers

**Important**: This tool provides estimates for planning purposes only. For official retirement calculations, consult:

- [OPM Retirement Services](https://www.opm.gov/retirement-center/)
- [Social Security Administration](https://www.ssa.gov/)
- [TSP](https://www.tsp.gov/)

The tool implements federal retirement rules as of 2025, including the repeal of WEP/GPO provisions.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- Report issues: [GitHub Issues](https://github.com/yourusername/ferex_cli/issues)
- Documentation: This README and inline help (`ferex --help`)
- Federal retirement resources: [OPM.gov](https://www.opm.gov/)