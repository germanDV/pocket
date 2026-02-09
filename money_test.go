package pocket

import (
	"math"
	"testing"
)

func TestMoney_StringAndFormat(t *testing.T) {
	tests := []struct {
		name       string
		amount     int64
		currency   string
		precision  int
		wantString string
		wantFormat string
	}{
		{
			name:       "zero",
			amount:     0,
			currency:   "USD",
			precision:  2,
			wantString: "0.00",
			wantFormat: "0.00 USD",
		},
		{
			name:       "positive",
			amount:     10099,
			currency:   "USD",
			precision:  2,
			wantString: "100.99",
			wantFormat: "100.99 USD",
		},
		{
			name:       "negative",
			amount:     -10099,
			currency:   "USD",
			precision:  2,
			wantString: "-100.99",
			wantFormat: "-100.99 USD",
		},
		{
			name:       "negative precision",
			amount:     -10099,
			currency:   "USD",
			precision:  8,
			wantString: "-10099.00000000",
			wantFormat: "-10099.00000000 USD",
		},
		{
			name:       "0 precision",
			amount:     10099,
			currency:   "JPY",
			precision:  0,
			wantString: "10099",
			wantFormat: "10099 JPY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMoney(tt.amount, tt.currency, tt.precision)
			AssertNil(t, err)
			AssertEqual(t, tt.wantString, m.String())
			AssertEqual(t, tt.wantFormat, m.Format())
		})
	}
}

func TestNewUSD(t *testing.T) {
	m := NewUSD(10050)
	AssertEqual(t, m.Amount(), int64(10050))
	AssertEqual(t, m.Currency(), "USD")
	AssertEqual(t, m.Precision(), 2)
	AssertEqual(t, m.String(), "100.50")
}

func TestNewARS(t *testing.T) {
	m := NewARS(99900)
	AssertEqual(t, m.Amount(), int64(99900))
	AssertEqual(t, m.Currency(), "ARS")
	AssertEqual(t, m.Precision(), 2)
	AssertEqual(t, m.String(), "999.00")
}

func TestNewMoney_Validation(t *testing.T) {
	tests := []struct {
		name      string
		amount    int64
		currency  string
		precision int
		wantError bool
	}{
		{
			name:      "valid with precision 0",
			amount:    100,
			currency:  "JPY",
			precision: 0,
			wantError: false,
		},
		{
			name:      "valid with precision 2",
			amount:    10099,
			currency:  "USD",
			precision: 2,
			wantError: false,
		},
		{
			name:      "valid with max precision 8",
			amount:    100,
			currency:  "BTC",
			precision: 8,
			wantError: false,
		},
		{
			name:      "invalid negative precision",
			amount:    100,
			currency:  "USD",
			precision: -1,
			wantError: true,
		},
		{
			name:      "invalid precision too high",
			amount:    100,
			currency:  "USD",
			precision: 9,
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMoney(tt.amount, tt.currency, tt.precision)
			if tt.wantError {
				AssertNotNil(t, err)
			} else {
				AssertNil(t, err)
				AssertEqual(t, m.Amount(), tt.amount)
				AssertEqual(t, m.Currency(), tt.currency)
				AssertEqual(t, m.Precision(), tt.precision)
			}
		})
	}
}

func TestMoney_Getters(t *testing.T) {
	m, err := NewMoney(50050, "EUR", 2)
	AssertNil(t, err)
	AssertEqual(t, m.Amount(), int64(50050))
	AssertEqual(t, m.Currency(), "EUR")
	AssertEqual(t, m.Precision(), 2)
}

func TestMoney_Plus(t *testing.T) {
	tests := []struct {
		name      string
		m1        Money
		m2        Money
		want      Money
		wantError bool
	}{
		{
			name: "add same currency",
			m1:   NewUSD(10000),
			m2:   NewUSD(5000),
			want: NewUSD(15000),
		},
		{
			name: "add with negative",
			m1:   NewUSD(10000),
			m2:   NewUSD(-3000),
			want: NewUSD(7000),
		},
		{
			name: "add zero",
			m1:   NewUSD(10000),
			m2:   NewUSD(0),
			want: NewUSD(10000),
		},
		{
			name:      "add different currencies",
			m1:        NewUSD(10000),
			m2:        NewARS(5000),
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.m1.Plus(tt.m2)
			if tt.wantError {
				AssertNotNil(t, err)
			} else {
				AssertNil(t, err)
				AssertTrue(t, result.Equals(tt.want))
			}
		})
	}
}

