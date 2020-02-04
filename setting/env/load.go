package env

import (
	"os"
)
func ReplaceConfigEnv(content []byte) []byte {
	return []byte(os.ExpandEnv(string(content)))
}