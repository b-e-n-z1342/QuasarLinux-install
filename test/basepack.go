
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// runCommand запускает команду и выводит вывод в реальном времени
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// runFastChroot запускает команду в chroot через fast-chroot
func runFastChroot(args ...string) error {
	chrootArgs := append([]string{"/mnt"}, args...)
	return runCommand("fast-chroot", chrootArgs...)
}

// showDialog показывает диалог выбора ядра
func showDialog() (string, error) {
	cmd := exec.Command("dialog", "--title", "Выберите ядро", "--menu", "Пакеты", "15", "50", "5",
			    "1", "linux-zen",
		     "2", "linux-lts",
		     "3", "linux",
	)

	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("диалог отменен")
	}

	choice := strings.TrimSpace(string(output))
	switch choice {
		case "1":
			return "zen", nil
		case "2":
			return "lts", nil
		case "3":
			return "standard", nil
		default:
			return "", fmt.Errorf("неверный выбор")
	}
}

// installBase устанавливает базовые пакеты для выбранного ядра
func installBase(kernel string) error {
	fmt.Println("Установка базовой системы с ядром", kernel)

	basePackages := []string{
		"terminus-font", "iptables-nft", "base", "base-devel",
		"mkinitcpio", "openrc", "dbus", "dbus-openrc",
		"elogind-openrc", "linux-firmware", "dialog", "acpid",
		"linux-api-headers", "flatpak", "acpid-openrc",
	}

	// Добавляем ядро и заголовки
	switch kernel {
		case "zen":
			basePackages = append(basePackages, "linux-zen", "linux-zen-headers")
		case "lts":
			basePackages = append(basePackages, "linux-lts", "linux-lts-headers")
		case "standard":
			basePackages = append(basePackages, "linux", "linux-headers")
	}

	// Устанавливаем пакеты
	args := append([]string{"/mnt"}, basePackages...)
	if err := runCommand("basestrap", args...); err != nil {
		return fmt.Errorf("ошибка установки пакетов: %w", err)
	}

	// Добавляем репозиторий Flatpak
	fmt.Println("Добавление репозитория Flatpak")
	if err := runFastChroot("flatpak", "remote-add", "--if-not-exists", "flathub", "https://flathub.org/repo/flathub.flatpakrepo"); err != nil {
		return fmt.Errorf("ошибка добавления Flatpak: %w", err)
	}

	return nil
}

// createOSRelease создает файлы идентификации ОС
func createOSRelease() error {
	fmt.Println("Создание файлов идентификации ОС")

	// /usr/lib/os-release
	osRelease := `NAME="Quasar Linux"
	PRETTY_NAME="Quasar Linux (Artix base)"
	ID=quasar
	ID_LIKE=artix
	ANACONDA_ID="quasar"
	VERSION="REV-1"
	VERSION_ID="1.0"
	BUILD_ID="rolling"
	ANSI_COLOR="0;36"
	HOME_URL="https://b-e-n-z1342.github.io"
	LOGO=quasar-logo`

	if err := os.WriteFile("/mnt/usr/lib/os-release", []byte(osRelease), 0644); err != nil {
		return fmt.Errorf("ошибка создания os-release: %w", err)
	}

	// /etc/lsb-release
	lsbRelease := `DISTRIB_ID=Quasar
	DISTRIB_RELEASE=1.0
	DISTRIB_DESCRIPTION="Quasar Linux"
	DISTRIB_CODENAME=rolling`

	if err := os.WriteFile("/mnt/etc/lsb-release", []byte(lsbRelease), 0644); err != nil {
		return fmt.Errorf("ошибка создания lsb-release: %w", err)
	}

	// issue files
	os.WriteFile("/mnt/etc/issue", []byte("Quasar Linux \\r \\l"), 0644)
	os.WriteFile("/mnt/etc/issue.net", []byte("Quasar Linux"), 0644)
	os.WriteFile("/mnt/etc/motd", []byte("Welcome to Quasar Linux!"), 0644)

	return nil
}

// mountVirtualFS монтирует виртуальные файловые системы
func mountVirtualFS() error {
	fmt.Println("Монтирование виртуальных файловых систем")

	mounts := []struct {
		source string
		target string
		fsType string
		options string
	}{
		{"/proc", "/mnt/proc", "proc", ""},
		{"/sys", "/mnt/sys", "sysfs", "rbind"},
		{"/dev", "/mnt/dev", "devtmpfs", "rbind"},
		{"/run", "/mnt/run", "tmpfs", "rbind"},
	}

	for _, m := range mounts {
		var args []string
		if m.options != "" {
			args = append(args, "--"+m.options)
		}
		if m.fsType != "" {
			args = append(args, "-t", m.fsType)
		}
		args = append(args, m.source, m.target)

		if err := runCommand("mount", args...); err != nil {
			return fmt.Errorf("ошибка монтирования %s: %w", m.target, err)
		}
	}

	return nil
}

// generateFstab генерирует fstab
func generateFstab() error {
	fmt.Println("Генерация fstab")
	return runCommand("fstabgen", "-U", "/mnt", ">>", "/mnt/etc/fstab")
}

// setupService настраивает сервис OpenRC
func setupService(service string, packages []string) error {
	fmt.Println("Настройка сервиса", service)

	// Пытаемся добавить в автозагрузку
	if err := runFastChroot("rc-update", "add", service, "default"); err != nil {
		fmt.Println("Установка пакетов для", service)

		// Устанавливаем пакеты если нужно
		if len(packages) > 0 {
			installArgs := append([]string{"-S"}, packages...)
			if err := runFastChroot(installArgs...); err != nil {
				return fmt.Errorf("ошибка установки %s: %w", service, err)
			}
		}

		// Снова пробуем добавить в автозагрузку
		if err := runFastChroot("rc-update", "add", service, "default"); err != nil {
			return fmt.Errorf("ошибка настройки %s: %w", service, err)
		}
	}

	fmt.Println(service, "настроен")
	return nil
}

