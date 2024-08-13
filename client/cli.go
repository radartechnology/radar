package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func stringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)

	for {
		_, err := fmt.Fprint(os.Stderr, label+" ")
		if err != nil {
			break
		}

		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	return strings.TrimSpace(s)
}
