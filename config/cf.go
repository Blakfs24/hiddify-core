package config

import (
	"fmt"
	"math/rand"
	"net"
	"sort"
	"time"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

// CloudflareIPRanges 存储Cloudflare的IP地址段
var CloudflareIPRanges = []string{
	"173.245.48.0/20",
	"103.21.244.0/22",
	"103.22.200.0/22",
	"103.31.4.0/22",
	"141.101.64.0/18",
	"108.162.192.0/18",
	"190.93.240.0/20",
	"188.114.96.0/20",
	"197.234.240.0/22",
	"198.41.128.0/17",
	"162.158.0.0/15",
	"104.16.0.0/13",
	"104.24.0.0/14",
	"172.64.0.0/13",
	"131.0.72.0/22",
}

// IPToUint32 将IP地址转换为32位无符号整数
func IPToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// CIDRToRange 将CIDR表示的IP地址段转换为起始和结束IP的32位无符号整数表示
func CIDRToRange(cidr string) (uint32, uint32, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0, 0, err
	}
	startIP := IPToUint32(ipnet.IP)
	mask := IPToUint32(net.IP(ipnet.Mask))
	endIP := startIP | ^mask
	return startIP, endIP, nil
}

// IsIPInCloudflareRanges 检查IP地址是否在Cloudflare的IP地址段中
func IsIPInCloudflareRanges(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	ipUint32 := IPToUint32(ip)

	ranges := make([][2]uint32, len(CloudflareIPRanges))
	for i, cidr := range CloudflareIPRanges {
		start, end, err := CIDRToRange(cidr)
		if err != nil {
			return false
		}
		ranges[i] = [2]uint32{start, end}
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][0] < ranges[j][0]
	})

	index := sort.Search(len(ranges), func(i int) bool {
		return ranges[i][1] >= ipUint32
	})

	if index < len(ranges) && ranges[index][0] <= ipUint32 && ranges[index][1] >= ipUint32 {
		return true
	}
	return false
}

func patchCloudflare(base *option.Outbound, configOpt *ConfigOptions) error {
	fmt.Println(base.SocksOptions.Server, " ", base.Tag)
	return nil
}

func RandomSelect[T any](list []T, n int) []T {
	if n >= len(list) {
		return list
	}
	// use current time as seed
	rand.Seed(time.Now().UnixNano())
	// copy list to avoid modifying the original list
	copiedList := append([]T(nil), list...)
	// shuffle the copied list
	for i := len(copiedList) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		copiedList[i], copiedList[j] = copiedList[j], copiedList[i]
	}
	return copiedList[:n]
}

func patchXragConfig(base option.Outbound, configOpt ConfigOptions, obj outboundMap) outboundMap {
	if base.Type == C.TypeVLESS && configOpt.XrayOptions.EnableVless {
		obj["type"] = "xvless"
	}
	if base.Type == C.TypeVMess && configOpt.XrayOptions.EnableVmess {
		obj["type"] = "xvmess"
	}
	if base.Type == C.TypeTrojan && configOpt.XrayOptions.EnableTrojan {
		obj["type"] = "xtrojan"
	}
	return obj
}
