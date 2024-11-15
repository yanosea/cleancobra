package util

import (
	"fmt"
	"io"
)

func PrintlnWithWriter(writer io.Writer, a ...any) {
	fmt.Fprintf(writer, "%s\n", a[0])
}
