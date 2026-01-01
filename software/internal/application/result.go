package application

// Result ist ein generisches Result-Pattern für Service-Rückgaben.
// Wails-kompatibel: Verwendet nur primitive Typen (keine Go error interface).
type Result[T any] struct {
	Data    T      `json:"data"`
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

func OK[T any](data T) Result[T] {
	return Result[T]{Data: data, Success: true, Error: ""}
}

func Fail[T any](err error) Result[T] {
	var zero T
	return Result[T]{Data: zero, Success: false, Error: err.Error()}
}

func FailWithMessage[T any](message string) Result[T] {
	var zero T
	return Result[T]{Data: zero, Success: false, Error: message}
}
