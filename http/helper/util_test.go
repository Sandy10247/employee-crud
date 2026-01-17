// helper/helper_test.go
package helper

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateNetSalary(t *testing.T) {
	tests := []struct {
		name            string
		grossSalary     float64
		totalDeductions float64
		want            float64
	}{
		{
			name:            "normal case",
			grossSalary:     100000,
			totalDeductions: 28000,
			want:            72000,
		},
		{
			name:            "deductions equal to gross",
			grossSalary:     50000,
			totalDeductions: 50000,
			want:            0,
		},
		{
			name:            "deductions exceed gross → should return 0",
			grossSalary:     45000,
			totalDeductions: 52000,
			want:            0,
		},
		{
			name:            "zero gross",
			grossSalary:     0,
			totalDeductions: 10000,
			want:            0,
		},
		{
			name:            "negative deductions (should still subtract)",
			grossSalary:     80000,
			totalDeductions: -5000,
			want:            85000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateNetSalary(tt.grossSalary, tt.totalDeductions)
			assert.Equal(t, tt.want, got, "net salary mismatch")
		})
	}
}

func TestCalculatePercentage(t *testing.T) {
	tests := []struct {
		name    string
		percent float64
		total   float64
		want    float64
	}{
		{
			name:    "25% of 400",
			percent: 25,
			total:   400,
			want:    100,
		},
		{
			name:    "0% of anything",
			percent: 0,
			total:   999999,
			want:    0,
		},
		{
			name:    "100% ",
			percent: 100,
			total:   7500,
			want:    7500,
		},
		{
			name:    "12.5% of 800",
			percent: 12.5,
			total:   800,
			want:    100,
		},
		{
			name:    "negative percent",
			percent: -20,
			total:   500,
			want:    -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePercentage(tt.percent, tt.total)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloatToNumeric(t *testing.T) {
	tests := []struct {
		name      string
		input     float64
		precision int
		wantStr   string // expected string representation
		wantValid bool
	}{
		{
			name:      "simple integer",
			input:     45000.0,
			precision: 2,
			wantStr:   "45000.00",
			wantValid: true,
		},
		{
			name:      "two decimal places",
			input:     12345.6789,
			precision: 2,
			wantStr:   "12345.68",
			wantValid: true,
		},
		{
			name:      "zero precision → rounds to integer",
			input:     9876.54321,
			precision: 0,
			wantStr:   "9877",
			wantValid: true,
		},
		{
			name:      "very small number",
			input:     0.00034567,
			precision: 8,
			wantStr:   "0.00034567",
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FloatToNumeric(tt.input, tt.precision)

			if tt.wantValid {
				require.NoError(t, err)
				require.True(t, got.Valid)

				// Compare string representation
				val, err := got.Value()
				if err != nil {
					t.Fatalf("failed to get value from sql.NullString: %v", err)
				}
				gotStr := val.(string)
				assert.Equal(t, tt.wantStr, gotStr, "string representation mismatch")
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	const testKey = "TEST_HELPER_ENV_VAR"

	tests := []struct {
		name         string
		setEnv       bool
		envValue     string
		defaultValue string
		want         string
	}{
		{
			name:         "env exists",
			setEnv:       true,
			envValue:     "mysql://user:pass@localhost:3306",
			defaultValue: "default-dsn",
			want:         "mysql://user:pass@localhost:3306",
		},
		{
			name:         "env not set → use default",
			setEnv:       false,
			envValue:     "",
			defaultValue: "postgres://default",
			want:         "postgres://default",
		},
		{
			name:         "empty env value",
			setEnv:       true,
			envValue:     "",
			defaultValue: "fallback",
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(testKey, tt.envValue)
			}

			got := GetEnv(testKey, tt.defaultValue)
			assert.Equal(t, tt.want, got)
			t.Cleanup(func() {
				os.Unsetenv(testKey)
			})
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	const testKey = "TEST_INT_ENV"

	t.Cleanup(func() {
		os.Unsetenv(testKey)
	})

	tests := []struct {
		name         string
		setEnv       bool
		envValue     string
		defaultValue int
		want         int
	}{
		{
			name:         "valid integer",
			setEnv:       true,
			envValue:     "8080",
			defaultValue: 9000,
			want:         8080,
		},
		{
			name:         "invalid integer → fallback",
			setEnv:       true,
			envValue:     "not-a-number",
			defaultValue: 3000,
			want:         3000,
		},
		{
			name:         "env not set",
			setEnv:       false,
			envValue:     "",
			defaultValue: 5432,
			want:         5432,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(testKey, tt.envValue)
			}

			got := GetEnvInt(testKey, tt.defaultValue)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTaxRatePerCountry(t *testing.T) {
	const taxKeyIndia = "india"
	const taxKeySingapore = "singapore"

	t.Cleanup(func() {
		os.Unsetenv(taxKeyIndia)
		os.Unsetenv(taxKeySingapore)
	})

	tests := []struct {
		name     string
		country  string
		setEnv   bool
		envValue string
		want     float64
	}{
		{
			name:     "found in env - India",
			country:  "India",
			setEnv:   true,
			envValue: "0.30",
			want:     0.30,
		},
		{
			name:     "found in env - lowercase",
			country:  "singapore",
			setEnv:   true,
			envValue: "0.22",
			want:     0.22,
		},
		{
			name:     "not found → returns 0",
			country:  "Japan",
			setEnv:   false,
			envValue: "",
			want:     0.0,
		},
		{
			name:     "invalid float in env → returns 0",
			country:  "germany",
			setEnv:   true,
			envValue: "thirty percent",
			want:     0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(strings.ToLower(tt.country), tt.envValue)
			}

			got := GetTaxRatePerCountry(tt.country)
			assert.Equal(t, tt.want, got, "tax rate mismatch")
		})
	}
}
