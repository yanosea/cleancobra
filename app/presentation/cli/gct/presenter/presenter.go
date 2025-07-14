package presenter

import (
	"fmt"
	"io"
)

func Present(writer io.Writer, output string) {
	if output != "" {
		_, _ = fmt.Fprintf(writer, "%s\n", output)
	}
}
