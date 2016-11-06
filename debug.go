package goerrors

var (
    err_debug bool = false
)

// Return the Debug boolean flag which indicates that the stack trace will be provided in errors
func GetDebug() bool {
    return err_debug
}

// Modify the Debug boolean flag for enable or disable the stack trace in errors.
// If the `debug` parameter is true, so the stack trace will be provided in errors.
func SetDebug(debug bool) {
    err_debug = debug
}
