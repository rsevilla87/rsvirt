package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y]/n: ", s)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" || response == "" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
