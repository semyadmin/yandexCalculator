package hostname

import (
	"os"
	"strconv"
)

func GetHostname() string {
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	return hostname + "_" + strconv.Itoa(pid)
}
