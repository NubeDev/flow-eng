package bacnetio

import (
	"errors"
	"github.com/NubeDev/flow-eng/helpers/rubixos"
	log "github.com/sirupsen/logrus"
	"net"
)

func (inst *Server) getServerConfig() (int, int, string, error) {
	config, err := rubixos.GetBACnetConfig()
	if err != nil {
		return 0, 0, "", err
	}
	if config == nil {
		errMsg := "unable to find bacnet config"
		log.Error(errMsg)
		return 0, 0, "", errors.New(errMsg)
	}
	ip := getIp(config.Iface)
	if ip == "" {
		errMsg := "bacnet config get machine ip: was unable to get ip address"
		return 0, 0, "", errors.New(errMsg)
	}
	return config.Port, config.DeviceId, ip, nil
}

func getIp(port string) string {
	itf, _ := net.InterfaceByName(port)
	item, _ := itf.Addrs()
	var ip net.IP
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if !v.IP.IsLoopback() {
				if v.IP.To4() != nil {
					ip = v.IP
				}
			}
		}
	}
	if ip != nil {
		return ip.String()
	} else {
		return ""
	}
}
