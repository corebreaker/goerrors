package goerrors

var (
	errDebug bool = false
)

// Return the Debug boolean flag which indicates that the stack trace will be provided in errors
func GetDebug() bool {
	return errDebug
}

// Modify the Debug boolean flag for enable or disable the stack trace in errors.
// If the `debug` parameter is true, so the stack trace will be provided in errors.
func SetDebug(debug bool) {
	errDebug = debug
}
