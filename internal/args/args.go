package args

import (
	"fmt"
	"strings"
)

func ParseCommaSeparatedList(input string) (result map[string]string, err error) {
	result = make(map[string]string)
	args := strings.Split(input, ",")
	for _, arg := range args {
		if arg == "" {
			continue
		}
		varVal := strings.Split(arg, "=")
		if len(varVal) != 2 {
			return nil, fmt.Errorf("expected variable and value separated with an '=' sign")
		}
		if varVal[0] == "" {
			return nil, fmt.Errorf("key must not be empty")
		}
		if _, ok := result[varVal[0]]; ok {
			return nil, fmt.Errorf("variable %s can only be assigned once", varVal[0])
		}
		result[varVal[0]] = varVal[1]
	}
	return
}
