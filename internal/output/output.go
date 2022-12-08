package output

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func Error(message string, err error) {
	fmt.Println(Red(fmt.Sprintf("Error: %s\n%s", message, err)))
}

func Info(message string) {
	fmt.Printf("Info: %s\n", message)
}

func Warning(message string) {
	fmt.Println(Yellow(fmt.Sprintf("Warning: %s\n", message)))
}

func Complete(message string) {
	fmt.Println(Green(fmt.Sprintf("Complete: %s\n", message)))
}
