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
