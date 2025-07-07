# Ferex CLI Usage Guide

## Command Reference

### Global Flags
- `--config string`: Config file (default: $HOME/.ferex.yaml)
- `--format string`: Output format (table, json, csv, yaml) (default: "table")
- `--verbose`: Verbose output
- `--help`: Show help

### Commands

#### `ferex init`
Generate configuration templates.

**Flags:**
- `--template string`: Template type (basic, advanced, csrs) (default: "basic")

**Examples:**
```bash
# Generate basic FERS template
ferex init > my-plan.yaml

# Generate advanced template with all options
ferex init --template advanced > advanced-plan.yaml

# Generate CSRS template
ferex init --template csrs > csrs-plan.yaml
```

#### `ferex validate`
Validate configuration files.

**Usage:** `ferex validate [config-file]`

**Flags:**
- `--fix-interactive`: Interactively fix validation issues

**Examples:**
```bash
# Basic validation
ferex validate my-plan.yaml

# Interactive validation with fixes
ferex validate my-plan.yaml --fix-interactive
```

#### `ferex calc`
Calculate retirement projections.

**Usage:** `ferex calc [config-file]`

**Flags:**
- `--output string`: Output file (default: stdout)

**Examples:**
```bash
# Display results in terminal
ferex calc my-plan.yaml

# Save to CSV file
ferex calc my-plan.yaml --format csv --output results.csv

# Verbose output with detailed projections
ferex calc my-plan.yaml --verbose

# JSON output for data processing
ferex calc my-plan.yaml --format json --output data.json
```

#### `ferex compare`
Compare different retirement scenarios.

**Usage:** `ferex compare [config-file]`

**Flags:**
- `--ages stringSlice`: Retirement ages to compare (default: [57,62])
- `--output string`: Output file (default: stdout)

**Examples:**
```bash
# Compare default ages (57 vs 62)
ferex compare my-plan.yaml

# Compare three ages
ferex compare my-plan.yaml --ages 57,60,62

# Export comparison to CSV
ferex compare my-plan.yaml --ages 55,57,60,62 --format csv --output comparison.csv
```

## Configuration File Structure

### Required Sections

#### Personal Information
```yaml
personal:
  name: "Your Name"                    # Full name
  birth_date: "1967-03-15T00:00:00Z"  # ISO 8601 format
  current_age: 57                      # Current age (auto-calculated if omitted)
  retirement_system: "FERS"            # "FERS" or "CSRS"
```

#### Employment Information
```yaml
employment:
  hire_date: "1999-01-15T00:00:00Z"   # Federal service start date
  current_salary: 85000               # Current annual salary
  high_3_salary: 82000               # High-3 average (auto-calculated if omitted)
  creditable_service:
    total_years: 25                   # Total creditable service years
    part_time_periods: []             # Part-time service periods (optional)
    military_service:                 # Military service (optional)
      years: 4
      bought_back: true
    unused_sick_leave: 0             # Hours of unused sick leave (optional)
```

#### Retirement Planning
```yaml
retirement:
  target_age: 62                      # Planned retirement age
  survivor_benefit: "full"            # "full", "partial", or "none"
  early_retirement:                   # Early retirement options (optional)
    type: "MRA+10"                   # "MRA+10", "VERA", or "DSR"
    postponed_start: false           # Postpone annuity start
```

#### TSP Information
```yaml
tsp:
  current_balance: 500000             # Total TSP balance
  traditional_balance: 400000         # Traditional TSP balance
  roth_balance: 100000               # Roth TSP balance
  withdrawal_strategy: "life_expectancy"  # "fixed_amount", "life_expectancy", "lump_sum"
  withdrawal_amount: 0               # For fixed_amount strategy
  growth_rate: 0.07                  # Annual growth rate assumption
```

#### Social Security
```yaml
social_security:
  estimated_pia: 2800                # Primary Insurance Amount from SSA
  claiming_age: 67                   # Age when you'll claim SS benefits
  spouse_benefit:                    # Spouse information (optional)
    estimated_pia: 2200
    claiming_age: 67
```

### Optional Sections

