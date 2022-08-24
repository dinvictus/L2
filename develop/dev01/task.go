package main

import (
	"os"

	"github.com/beevik/ntp"
)

// GetCurTime Функция для получения точного времени
func GetCurTime() string {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)

	}
	return time.Format("15:04:05") + "\n"
}

func main() {
	os.Stdout.WriteString(GetCurTime())
}