// copyConfigs копирует конфигурационные файлы
func copyConfigs() error {
	fmt.Println("Копирование конфигурационных файлов")

	// pacman.conf
	if _, err := os.Stat("/installer/configs/pacman.conf"); err == nil {
		os.Remove("/mnt/etc/pacman.conf")
		if err := runCommand("cp", "/installer/configs/pacman.conf", "/mnt/etc/"); err != nil {
			fmt.Println("Предупреждение: ошибка копирования pacman.conf:", err)
		}
	} else {
		fmt.Println("Предупреждение: pacman.conf не найден")
	}

	// mirrorlist
	if _, err := os.Stat("/installer/configs/pacman.d/mirrorlist"); err == nil {
		os.Remove("/mnt/etc/pacman.d/mirrorlist")
		if err := runCommand("cp", "/installer/configs/pacman.d/mirrorlist", "/mnt/etc/pacman.d/"); err != nil {
			fmt.Println("Предупреждение: ошибка копирования mirrorlist:", err)
		}
	} else {
		fmt.Println("Предупреждение: mirrorlist не найден")
	}

	// post-install скрипты
	if _, err := os.Stat("/installer/post"); err == nil {
		if err := runCommand("cp", "-r", "/installer/post", "/mnt/"); err != nil {
			fmt.Println("Предупреждение: ошибка копирования post-скриптов:", err)
		}
	}

	// региональные настройки
	if _, err := os.Stat("/installer/regions"); err == nil {
		if err := runCommand("cp", "-r", "/installer/regions", "/mnt/"); err != nil {
			fmt.Println("Предупреждение: ошибка копирования региональных настроек:", err)
		}
	}

	return nil
}

// createPostInstall создает пост-инсталляционный скрипт
func createPostInstall() error {
	fmt.Println("Создание пост-инсталляционного скрипта")

	postScript := `#!/bin/bash
	git_post() {
	git clone https://github.com/b-e-n-z1342/QuasarLinux-install.git
	cp -r QuasarLinux-install/modules/post  ~/
	cp -r QuasarLinux-install/regions /
	rm -rf QuasarLinux-install
	cd ~/post
	./post_install
}
cd /post
./post_install || git_post`

if err := os.WriteFile("/mnt/usr/local/bin/post_install", []byte(postScript), 0755); err != nil {
	return fmt.Errorf("ошибка создания post_install: %w", err)
}

// Добавляем в .bashrc
bashrcEntry := `
# Auto-run post_install on first login
if [ ! -f /root/.post_install_done ]; then
	post_install && touch /root/.post_install_done
	fi`

	file, err := os.OpenFile("/mnt/root/.bashrc", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ошибка открытия .bashrc: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(bashrcEntry); err != nil {
		return fmt.Errorf("ошибка записи в .bashrc: %w", err)
	}

	return nil
}

// cleanup выполняет финальную очистку
func cleanup() error {
	fmt.Println("Финальная настройка")

	// Убираем автогенерацию motd
	if _, err := os.Stat("/mnt/etc/update-motd.d"); err == nil {
		os.RemoveAll("/mnt/etc/update-motd.d")
	}

	// Симлинк для совместимости
	os.Remove("/mnt/etc/os-release")
	os.Symlink("/usr/lib/os-release", "/mnt/etc/os-release")

	return nil
}

func main() {
	fmt.Println("Установка Quasar Linux")
	fmt.Println()

	// Выбор ядра через dialog
	kernel, err := showDialog()
	if err != nil {
		fmt.Println("Ошибка выбора ядра:", err)
		os.Exit(1)
	}

	// Установка базовой системы
	if err := installBase(kernel); err != nil {
		fmt.Println("Ошибка установки:", err)
		os.Exit(1)
	}

	// Создание идентификации ОС
	if err := createOSRelease(); err != nil {
		fmt.Println("Ошибка создания идентификации:", err)
		os.Exit(1)
	}

	// Монтирование виртуальных ФС
	if err := mountVirtualFS(); err != nil {
		fmt.Println("Ошибка монтирования:", err)
		os.Exit(1)
	}

	// Генерация fstab
	if err := generateFstab(); err != nil {
		fmt.Println("Ошибка генерации fstab:", err)
		os.Exit(1)
	}

	// Настройка сервисов
	services := []struct {
		name string
		pkgs []string
	}{
		{"dbus", []string{"dbus", "dbus-openrc"}},
		{"udev", []string{"udev"}},
		{"elogind", []string{"elogind", "elogind-openrc"}},
		{"acpid", []string{"acpid", "acpid-openrc"}},
	}

	for _, service := range services {
		if err := setupService(service.name, service.pkgs); err != nil {
			fmt.Println("Ошибка настройки", service.name+":", err)
		}
	}

	// Копирование конфигов
	if err := copyConfigs(); err != nil {
		fmt.Println("Ошибка копирования конфигов:", err)
	}

	// Создание пост-инсталляционного скрипта
	if err := createPostInstall(); err != nil {
		fmt.Println("Ошибка создания post-install:", err)
		os.Exit(1)
	}

	// Финальная очистка
	if err := cleanup(); err != nil {
		fmt.Println("Ошибка очистки:", err)
	}

	fmt.Println()
	fmt.Println("Базовая установка завершена!")
}
