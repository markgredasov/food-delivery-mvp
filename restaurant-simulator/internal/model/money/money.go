// Package money provides a small decimal-backed Money value object used
// everywhere a monetary amount crosses a domain boundary.
//
// Amounts are represented with shopspring/decimal to avoid the rounding
// errors floating point would introduce when summing item prices. The MVP
// assumes a single implicit currency (RUB) — no currency field is carried,
// which is documented as a deliberate simplification in the README.
package money

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

// Money is a non-negative monetary amount with 2 decimal places.
type Money struct {
	amount decimal.Decimal
}

// Zero is the additive identity.
var Zero = Money{amount: decimal.Zero}

// New builds a Money from a decimal, rejecting negative amounts.
func New(d decimal.Decimal) (Money, error) {
	if d.IsNegative() {
		return Money{}, fmt.Errorf("money: amount must not be negative, got %s", d.String())
	}
	return Money{amount: d.Round(2)}, nil
}

// FromString parses a decimal string (e.g. "350.00") into Money.
func FromString(s string) (Money, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Money{}, fmt.Errorf("money: invalid amount %q: %w", s, err)
	}
	return New(d)
}

// FromInt builds Money from an integer amount of the major unit (e.g. rubles).
func FromInt(v int64) Money {
	m, _ := New(decimal.NewFromInt(v))
	return m
}

// MustFromString is like FromString but panics on error; intended for
// constants/tests, never for user input.
func MustFromString(s string) Money {
	m, err := FromString(s)
	if err != nil {
		panic(err)
	}
	return m
}

// Decimal returns the underlying decimal value.
func (m *Money) Decimal() decimal.Decimal { return m.amount }

// String renders the amount as a fixed 2-decimal string.
func (m *Money) String() string { return m.amount.StringFixed(2) }

// Add returns the sum of m and other.
func (m *Money) Add(other Money) Money {
	return Money{amount: m.amount.Add(other.amount)}
}

// Mul multiplies the amount by a non-negative integer quantity.
func (m *Money) Mul(qty int) Money {
	return Money{amount: m.amount.Mul(decimal.NewFromInt(int64(qty)))}
}

// Equal reports whether m and other represent the same amount.
func (m *Money) Equal(other Money) bool { return m.amount.Equal(other.amount) }

// IsZero reports whether the amount is zero.
func (m *Money) IsZero() bool { return m.amount.IsZero() }

// MarshalJSON renders the amount as a JSON string ("350.00"), matching the
// OpenAPI schema which represents money as a decimal string to avoid float
// precision issues on the wire.
func (m *Money) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.amount.StringFixed(2) + `"`), nil
}

// UnmarshalJSON accepts either a JSON string ("350.00") or a JSON number.
func (m *Money) UnmarshalJSON(data []byte) error {
	var d decimal.Decimal
	if err := json.Unmarshal(data, &d); err != nil {
		return fmt.Errorf("money: %w", err)
	}
	v, err := New(d)
	if err != nil {
		return err
	}
	*m = v
	return nil
}

// Value implements [driver.Valuer] for pgx/database-sql NUMERIC columns.
func (m Money) Value() (driver.Value, error) {
	return m.amount.StringFixed(2), nil
}

// Scan implements sql.Scanner for pgx/database-sql NUMERIC columns.
func (m *Money) Scan(src any) error {
	d := decimal.Decimal{}
	if err := d.Scan(src); err != nil {
		return err
	}
	v, err := New(d)
	if err != nil {
		return err
	}
	*m = v
	return nil
}
