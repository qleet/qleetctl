package output

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func Error(message string, err error) {
	fmt.Printf(Red("Error: %s\n%s", message, err))
}

func Info(message string) error {
	fmt.Printf("Info: %s\n", message)
}

func Warning(message string) error {
	fmt.Printf(Yellow("Warning: %s\n", message))

func Complete(message string) error {
	fmt.Printf(Green("Complete: %s\n", message))
}
