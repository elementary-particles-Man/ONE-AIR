package internal

import (
    "log"

    "github.com/grandcat/zeroconf"
)

// PublishMDNS announces ONE-AIR over mDNS and returns a shutdown function.
func PublishMDNS(port int) func() {
    server, err := zeroconf.Register("oneair", "_http._tcp", "local.", port, nil, nil)
    if err != nil {
        log.Println("mDNS publish failed:", err)
        return func() {}
    }
    log.Printf("mDNS publish enabled on port %d", port)
    return func() {
        server.Shutdown()
    }
}
