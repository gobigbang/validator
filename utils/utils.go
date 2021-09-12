package utils

import (
	"encoding/json"
	"regexp"
)

func JsonPrettyPrint(input interface{}) string {
	str, _ := json.MarshalIndent(input, "", "   ")
	return string(str)
}

func FilterStrings(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

var FuncRegex = regexp.MustCompile(`(?P<RuleName>[a-z0-9_\-]+)(:)?(?P<Args>.+)?`)

func ParseFuncRegex(expr string) (paramsMap map[string]string) {

	match := FuncRegex.FindStringSubmatch(expr)

	if len(match) == 0 {
		return nil
	}

	paramsMap = make(map[string]string)
	for i, name := range FuncRegex.SubexpNames() {
		if name != "" {
			if i > 0 && i <= len(match) {
				paramsMap[name] = match[i]
			}
		}
	}
	return paramsMap
}
