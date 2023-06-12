package config

import (
	"fmt"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{}, // TODO: Add test cases.
	}
	fmt.Println(os.Args)
	var args []string
	for _, arg := range os.Args {
		if arg != "^TestInit$" {
			args = append(args, arg)
		}
	}
	os.Args = args
	configPath := os.Getenv("CONFIG_FILE")
	fmt.Println(configPath)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init()
			fmt.Println(BotConfig)
			fmt.Println(Claudeconfig)
			fmt.Println(ForbiddenWord)

		})
	}
}
