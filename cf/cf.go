package cf

import (
	"fmt"
	"time"

	"github.com/XIU2/CloudflareSpeedTest/task"
	"github.com/XIU2/CloudflareSpeedTest/utils"
)

var minDelay, maxDelay, downloadTime int
var maxLossRate float64

type CloudFlareOptions struct {
	EnableCloudFlare bool `json:"enable-cloudflare"`
	CloudFlareIPNum  int  `json:"cloudflare-ip-num"`
	CloudFlareIPs    []string
}

func init() {
	fmt.Println("cf init")
	minDelay = 0
	maxDelay = 9999
	downloadTime = 10
	maxLossRate = 1.0
	utils.InputMaxDelay = time.Duration(maxDelay) * time.Millisecond
	utils.InputMinDelay = time.Duration(minDelay) * time.Millisecond
	utils.InputMaxLossRate = float32(maxLossRate)
	task.Timeout = time.Duration(downloadTime) * time.Second
	task.HttpingCFColomap = task.MapColoMap()
	task.Routines = 200
	task.PingTimes = 4
	task.TestCount = 10
	task.TCPPort = 443
	task.URL = "https://cf.xiu2.xyz/url"
	task.Httping = false
	task.HttpingStatusCode = 0
	task.HttpingCFColo = ""
	utils.PrintNum = 10
	task.IPFile = "ip.txt"
	utils.Output = "result.csv"
	task.Disable = false
	task.TestAll = false
	task.TestIP = false
	task.TestIPNum = 200
	task.IPText = "173.245.48.0/20, 103.21.244.0/22, 103.22.200.0/22, 103.31.4.0/22, 141.101.64.0/18, 108.162.192.0/18, 190.93.240.0/20, 188.114.96.0/20, 197.234.240.0/22, 198.41.128.0/17, 162.158.0.0/15, 104.16.0.0/13, 104.24.0.0/14, 172.64.0.0/13, 131.0.72.0/22"
}

func Run(base CloudFlareOptions) CloudFlareOptions {
	task.InitRandSeed() // set random seed
	task.TestIP = base.EnableCloudFlare
	task.TestIPNum = base.CloudFlareIPNum
	pingData := task.NewPing().Run().FilterDelay().FilterLossRate()
	// jsonData, _ := json.Marshal(pingData)
	var CloudflareIPList []string
	if len(pingData) == 0 {
		return base
	}
	for _, v := range pingData[:10] {
		CloudflareIPList = append(CloudflareIPList, v.IP.String())
	}
	base.CloudFlareIPs = CloudflareIPList
	return base
}
