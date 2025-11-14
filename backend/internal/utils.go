package internal

import (
    "crypto/rand"
    "log"
    "math/big"
    "net"
)

// LocalIP returns the first non-loopback IPv4 address.
func LocalIP() string {
    ifaces, err := net.Interfaces()
    if err != nil {
        log.Println("failed to list interfaces:", err)
        return "127.0.0.1"
    }
    for _, iface := range ifaces {
        addrs, err := iface.Addrs()
        if err != nil {
            continue
        }
        for _, addr := range addrs {
            ipnet, ok := addr.(*net.IPNet)
            if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return "127.0.0.1"
}

// RandomPort returns a TCP port between 50000–59999.
func RandomPort() int {
    n, err := rand.Int(rand.Reader, big.NewInt(10000))
    if err != nil {
        return 50000
    }
    return int(n.Int64()) + 50000
}

// ShortCode generates a 6-char code for smartphone redirect.
func ShortCode() string {
    const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
    b := make([]byte, 6)
    for i := range b {
        n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
        if err != nil {
            b[i] = 'A'
            continue
        }
        b[i] = alphabet[n.Int64()]
    }
    return string(b)
}
