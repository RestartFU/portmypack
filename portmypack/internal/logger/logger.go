package logger

import "fmt"

func Debugf(str string, args ...any) {
	fmt.Printf("DEBU| %s\n", fmt.Sprintf(str, args...))
}