func TestMoney_Minus(t *testing.T) {
	tests := []struct {
		name      string
		m1        Money
		m2        Money
		want      Money
		wantError bool
	}{
		{
			name: "subtract same currency",
			m1:   NewUSD(10000),
			m2:   NewUSD(3000),
			want: NewUSD(7000),
		},
		{
			name: "subtract larger amount",
			m1:   NewUSD(10000),
			m2:   NewUSD(15000),
			want: NewUSD(-5000),
		},
		{
			name: "subtract zero",
			m1:   NewUSD(10000),
			m2:   NewUSD(0),
			want: NewUSD(10000),
		},
		{
			name:      "subtract different currencies",
			m1:        NewUSD(10000),
			m2:        NewARS(3000),
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.m1.Minus(tt.m2)
			if tt.wantError {
				AssertNotNil(t, err)
			} else {
				AssertNil(t, err)
				AssertTrue(t, result.Equals(tt.want))
			}
		})
	}
}

func TestMoney_Inc(t *testing.T) {
	tests := []struct {
		name   string
		m      Money
		amount int64
		want   Money
	}{
		{
			name:   "increment positive",
			m:      NewUSD(10000),
			amount: 5000,
			want:   NewUSD(15000),
		},
		{
			name:   "increment negative",
			m:      NewUSD(10000),
			amount: -2000,
			want:   NewUSD(8000),
		},
		{
			name:   "increment zero",
			m:      NewUSD(10000),
			amount: 0,
			want:   NewUSD(10000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.m.Inc(tt.amount)
			AssertNil(t, err)
			AssertTrue(t, result.Equals(tt.want))
		})
	}
}

func TestMoney_Dec(t *testing.T) {
	tests := []struct {
		name   string
		m      Money
		amount int64
		want   Money
	}{
		{
			name:   "decrement positive",
			m:      NewUSD(10000),
			amount: 3000,
			want:   NewUSD(7000),
		},
		{
			name:   "decrement negative",
			m:      NewUSD(10000),
			amount: -2000,
			want:   NewUSD(12000),
		},
		{
			name:   "decrement zero",
			m:      NewUSD(10000),
			amount: 0,
			want:   NewUSD(10000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.m.Dec(tt.amount)
			AssertNil(t, err)
			AssertTrue(t, result.Equals(tt.want))
		})
	}
}

func TestMoney_Times(t *testing.T) {
	tests := []struct {
		name   string
		m      Money
		factor int64
		want   Money
	}{
		{
			name:   "multiply by positive",
			m:      NewUSD(1000),
			factor: 5,
			want:   NewUSD(5000),
		},
		{
			name:   "multiply by negative",
			m:      NewUSD(1000),
			factor: -3,
			want:   NewUSD(-3000),
		},
		{
			name:   "multiply by zero",
			m:      NewUSD(1000),
			factor: 0,
			want:   NewUSD(0),
		},
		{
			name:   "multiply by one",
			m:      NewUSD(1000),
			factor: 1,
			want:   NewUSD(1000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.m.Times(tt.factor)
			AssertNil(t, err)
			AssertTrue(t, result.Equals(tt.want))
		})
	}
}

func TestMoney_DividedBy(t *testing.T) {
	tests := []struct {
		name    string
		m       Money
		divisor int64
		want    string
		wantErr bool
	}{
		{
			name:    "dividing by zero should return error",
			m:       NewUSD(10_00),
			divisor: 0,
			wantErr: true,
			want:    "",
		},
		{
			name:    "round down",
			m:       NewUSD(100_00),
			divisor: 3,
			wantErr: false,
			want:    "33.33 USD",
		},
		{
			name:    "round up",
			m:       NewUSD(200_00),
			divisor: 3,
			wantErr: false,
			want:    "66.67 USD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.m.DividedBy(tt.divisor)
			if tt.wantErr {
				AssertNotNil(t, err)
			} else {
				AssertNil(t, err)
				AssertEqual(t, result.Format(), tt.want)
			}
		})
	}
}

func TestMoney_Equals(t *testing.T) {
	tests := []struct {
		name string
		m1   Money
		m2   Money
		want bool
	}{
		{
			name: "same amount and currency",
			m1:   NewUSD(10000),
			m2:   NewUSD(10000),
			want: true,
		},
		{
			name: "different amount",
			m1:   NewUSD(10000),
			m2:   NewUSD(5000),
			want: false,
		},
		{
			name: "different currency",
			m1:   NewUSD(10000),
			m2:   NewARS(10000),
			want: false,
		},
		{
			name: "different precision",
			m1:   Must(NewMoney(10000, "USD", 2)),
			m2:   Must(NewMoney(1000000, "USD", 4)),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m1.Equals(tt.m2)
			AssertEqual(t, got, tt.want)
		})
	}
}

