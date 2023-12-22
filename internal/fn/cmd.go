package fn

import (
	"net"
	"os"
	"strconv"
	"strings"
)

func GetIntOrDefaultFromEnv(key string, defaultValue int) int {
	value := strings.TrimSpace(os.Getenv(key))
	num, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return defaultValue
	}
	return int(num)
}

func GetStringOrDefaultFromEnv(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}

func IsValidIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}
