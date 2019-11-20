package config

import (
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

var selection *ini.Section

func init() {
	cfg, err := ini.Load("config.conf")
	if err != nil {
		panic(err)
	}
	selection = cfg.Section("")
	initConfig()
}

func initConfig() {
	for _, key := range selection.KeyStrings() {
		value := selection.Key(key).String()
		envValue := expandValueEnv(value)
		if strings.Compare(value, envValue) != 0 {
			selection.Key(key).SetValue(envValue)
		}
	}
}

func expandValueEnv(value string) (realValue string) {
	realValue = value
	vLen := len(value)
	// 3 = ${}
	if vLen < 3 {
		return
	}
	// Need start with "${" and end with "}", then return.
	if value[0] != '$' || value[1] != '{' || value[vLen-1] != '}' {
		return
	}
	key := ""
	defaultV := ""
	// value start with "${"
	for i := 2; i < vLen; i++ {
		if value[i] == '|' && (i+1 < vLen && value[i+1] == '|') {
			key = value[2:i]
			defaultV = value[i+2 : vLen-1] // other string is default value.
			break
		} else if value[i] == '}' {
			key = value[2:i]
			break
		}
	}
	realValue = os.Getenv(key)
	if realValue == "" {
		realValue = defaultV
	}
	return
}

func String(key string) string {
	return selection.Key(key).String()
}