func TestMoney_Uninitialized(t *testing.T) {
	var uninitialized Money

	t.Run("String returns empty", func(t *testing.T) {
		AssertEqual(t, uninitialized.String(), "")
	})

	t.Run("Format returns empty", func(t *testing.T) {
		AssertEqual(t, uninitialized.Format(), "")
	})

	t.Run("Plus returns error", func(t *testing.T) {
		_, err := uninitialized.Plus(NewUSD(100))
		AssertNotNil(t, err)
	})

	t.Run("Minus returns error", func(t *testing.T) {
		_, err := uninitialized.Minus(NewUSD(100))
		AssertNotNil(t, err)
	})

	t.Run("Inc returns error", func(t *testing.T) {
		_, err := uninitialized.Inc(100)
		AssertNotNil(t, err)
	})

	t.Run("Dec returns error", func(t *testing.T) {
		_, err := uninitialized.Dec(100)
		AssertNotNil(t, err)
	})

	t.Run("Times returns error", func(t *testing.T) {
		_, err := uninitialized.Times(2)
		AssertNotNil(t, err)
	})
}

func TestMoney_Overflow(t *testing.T) {
	t.Run("Plus overflow", func(t *testing.T) {
		m1, _ := NewMoney(math.MaxInt64-100, "USD", 2)
		m2 := NewUSD(1000)
		_, err := m1.Plus(m2)
		AssertNotNil(t, err)
	})

	t.Run("Minus underflow", func(t *testing.T) {
		m1, _ := NewMoney(math.MinInt64+100, "USD", 2)
		m2 := NewUSD(1000)
		_, err := m1.Minus(m2)
		AssertNotNil(t, err)
	})

	t.Run("Times overflow", func(t *testing.T) {
		m, _ := NewMoney(math.MaxInt64/2+1000, "USD", 2)
		_, err := m.Times(3)
		AssertNotNil(t, err)
	})
}

func Must(m Money, err error) Money {
	if err != nil {
		panic(err)
	}
	return m
}

func TestNewMoneyFromString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      Money
		wantError bool
	}{
		{
			name:  "valid USD with 2 decimals",
			input: "100.99 USD",
			want:  NewUSD(10099),
		},
		{
			name:  "valid ARS with 2 decimals",
			input: "100.00 ARS",
			want:  NewARS(10000),
		},
		{
			name:  "valid BTC with 8 decimals",
			input: "1.00000000 BTC",
			want:  Must(NewMoney(100000000, "BTC", 8)),
		},
		{
			name:  "lowercase currency is uppercased",
			input: "100.99 usd",
			want:  NewUSD(10099),
		},
		{
			name:  "negative amount",
			input: "-100.50 USD",
			want:  NewUSD(-10050),
		},
		{
			name:  "zero amount",
			input: "0.00 USD",
			want:  NewUSD(0),
		},
		{
			name:  "single decimal digit",
			input: "100.5 USD",
			want:  Must(NewMoney(1005, "USD", 1)),
		},
		{
			name:  "8 decimal precision",
			input: "0.12345678 BTC",
			want:  Must(NewMoney(12345678, "BTC", 8)),
		},
		{
			name:      "missing space separator",
			input:     "100.99USD",
			wantError: true,
		},
		{
			name:      "too many parts",
			input:     "100.99 USD extra",
			wantError: true,
		},
		{
			name:      "missing decimal point",
			input:     "100 USD",
			wantError: true,
		},
		{
			name:      "too many decimal points",
			input:     "100.99.99 USD",
			wantError: true,
		},
		{
			name:      "precision too high",
			input:     "1.000000000 BTC",
			wantError: true,
		},
		{
			name:      "invalid integer part",
			input:     "abc.99 USD",
			wantError: true,
		},
		{
			name:      "invalid fractional part",
			input:     "100.xyz USD",
			wantError: true,
		},
		{
			name:      "empty string",
			input:     "",
			wantError: true,
		},
		{
			name:      "only currency",
			input:     "USD",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMoneyFromString(tt.input)
			if tt.wantError {
				AssertNotNil(t, err)
			} else {
				AssertNil(t, err)
				AssertTrue(t, got.Equals(tt.want))
			}
		})
	}
}
