package goerrors

var (
    err_debug bool = false
)

func GetDebug() bool {
    return err_debug
}

func SetDebug(debug bool) {
    err_debug = debug
}
