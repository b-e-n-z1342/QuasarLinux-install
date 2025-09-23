package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// runChroot — запускает команду внутри chroot через fast-chroot
func runChroot(args ...string) error {
	cmdArgs := append([]string{"/mnt"}, args...)
	cmd := exec.Command("fast-chroot", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// input — читает строку с клавиатуры
func input(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// createUser — создаёт пользователя и группу
func createUser() error {
	username := input("→ Введите имя пользователя: ")
	if username == "" {
		return fmt.Errorf("имя пользователя не может быть пустым")
	}

	fmt.Printf(" Создаём группу %s...\n", username)
	runChroot("groupadd", username) // игнорируем ошибку, если группа есть

	fmt.Printf(" Создаём пользователя %s...\n", username)
	if err := runChroot("useradd", "-m", "-g", username, username); err != nil {
		return fmt.Errorf("не удалось создать пользователя: %w", err)
	}

	fmt.Printf(" Установите пароль для %s:\n", username)
	if err := runChroot("passwd", username); err != nil {
		return fmt.Errorf("не удалось установить пароль: %w", err)
	}

	return nil
}

// setRootPassword — устанавливает пароль root
func setRootPassword() error {
	fmt.Println("  Установите пароль для root:")
	return runChroot("passwd")
}

func main() {
	fmt.Println("=====================================")
	fmt.Println("      Настройка пользователей")
	fmt.Println("=====================================")

	if err := createUser(); err != nil {
		fmt.Fprintf(os.Stderr, " Ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	if err := setRootPassword(); err != nil {
		fmt.Fprintf(os.Stderr, " Ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("  Пользователи успешно настроены!")
}
