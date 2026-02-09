package pocket

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Money represents a monetary value.
// A Money instance is immutable, operations return a new Money instance.
// Precision is limited to 8 digits to accommdate fairly large values without overflowing.
type Money struct {
	amount    int64
	currency  string
	precision int
	// Sentinel value to ensure Money instances are created with the constructor.
	initialized bool
}

// NewUSD creates a new Money instance with USD currency.
func NewUSD(amount int64) Money {
	return Money{
		amount:      amount,
		currency:    "USD",
		precision:   2,
		initialized: true,
	}
}

// NewARS creates a new Money instance with ARS currency.
func NewARS(amount int64) Money {
	return Money{
		amount:      amount,
		currency:    "ARS",
		precision:   2,
		initialized: true,
	}
}

// NewMoney creates a new Money instance.
func NewMoney(amount int64, currency string, precision int) (Money, error) {
	if precision < 0 {
		return Money{}, fmt.Errorf("precision must be non-negative")
	}
	if precision > 8 {
		return Money{}, fmt.Errorf("precision must be less than or equal to 8")
	}

	return Money{
		amount:      amount,
		currency:    currency,
		precision:   precision,
		initialized: true,
	}, nil
}

// NewMoneyFromString creates a new Money instance from a string. The string must be in the format "amount currency".
// The number of decimal places determines the precision. So be sure to include 0s if necessary.
// And be careful not to use unsanitized user input as "100 USD" will be different from "100.00 USD".
// e.g., "100.99 USD" // precision=2
// e.g., "100.00 ARS" // precision=2
// e.g., "100 ARS" // precision=0, not what you want for ARS (and most FIAT currencies)
// e.g., "1.00000000 BTC" // precision=8
func NewMoneyFromString(s string) (Money, error) {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return Money{}, fmt.Errorf("invalid string format: %s", s)
	}

	amount := parts[0]
	currency := strings.ToUpper(parts[1])

	amountParts := strings.Split(amount, ".")
	if len(amountParts) != 2 {
		return Money{}, fmt.Errorf("invalid amount format: %s - expected a '.'", amount)
	}

	precision := len(amountParts[1])
	if precision > 8 {
		return Money{}, fmt.Errorf("invalid amount format: %s - precision must be less than or equal to 8", amount)
	}

	amountInt, err := strconv.ParseInt(amountParts[0], 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("invalid amount format: %s - %w", amount, err)
	}
	amountFrac, err := strconv.ParseInt(amountParts[1], 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("invalid amount format: %s - %w", amount, err)
	}

	multiplier := int64(math.Pow10(precision))

	total := amountInt * multiplier
	if amountInt < 0 {
		total -= amountFrac
	} else {
		total += amountFrac
	}

	return NewMoney(total, currency, precision)
}

// Returns the currency of the money.
func (m Money) Currency() string {
	return m.currency
}

// Returns the precision of the money.
func (m Money) Precision() int {
	return m.precision
}

// Returns the amount of money in the smallest unit of the currency.
// For example, if money is `Money{amount: 10099, currency: "USD"}`, the amount will be 10099.
func (m Money) Amount() int64 {
	return m.amount
}

// String returns the amount in major units with proper decimal places.
// e.g., amount=10099, precision=2 → "100.99"
// e.g., amount=1000, precision=2 → "10.00"
// e.g., amount=-10099, precision=2 → "-100.99"
func (m Money) String() string {
	if !m.initialized {
		return ""
	}

	if m.precision == 0 {
		return fmt.Sprintf("%d", m.amount)
	}

	divisor := int64(1)
	for i := 0; i < m.precision; i++ {
		divisor *= 10
	}

	// If |amount| < divisor, the amount is already in major units
	// (e.g., -10099 with precision 8 means -10099.00000000, not -0.00010099)
	if m.amount < 0 && -m.amount < divisor {
		return fmt.Sprintf("%d.%0*d", m.amount, m.precision, 0)
	}
	if m.amount >= 0 && m.amount < divisor {
		return fmt.Sprintf("%d.%0*d", m.amount, m.precision, 0)
	}

	major := m.amount / divisor
	minor := m.amount % divisor

	if minor < 0 {
		minor = -minor
	}

	format := fmt.Sprintf("%%d.%%0%dd", m.precision)
	return fmt.Sprintf(format, major, minor)
}

