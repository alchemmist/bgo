package app

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func handleError(in io.Reader, out io.Writer, err error, message string) {
	fmt.Fprintln(out, message)
	fmt.Fprint(out, "Do you want to see the full error? [y/N]: ")

	reader := bufio.NewReader(in)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer == "y" {
		fmt.Fprintf(out, "%v\n", err)
	}
}
