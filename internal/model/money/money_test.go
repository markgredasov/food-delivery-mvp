package money_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/markgredasov/food-delivery-mvp/internal/model/money"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		amount  decimal.Decimal
		want    string
		wantErr bool
	}{
		{
			name:    "valid positive amount",
			amount:  decimal.NewFromFloat(350.50),
			want:    "350.50",
			wantErr: false,
		},
		{
			name:    "valid zero amount",
			amount:  decimal.Zero,
			want:    "0.00",
			wantErr: false,
		},
		{
			name:    "rounds to 2 decimal places",
			amount:  decimal.NewFromFloat(10.555),
			want:    "10.56",
			wantErr: false,
		},
		{
			name:    "negative amount",
			amount:  decimal.NewFromFloat(-5.00),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := money.New(tt.amount)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "must not be negative")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, m.String())
			}
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid with two decimals",
			input:   "350.50",
			want:    "350.50",
			wantErr: false,
		},
		{
			name:    "valid with one decimal",
			input:   "350.5",
			want:    "350.50",
			wantErr: false,
		},
		{
			name:    "valid integer",
			input:   "350",
			want:    "350.00",
			wantErr: false,
		},
		{
			name:    "valid zero",
			input:   "0.00",
			want:    "0.00",
			wantErr: false,
		},
		{
			name:    "negative amount",
			input:   "-5.00",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "abc",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := money.FromString(tt.input)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, m.String())
			}
		})
	}
}

func TestFromInt(t *testing.T) {
	tests := []struct {
		name string
		v    int64
		want string
	}{
		{"positive", 350, "350.00"},
		{"zero", 0, "0.00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.FromInt(tt.v)
			assert.Equal(t, tt.want, m.String())
		})
	}
}

func TestMoney_Add(t *testing.T) {
	m1 := money.MustFromString("10.50")
	m2 := money.MustFromString("5.25")
	m3 := money.MustFromString("0.00")

	result := m1.Add(m2)
	assert.Equal(t, "15.75", result.String())

	resultZero := m1.Add(m3)
	assert.Equal(t, "10.50", resultZero.String())
}

func TestMoney_Mul(t *testing.T) {
	m := money.MustFromString("10.50")

	tests := []struct {
		qty  int
		want string
	}{
		{1, "10.50"},
		{2, "21.00"},
		{3, "31.50"},
		{0, "0.00"},
		{-1, "-10.50"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := m.Mul(tt.qty)
			assert.Equal(t, tt.want, result.String())
		})
	}
}

func TestMoney_Equal(t *testing.T) {
	m1 := money.MustFromString("10.50")
	m2 := money.MustFromString("10.50")
	m3 := money.MustFromString("10.51")
	m4 := money.MustFromString("10.5")

	assert.True(t, m1.Equal(m2))
	assert.True(t, m1.Equal(m4))
	assert.False(t, m1.Equal(m3))
}

func TestMoney_IsZero(t *testing.T) {
	m1 := money.MustFromString("0.00")
	m2 := money.MustFromString("10.50")
	m3 := money.MustFromString("0")

	assert.True(t, m1.IsZero())
	assert.False(t, m2.IsZero())
	assert.True(t, m3.IsZero())
}

func TestMoney_MarshalJSON(t *testing.T) {
	m := money.MustFromString("350.50")

	data, err := json.Marshal(m)
	require.NoError(t, err)
	assert.Equal(t, `"350.50"`, string(data))

	zero := money.Zero
	data, err = json.Marshal(zero)
	require.NoError(t, err)
	assert.Equal(t, `"0.00"`, string(data))
}

func TestMoney_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "string amount",
			input:   `"350.50"`,
			want:    "350.50",
			wantErr: false,
		},
		{
			name:    "integer number",
			input:   `350`,
			want:    "350.00",
			wantErr: false,
		},
		{
			name:    "float number",
			input:   `350.5`,
			want:    "350.50",
			wantErr: false,
		},
		{
			name:    "negative amount",
			input:   `"-5.00"`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   `invalid`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m money.Money
			err := json.Unmarshal([]byte(tt.input), &m)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, m.String())
			}
		})
	}
}

func TestMoney_Value(t *testing.T) {
	m := money.MustFromString("350.50")

	val, err := m.Value()
	require.NoError(t, err)
	assert.Equal(t, "350.50", val)
	assert.IsType(t, driver.Value(""), val)
}

func TestMoney_Scan(t *testing.T) {
	tests := []struct {
		name    string
		src     any
		want    string
		wantErr bool
	}{
		{
			name:    "valid decimal as string",
			src:     "350.50",
			want:    "350.50",
			wantErr: false,
		},
		{
			name:    "valid decimal as []byte",
			src:     []byte("350.50"),
			want:    "350.50",
			wantErr: false,
		},
		{
			name:    "valid decimal as int64",
			src:     int64(350),
			want:    "350.00",
			wantErr: false,
		},
		{
			name:    "negative amount",
			src:     "-5.00",
			want:    "",
			wantErr: true,
		},
		{
			name:    "nil source",
			src:     nil,
			want:    "0.00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m money.Money
			err := m.Scan(tt.src)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, m.String())
			}
		})
	}
}

func TestMoney_ScanValueRoundTrip(t *testing.T) {
	original := money.MustFromString("350.50")

	val, err := original.Value()
	require.NoError(t, err)

	var scanned money.Money
	err = scanned.Scan(val)
	require.NoError(t, err)

	assert.True(t, original.Equal(scanned))
}