// Format returns "amount currency" format.
// e.g., "100.99 USD"
func (m Money) Format() string {
	if !m.initialized {
		return ""
	}
	return fmt.Sprintf("%s %s", m.String(), m.currency)
}

// Plus returns a new Money with the sum of the two amounts.
// Returns an error if the currencies don't match or if overflow occurs.
func (m Money) Plus(other Money) (Money, error) {
	if !m.initialized || !other.initialized {
		return Money{}, errors.New("Money instances must be created with the constructor")
	}

	if m.currency != other.Currency() {
		return Money{}, fmt.Errorf("cannot add %s to %s: currencies must match", other.Currency(), m.currency)
	}

	sum, err := TrySafeAdd(m.amount, other.Amount())
	if err != nil {
		return Money{}, fmt.Errorf("cannot add amounts: %w", err)
	}

	return NewMoney(sum, m.currency, m.precision)
}

// Minus returns a new Money with the difference of the two amounts.
// Returns an error if the currencies don't match or if overflow occurs.
func (m Money) Minus(other Money) (Money, error) {
	if !m.initialized || !other.initialized {
		return Money{}, errors.New("Money instances must be created with the constructor")
	}

	if m.currency != other.Currency() {
		return Money{}, fmt.Errorf("cannot subtract %s from %s: currencies must match", other.Currency(), m.currency)
	}

	diff, err := TrySafeSub(m.amount, other.Amount())
	if err != nil {
		return Money{}, fmt.Errorf("cannot subtract amounts: %w", err)
	}

	return NewMoney(diff, m.currency, m.precision)
}

// Inc adds the given amount to the money.
func (m Money) Inc(amount int64) (Money, error) {
	if !m.initialized {
		return Money{}, errors.New("Money instances must be created with the constructor")
	}

	sum, err := TrySafeAdd(m.amount, amount)
	if err != nil {
		return Money{}, fmt.Errorf("cannot add amounts: %w", err)
	}

	return NewMoney(sum, m.currency, m.precision)
}

// Dec subtracts the given amount from the money.
func (m Money) Dec(amount int64) (Money, error) {
	if !m.initialized {
		return Money{}, errors.New("Money instances must be created with the constructor")
	}

	diff, err := TrySafeSub(m.amount, amount)
	if err != nil {
		return Money{}, fmt.Errorf("cannot subtract amounts: %w", err)
	}

	return NewMoney(diff, m.currency, m.precision)
}

// Times returns a new Money with the product of the two amounts.
func (m Money) Times(amount int64) (Money, error) {
	if !m.initialized {
		return Money{}, errors.New("Money instances must be created with the constructor")
	}

	prod, err := TrySafeMul(m.amount, amount)
	if err != nil {
		return Money{}, fmt.Errorf("cannot multiply amounts: %w", err)
	}

	return NewMoney(prod, m.currency, m.precision)
}

// DividedBy returns a new Money instance with the amount divided by the given divisor.
// Uses half-up rounding: fractions >= 0.5 round up, < 0.5 round down.
func (m Money) DividedBy(divisor int64) (Money, error) {
	if !m.initialized {
		return Money{}, errors.New("Money instances must be created with the constructor")
	}

	quotient, err := TrySafeDiv(m.amount, divisor)
	if err != nil {
		return Money{}, fmt.Errorf("cannot multiply amounts: %w", err)
	}

	remainder := m.amount % divisor

	// For half-up rounding, we need to check if abs(remainder) >= abs(divisor)/2
	// Handle both positive and negative cases
	absReminder := remainder
	if absReminder < 0 {
		absReminder = -absReminder
	}
	absDivisor := divisor
	if absDivisor < 0 {
		absDivisor = -absDivisor
	}

	// Check if we should round up
	// remainder * 2 >= divisor (avoiding division for precision)
	if absReminder*2 >= absDivisor {
		if (m.amount >= 0 && divisor > 0) || (m.amount < 0 && divisor < 0) {
			quotient++
		} else {
			quotient--
		}
	}

	return NewMoney(quotient, m.currency, m.precision)
}

// Equals returns true if the two moneys have the same amount, currency, and precision.
func (m Money) Equals(other Money) bool {
	return m.amount == other.Amount() && m.currency == other.Currency() && m.precision == other.Precision()
}
