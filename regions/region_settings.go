package main

import (
	"fmt"
	"os/exec"
)
func main() {
	fmt.Println("Выберите регион:")
	
	fmt.Println("1)  Северная Америка")
	fmt.Println("2)  Европа")
	fmt.Println("3)  Азия")
	fmt.Println("4)  Африка")
	fmt.Println("5)  Антарктида")
	fmt.Println("6)  Арктика")
	fmt.Println("7)  Атлантика")
	fmt.Println("8)  Австралия")
	fmt.Println("9)  Бразилия")
	fmt.Println("10) Канада")
	fmt.Println("11) Etc")	
	fmt.Println("Выберите [1-11]:  ")
	var choice int
	fmt.Scan(&choice)


	switch choice {
	case 1:
		exec.Command("./region", "US").Run()
	case 2:
		exec.Command("./region", "Europe").Run()
	case 3: 
		exec.Command("./region", "Asia").Run()
	case 4:
		exec.Command("./region", "Africa").Run()
	case 5:
		exec.Command("./region", "Antarctica").Run()
	case 6:
		exec.Command("./region", "Arctica").Run()
	case 7:
		exec.Command("./region", "Atlantic").Run()
	case 8:
		exec.Command("./region", "Australia").Run()
	case 9:
		exec.Command("./region", "Brazil").Run()
	case 10:
		exec.Command("./region", "Canada").Run()
	case 11:
		exec.Command("./region", "Etc").Run()
		
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
