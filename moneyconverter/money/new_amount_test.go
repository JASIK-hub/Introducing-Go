package money

import (
	"errors"
	"testing"
)

func TestNewAmount(t *testing.T) {
	tt := map[string]struct {
		quantity Decimal
		currency Currency
		expected Amount
		err      error
	}{
		"quantity precision too presice": {
			quantity: Decimal{subunits: 192345, precision: 4},
			currency: Currency{code: "USD", precision: 2},
			expected: Amount{},
			err:      ErrTooPrecise,
		}, "precision equals currency precision": {
			quantity: Decimal{subunits: 1234, precision: 2},
			currency: Currency{code: "USD", precision: 2},
			expected: Amount{
				quantity: Decimal{subunits: 1234, precision: 2},
				currency: Currency{code: "USD", precision: 2},
			},
			err: nil,
		},
		"zero quantity": {
			quantity: Decimal{subunits: 0, precision: 2},
			currency: Currency{code: "USD", precision: 2},
			expected: Amount{
				quantity: Decimal{subunits: 0, precision: 2},
				currency: Currency{code: "USD", precision: 2},
			},
			err: nil,
		},
		"quantity precision is lower": {
			quantity: Decimal{subunits: 12, precision: 0},
			currency: Currency{code: "USD", precision: 2},
			expected: Amount{
				quantity: Decimal{subunits: 12, precision: 2},
				currency: Currency{code: "USD", precision: 2},
			},
			err: nil,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := NewAmount(tc.quantity, tc.currency)
			if !errors.Is(err, tc.err) {
				t.Errorf("expected %v, got %v", tc.err, err)
			}
			if got != tc.expected {
				t.Errorf("expected amount %+v, got %+v", tc.expected, got)
			}
		})
	}
}
