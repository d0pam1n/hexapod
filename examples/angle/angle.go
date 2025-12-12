package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/d0pam1n/dynamixel/network"
	"github.com/d0pam1n/hexapod/servos"
	"github.com/jacobsa/go-serial/serial"
)

var (
	portName  = flag.String("port", "/dev/ttyUSB0", "the serial port path")
	legBaseId = flag.Int("id", 10, "The base ID of the whole leg (4 servos per leg)")
	angle     = flag.Float64("angle", 512, "the goal angle to set")
	debug     = flag.Bool("debug", false, "show serial traffic")
)

func main() {
	flag.Parse()

	options := serial.OpenOptions{
		PortName:              *portName,
		BaudRate:              1000000,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	}

	serial, err := serial.Open(options)
	if err != nil {
		fmt.Printf("open error: %s\n", err)
		os.Exit(1)
	}

	network := network.New(serial)
	if *debug {
		network.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	network.Timeout = 1 * time.Second

	s, err := servos.New(network, *legBaseId)
	if err != nil {
		fmt.Printf("servos.New error: %s\n", err)
		os.Exit(1)
	}

	err = s.MoveTo(*angle)
	if err != nil {
		fmt.Printf("servos.MoveTo error: %s\n", err)
		os.Exit(1)
	}
}
