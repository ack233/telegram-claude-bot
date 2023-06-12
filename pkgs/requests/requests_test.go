package requests

import (
	"fmt"
	"tebot/pkgs/logtool"
	"testing"
)

func TestEcocde_json(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logtool.InitEvent("info", "/tmp/cc")
			logtool.SugLog.Info("开始测试")
			Client = getclinet()
			resp, err := Client.R().EnableTrace().
				SetBody(map[string]string{"aa": "bb"}).
				SetHeader("aa", "m").
				Post("https://www.baidu.com")

			// Explore response object
			fmt.Println("Response Info:")
			fmt.Println("  Error      :", err)
			fmt.Println("  Status Code:", resp.StatusCode())
			fmt.Println("  Status     :", resp.Status())
			fmt.Println("  Proto      :", resp.Proto())
			fmt.Println("  Time       :", resp.Time())
			fmt.Println("  Received At:", resp.ReceivedAt())
			fmt.Println("  Body       :\n", resp)
			fmt.Println()

			// Explore trace info
			fmt.Println("Request Trace Info:")
			ti := resp.Request.TraceInfo()
			fmt.Println("  DNSLookup     :", ti.DNSLookup)
			fmt.Println("  ConnTime      :", ti.ConnTime)
			fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
			fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
			fmt.Println("  ServerTime    :", ti.ServerTime)
			fmt.Println("  ResponseTime  :", ti.ResponseTime)
			fmt.Println("  TotalTime     :", ti.TotalTime)
			fmt.Println("  IsConnReused  :", ti.IsConnReused)
			fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
			fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
			fmt.Println("  RequestAttempt:", ti.RequestAttempt)
			fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
		})
	}
}