#### Output Preferences
```yaml
output:
  format: "table"                    # "table", "csv", "json", "yaml"
  verbose: false                     # Include detailed projections
  output_file: ""                    # File to save results
```

#### Part-Time Service Periods
```yaml
employment:
  creditable_service:
    part_time_periods:
      - start_date: "2010-01-01T00:00:00Z"
        end_date: "2012-12-31T00:00:00Z"
        hours_per_week: 32
```

## Understanding Output

### Summary Section
- **Monthly/Annual Pension**: Your FERS/CSRS annuity
- **Pension Reduction**: Early retirement reduction percentage
- **Survivor Benefit Cost**: Annual cost of survivor benefit election
- **FERS Supplement**: Monthly supplement until age 62 (if eligible)
- **Social Security**: Monthly benefit at your claiming age
- **TSP Depletion Age**: When TSP balance reaches zero (if applicable)
- **Replacement Ratio**: Retirement income as percentage of current salary

### Annual Projections (CSV Export)
Each row represents one year of retirement with:
- **Income Sources**: Pension, FERS Supplement, Social Security, TSP withdrawals
- **Taxes**: Federal and state income tax estimates
- **Net Income**: Take-home pay after taxes and deductions
- **TSP Balance**: Account balance progression

### Comparison Analysis
Side-by-side comparison of different retirement ages showing:
- Financial impact of retiring earlier vs. later
- Lifetime income differences
- Replacement ratio variations

## Common Scenarios

### Scenario 1: Standard FERS Retirement
```yaml
# Age 62 with 25+ years service
retirement:
  target_age: 62
  survivor_benefit: "full"
```
- 1.1% multiplier (if 20+ years)
- No early retirement reduction
- Eligible for immediate pension

### Scenario 2: Early FERS Retirement (MRA+10)
```yaml
# MRA with 10-29 years service
retirement:
  target_age: 57  # Assuming 1967 birth year (MRA = 57)
  survivor_benefit: "partial"
  early_retirement:
    type: "MRA+10"
    postponed_start: false
```
- 5% reduction per year under 62
- Can postpone to reduce/eliminate penalty
- Not eligible for FERS Supplement

### Scenario 3: FERS with Military Service
```yaml
employment:
  creditable_service:
    total_years: 28
    military_service:
      years: 6
      bought_back: true
```
- Military time counts toward pension calculation
- Must pay deposit for post-1956 military service
- Affects eligibility and computation

### Scenario 4: CSRS Employee
```yaml
personal:
  retirement_system: "CSRS"
employment:
  creditable_service:
    total_years: 42
```
- Different calculation formula (tiered percentages)
- Different survivor benefit costs
- No TSP matching contributions

## Tips and Best Practices

### Accurate Data Entry
1. **Use SSA estimates**: Get your PIA from [ssa.gov/myaccount](https://ssa.gov/myaccount)
2. **Verify High-3**: Check with HR or payroll for accurate salary history
3. **Confirm service time**: Include all creditable service (military, transfers, etc.)
4. **TSP balances**: Use current statement values

### Planning Considerations
1. **Compare multiple ages**: Use the compare command to see financial impact
2. **Consider survivor needs**: Model surviving spouse income requirements
3. **Tax planning**: Review state tax implications for retirement location
4. **Healthcare costs**: Factor in FEHB premiums and Medicare decisions

### Excel Analysis
Export CSV data for advanced analysis:
```bash
ferex calc my-plan.yaml --format csv --output retirement-data.csv
```

Then create charts in Excel/Sheets for:
- Income progression over time
- TSP balance depletion
- Tax burden analysis
- Net worth projections

## Troubleshooting

### Common Validation Errors
- **"FERS eligibility not met"**: Check age and service requirements
- **"TSP balance inconsistency"**: Ensure traditional + roth = current balance
- **"Birth date must be before hire date"**: Verify date formats

### Getting Help
```bash
# Command-specific help
ferex calc --help
ferex compare --help

# Validate configuration
ferex validate my-plan.yaml

# Check with verbose output
ferex calc my-plan.yaml --verbose
```

### Known Limitations
- State tax calculations are simplified estimates
- COLA assumptions use average historical rates
- Some complex scenarios (VERA/DSR) are simplified
- Federal tax calculations use current year brackets