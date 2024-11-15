package presenter

import (
	"fmt"
	"io"
)

func Present(writer io.Writer, output string) {
	if output != "" {
		fmt.Fprintf(writer, "%s\n", output)
	}
}
