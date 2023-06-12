package claude

import (
	"fmt"
	"testing"

	"tebot/pkgs/config"
	"tebot/pkgs/initfunc"
	"tebot/pkgs/logtool"
)

func TestGetMessFromClaude(t *testing.T) {
	type args struct {
		prompt         string
		conversationID string
	}
	tests := []struct {
		name string
		args args
	}{
		{args: args{
			prompt:         "详细的春游计划",
			conversationID: "",
		},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Init()
			logtool.InitEvent("debug", "/tmp/test.log")
			initfunc.InitFun()
			var evendata string
			logtool.SugLog.Info("start test")
			for r := range GetMessFromClaude(tt.args.prompt, tt.args.conversationID) {
				logtool.SugLog.Infof("%#v", r)
				newdata := r.Message.Content.Parts[0]
				if evendata == newdata {
					logtool.SugLog.Info("重复字符串")
					continue
				}

				evendata = newdata

				fmt.Println("-------------------------------------------")
				fmt.Println(newdata)
			}
		})
	}
}
