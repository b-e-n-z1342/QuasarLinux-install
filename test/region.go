package main

import (
	"fmt"
	"os"
	"os/exec"
)
var regions = map[string][]string{
	"Europe": {
		runCommand("ls /usr/share/zoneinfo/Europe/")
	},
	"Asia": {
		runCommand("ls /usr/share/zoneifo/Asia")
	},
	"North_America" {
		runCommand("ls /usr/share/zoneinfo/America")
	},
}
func main() {
	// Проверяем аргументы
	if len(os.Args) < 2 {
		fmt.Println("Использование: ./region <континент>")
		fmt.Println("Доступные континенты: Europe, Asia, US")
		os.Exit(1)
	}
	continent := os.Args[1]
	

	regionList, exists := regions[continent]
	if !exists {
		fmt.Printf("Неизвестный континент: %s\n", continent)
		fmt.Println("Доступные континенты: Europe, Asia, US")
		os.Exit(1)
	}

