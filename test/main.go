package main

import (
	"fmt"
	"os/exec"
)
func main() {
	fmt.Println("Выберите регион:")
	fmt.Println("1 Северная Америка")
	fmt.Println("2 Южная Америка")
	fmt.Println("3 Европа")
	fmt.Println("4 Азия")
	fmt.Println("5 Африка")
	
	var choice int
	fmt.Scan(&choice)


	switch choice {
	case 1:
		runCommand("./region", "US")
	case 2:
		runCommand("./region", "SA")
	case 3: 
		runCommand("./region", "Europe")
	case 4:
		runCommand("./region", "Asia")
	case 5: 
		runCommand("./region", "Africa")
	default:
		fmt.Println("Неверный выбор!")
	}
}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Ошибка: %v\n", err)
		return
	}
	runCommand("clear")
	fmt.Println("Результат:\n%s\n", output)
}
