package tools

import (
	"regexp"
	"strings"
	"tebot/pkgs/config"
	"tebot/pkgs/initfunc"
	//"tebot/pkgs/msg"
)

var (
	forbiddenWordRegexp config.ForbiddenWordStruct
)

func init() {
	initfunc.RegisterInitFunc(
		func() {
			forbiddenWordRegexp = config.ForbiddenWord
			forbiddenWordRegexp.Attack = []string{`(?i)` + strings.Join(forbiddenWordRegexp.Attack, "|")}
			forbiddenWordRegexp.Politics = []string{strings.Join(forbiddenWordRegexp.Politics, "|")}
		},
	)
}
func ChecText(t string) string {
	var content string
	if regexp.MustCompile(forbiddenWordRegexp.Politics[0]).MatchString(t) {
		content = "检测到政治敏感词汇"
	} else if regexp.MustCompile(forbiddenWordRegexp.Attack[0]).MatchString(t) {
		content = "检测到攻击性敏感词汇"
	}
	return content
}
