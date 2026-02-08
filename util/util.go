package util

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

var Log *log.Logger = log.New(os.Stdout, "magic ", log.Default().Flags())

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Generate a random port
func RandomPort(start, end uint) uint {
	return start + uint(rand.Intn(int(end-start+1)))
}

// Scan an individual port. Returns true when the creation of the listener succeeds.
func ScanPort(port uint) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err == nil {
		listener.Close()
		return true
	}
	return false
}
