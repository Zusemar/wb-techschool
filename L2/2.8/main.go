package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

// getTime печатает текущее время, полученное с NTP-сервера.
func getTime() {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка получения времени:", err)
		os.Exit(1)
	}
	fmt.Println("Точное время по NTP:", t)
}

func main() {
	getTime()
}
