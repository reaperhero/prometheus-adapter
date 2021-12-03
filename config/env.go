package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnvWithDeafult(name string, value string) string {
	v := os.Getenv(name)
	if v != "" {
		return v
	}
	return value
}

func GetArrayEnvWithDeafult(name string, value string) []string {
	v := os.Getenv(name)
	if v != "" {
		return strings.Split(v, ",")
	}
	return strings.Split(value, ",")
}

func GetEnvInt64WithDefault(name string, value int64) int64 {
	v := os.Getenv(name)
	if v != "" {
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	}
	return value
}

func GetEnvIntWithDefault(name string, value int) int {
	v := os.Getenv(name)
	if v != "" {
		i, _ := strconv.Atoi(v)
		return i
	}
	return value
}

func GetEnvTimeWithDefault(name string, value time.Duration) time.Duration {
	v := os.Getenv(name)
	if v != "" {
		i, _ := strconv.Atoi(v)
		return time.Second * time.Duration(i) / 1000
	}
	return value
}
