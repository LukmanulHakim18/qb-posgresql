package utility

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// get env string
func GetStrEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		return ""
	}
	return val
}

// get env bool
func GetBoolEnv(key string) bool {
	val := GetStrEnv(key)
	ret, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return ret
}

//

func DebugQueryAndParams(debugMode bool, query string, args []interface{}) {
	if debugMode {
		param := []string{}
		for k, v := range args {
			key := fmt.Sprintf("$%d", k+1)
			param = append(param, fmt.Sprintf("%s=>%v", key, v))
		}
		fmt.Println(fmt.Sprintf("query => %s ", query))
		fmt.Println(fmt.Sprintf("param => %s ", strings.Join(param, ",\n\t ")))
	}

}
