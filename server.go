package main

import (
    "flag"
    "fmt"
    "net"
    "os"
    "sync"
    "time"
)

var totalBytesReceived uint64
var mutex sync.Mutex
var startTransmission bool

func main() {
    listenPort := flag.Int("port", 12345, "Server port")
    flag.Parse()

    serverAddr := fmt.Sprintf(":%d", *listenPort)

    udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    udpConn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    defer udpConn.Close()

    fmt.Printf("Server listening on port %d\n", *listenPort)

    buf := make([]byte, 65535)
    for {
        n, addr, err := udpConn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("Error:", err)
            break
        }

        if string(buf[:n]) == "start" {
            fmt.Println("Received start signal from", addr)
            startTransmission = true
            go printReceivedTraffic()
        }

        if startTransmission {
            mutex.Lock()
            totalBytesReceived += uint64(n)
            mutex.Unlock()
        }
    }
}

func printReceivedTraffic() {
    for {
        time.Sleep(time.Second)
        mutex.Lock()
        receivedBytes := totalBytesReceived
        totalBytesReceived = 0
        mutex.Unlock()

        fmt.Printf("Received traffic: %.2f KB/s\n", float64(receivedBytes)/1024)
    }
}
