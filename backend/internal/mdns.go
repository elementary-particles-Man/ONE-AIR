package internal

import (
	"log"
	"strings"

	"github.com/grandcat/zeroconf"
)

// PublishMDNS announces ONE-AIR over mDNS with a resolvable hostname.
func PublishMDNS(host string, port int, ip string) func() {
	host = strings.TrimSpace(host)
	ip = strings.TrimSpace(ip)
	if host == "" || ip == "" {
		log.Println("mDNS publish skipped: host or IP missing")
		return func() {}
	}
	server, err := zeroconf.RegisterProxy("oneair", "_http._tcp", "local.", port, host, []string{ip}, nil, nil)
	if err != nil {
		log.Println("mDNS publish failed:", err)
		return func() {}
	}
	log.Printf("mDNS publish enabled for %s:%d (%s)", host, port, ip)
	return func() {
		server.Shutdown()
	}
}
