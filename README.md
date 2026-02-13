# Pocket

<img src="./pocket_logo.png" width="400">

Pocket provides utility functions for common patterns that aren't in the Go standard library but I find useful in almost every project.

No external dependencies, just the standard library.

You can add it as a dependency to your project or you can just copy the parts that you find useful.

## Installation

```bash
go get github.com/germanDV/pocket
```

## Directory Functions

### `HomeDir`
Returns the home directory of the current user.

```go
home, _ := pocket.HomeDir()
fmt.Println("Home:", home)
```

### `ConfigDir`
Returns the configuration directory.

```go
configDir, _ := pocket.ConfigDir()
fmt.Println("Config:", configDir)
```

### `DataDir`
Returns the data directory.

```go
dataDir, _ := pocket.DataDir()
fmt.Println("Data:", dataDir)
```

## Slice Functions

### `Map`
Applies a function to each element of a slice and returns a new slice with the results.

```go
numbers := []int{1, 2, 3}
doubled := pocket.Map(numbers, func(n int) int {
    return n * 2
})
// doubled = [2, 4, 6]
```

### `Filter`
Returns a new slice containing only elements for which the predicate function returns true.

```go
numbers := []int{1, 2, 3, 4, 5}
even := pocket.Filter(numbers, func(n int) bool {
    return n%2 == 0
})
// even = [2, 4]
```

## Safe Math Functions

### `SafeAdd`
Returns the sum of two integers, panicking if the result overflows or underflows.

```go
sum := pocket.SafeAdd(100, 200)            // 300
sum = pocket.SafeAdd(uint8(255), uint8(1)) // panics: unsigned integer overflow
```

### `SafeSub`
Returns the difference of two integers, panicking if the result overflows or underflows.

```go
diff := pocket.SafeSub(100, 50)           // 50
diff = pocket.SafeSub(uint8(0), uint8(1)) // panics: unsigned integer underflow
```

## Test Assertion Functions

- `AssertNotNil` Asserts that the given value is not nil.
- `AssertNil` Asserts that the given value is nil.
- `AssertTrue` Asserts that the given value is true.
- `AssertFalse` Asserts that the given value is false.
- `AssertEqual` Asserts that two values are deeply equal.
- `AssertNotEqual` Asserts that two values are not deeply equal.
- `AssertErrorIs` Asserts that an error is of the expected type using `errors.Is`.
- `AssertContains` Asserts that a string contains a substring.
- `AssertPanics` Asserts that the given function panics.

## Configuration Functions

### `LoadConfigFromEnv`
Populates a config struct from environment variables.
Uses `env` and `default` struct tags.

```go
type AppConfig struct {
    Port     int           `env:"PORT" default:"8080"`
    LogLevel string        `env:"LOG_LEVEL"`
    Timeout  time.Duration `env:"TIMEOUT" default:"10s"`
}

config, err := pocket.LoadConfigFromEnv[AppConfig]()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Port: %d\n", config.Port)
```

Supported types: `string`, `int`, `bool`, `time.Duration`

## String Functions

### `SafeCompare`
Performs a constant-time comparison of two strings to protect against timing attacks. It hashes both strings using SHA-256 to ensure they have the same length before comparison.

```go
result := pocket.SafeCompare("secret_token", "secret_token") // true
result = pocket.SafeCompare("token1", "token2")              // false
```

### `GenerateString`
Generates a random string of the specified length using `crypto/rand`. The result is base64 URL-encoded. Note: The returned string will be longer than the input length due to base64 encoding. Panics if random number generation fails.

```go
token := pocket.GenerateString(32) // Random URL-safe string
```

## Money Functions

Pocket provides a `Money` type for working with monetary values. Money instances are immutable and support safe arithmetic operations with overflow protection.

### `NewUSD`
Creates a new Money instance with USD currency (precision 2).

```go
m := pocket.NewUSD(100_99)
fmt.Println(m.Format())   // "100.99 USD"
```

### `NewARS`
Creates a new Money instance with ARS currency (precision 2).

```go
m := pocket.NewARS(50_00)
fmt.Println(m.Format())  // "50.00 ARS"
```

### `NewMoney`
Creates a new Money instance with custom currency and precision.

```go
m, err := pocket.NewMoney(100_99, "EUR", 2)
if err != nil {
    log.Fatal(err)
}
fmt.Println(m.String()) // "100.99"
```

### `NewMoneyFromString`
Creates a new Money instance from a string in the format "amount currency".
The number of decimal places determines the precision. So be sure to include 0s if necessary.
And be careful not to use unsanitized user input as "100 USD" will be different from "100.00 USD".

```go
m, err := pocket.NewMoneyFromString("100.99 USD")
if err != nil {
    log.Fatal(err)
}
fmt.Println(m.Amount()) // 10099

// Works with any precision
btc, _ := pocket.NewMoneyFromString("1.00000000 BTC")
fmt.Println(btc.Precision()) // 8

// Currency is case-insensitive
m2, _ := pocket.NewMoneyFromString("50.00 usd")
fmt.Println(m2.Currency()) // "USD"
```

### `Money.Currency`
Returns the currency code of the money.

```go
m := pocket.NewUSD(100)
fmt.Println(m.Currency()) // "USD"
```

### `Money.Precision`
Returns the precision (number of decimal places) of the money.

```go
m := pocket.NewUSD(100)
fmt.Println(m.Precision()) // 2
```

### `Money.Amount`
Returns the amount in the smallest unit of the currency (e.g., cents for USD).

```go
m := pocket.NewUSD(100_99)
fmt.Println(m.Amount()) // 10099 (represents $100.99)
```

### `Money.String`
Returns the amount as a formatted string with proper decimal places.

```go
m := pocket.NewUSD(100_99)
fmt.Println(m.String()) // "100.99"
```

### `Money.Format`
Returns the amount and currency in "amount currency" format.

```go
m := pocket.NewUSD(100_99)
fmt.Println(m.Format()) // "100.99 USD"
```

### `Money.Plus`
Returns a new Money with the sum of two amounts. Returns an error if currencies don't match or if overflow occurs.

```go
m1 := pocket.NewUSD(100_00)
m2 := pocket.NewUSD(50_00)
result, err := m1.Plus(m2)
// result = 15000 ($150.00)
```

### `Money.Minus`
Returns a new Money with the difference of two amounts. Returns an error if currencies don't match or if overflow occurs.

```go
m1 := pocket.NewUSD(100_00)
m2 := pocket.NewUSD(30_00)
result, err := m1.Minus(m2)
// result = 7000 ($70.00)
```

### `Money.Inc`
Adds a raw amount to the money and returns a new Money instance.

```go
m := pocket.NewUSD(100_00)
result, err := m.Inc(50_00)
// result = 15000 ($150.00)
```

### `Money.Dec`
Subtracts a raw amount from the money and returns a new Money instance.

```go
m := pocket.NewUSD(100_00)
result, err := m.Dec(30_00)
// result = 7000 ($70.00)
```

### `Money.Times`
Multiplies the money by a factor and returns a new Money instance.

```go
m := pocket.NewUSD(10_00)
result, err := m.Times(3)
// result = 3000 ($30.00)
```

### `Money.DividedBy`
Divides the money by a divisor and returns a new Money instance. Uses half-up rounding.

```go
m := pocket.NewUSD(100_00)
result, err := m.DividedBy(3)
// result = 3333 ($33.33)
```

### `Money.Equals`
Returns true if two Money instances have the same amount, currency, and precision.

```go
m1 := pocket.NewUSD(100_00)
m2 := pocket.NewUSD(100_00)
fmt.Println(m1.Equals(m2)) // true
```
