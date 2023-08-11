package main

import (
    "flag"
    "fmt"
    "net"
    "os"
    "time"
)

func main() {
    serverIP := flag.String("server", "127.0.0.1", "Server IP address")
    serverPort := flag.Int("port", 12345, "Server port")
    bandwidthKbps := flag.Int("bw", 1000, "Bandwidth in Kbps")
    interfaceName := flag.String("interface", "eth0", "Network interface name")
    flag.Parse()

    serverAddr := fmt.Sprintf("%s:%d", *serverIP, *serverPort)
    localAddr, err := getLocalIP(*interfaceName)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    fmt.Printf("Sending initial packet to start traffic synchronization to server %s from local address %s\n", serverAddr, localAddr)

    udpConn, err := net.Dial("udp", serverAddr)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    defer udpConn.Close()

    initialPacket := []byte("start")
    _, err = udpConn.Write(initialPacket)
    if err != nil {
        fmt.Println("Error sending initial packet:", err)
        os.Exit(1)
    }

    payloadSize := 1472 // Typical MTU size
    interval := time.Second / time.Duration(*bandwidthKbps)
    payload := make([]byte, payloadSize)

    for {
        start := time.Now()
        _, err := udpConn.Write(payload)
        if err != nil {
            fmt.Println("Error:", err)
            break
        }

        elapsed := time.Since(start)
        if elapsed < interval {
            time.Sleep(interval - elapsed)
        }
    }
}

func getLocalIP(interfaceName string) (string, error) {
    iface, err := net.InterfaceByName(interfaceName)
    if err != nil {
        return "", err
    }

    addrs, err := iface.Addrs()
    if err != nil {
        return "", err
    }

    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String(), nil
            }
        }
    }

    return "", fmt.Errorf("no suitable local IP address found")
}

