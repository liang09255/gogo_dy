package ggIP

import (
	"common/ggLog"
	"fmt"
	"net"
)

func GetIP() string {
	ip := getAllIP()[0]
	ggLog.Infof("IP: %v", ip)
	return ip
}

func getAllIP() (ips []string) {
	ips = make([]string, 0)
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		ggLog.Error("Get Interfaces failed, err:", err)
		return
	}

	// 遍历每个网络接口，获取IP地址
	for _, iface := range interfaces {
		// 跳过无效接口和回环接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口的地址信息
		addrs, err := iface.Addrs()
		if err != nil {
			ggLog.Error("Get Addrs failed, err:", err)
			continue
		}

		// 遍历每个地址，提取IP地址
		for _, addr := range addrs {
			// 检查地址是否是IP地址
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips = append(ips, fmt.Sprintf("%s", ipnet.IP))
				}
			}
		}
	}
	return ips
}
