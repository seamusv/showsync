package main

import (
	"fmt"
	"net"
)

func SingleInstance(addr string) error {
	if _, err := net.Listen("tcp", addr); err != nil {
		return fmt.Errorf("instance already running: %v", err)
	}
	return nil
}
