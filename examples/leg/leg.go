package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/d0pam1n/dynamixel/network"
	v1 "github.com/d0pam1n/dynamixel/protocol/v1"

	"github.com/d0pam1n/hexapod/components/legs"
	"github.com/d0pam1n/hexapod/math3d"
	"github.com/jacobsa/go-serial/serial"
)

var (
	portName  = flag.String("port", "/dev/ttyUSB0", "the serial port path")
	legBaseId = flag.Int("id", 10, "The base ID of the whole leg (4 servos per leg)")
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

	protocol := v1.New(network)

	var l *legs.Leg

	// Get origins from hexapod legs definitions
	// Path: components/legs/hexapod.go
	if *legBaseId == 10 {
		l = legs.NewLeg(network, 10, "FL", math3d.MakeVector3(-61.167, 24, 98), 300)
	}
	if *legBaseId == 20 {
		l = legs.NewLeg(network, 20, "ML", math3d.MakeVector3(-81, 24, 0), 270)
	}
	if *legBaseId == 30 {
		l = legs.NewLeg(network, 30, "BL", math3d.MakeVector3(-61.167, 24, -98), 240)
	}
	if *legBaseId == 40 {
		l = legs.NewLeg(network, 40, "BR", math3d.MakeVector3(61.167, 24, -98), 120)
	}
	if *legBaseId == 50 {
		l = legs.NewLeg(network, 50, "MR", math3d.MakeVector3(81, 24, 0), 90)
	}
	if *legBaseId == 60 {
		l = legs.NewLeg(network, 60, "FR", math3d.MakeVector3(61.167, 24, 98), 60)
	}

	if l == nil {
		fmt.Printf("unknown leg base ID: %d\n", *legBaseId)
		os.Exit(1)
	}

	v := homeFootPosition(&math3d.ZeroVector3, l, math3d.Pose{})

	l.SetGoal(v)

	err = protocol.Action()
	if err != nil {
		fmt.Printf("protocol.Action error: %s\n", err)
		os.Exit(1)
	}

}

// COPIED FROM components/legs/hexapod.go
// Did not have time to refactor properly.
var stepRadius = 240.0

// homeFootPosition returns a vector in the WORLD coordinate space for the home
// position of the given leg.
func homeFootPosition(offset *math3d.Vector3, leg *legs.Leg, pose math3d.Pose) math3d.Vector3 {
	hyp := math.Sqrt((leg.Origin.X * leg.Origin.X) + (leg.Origin.Z * leg.Origin.Z))
	v := pose.Add(math3d.Pose{*offset, 0, 0, 0}).Add(math3d.Pose{math3d.Vector3{0, 0, 10}, 0, 0, 0}).Add(math3d.Pose{*leg.Origin, leg.Angle, 0, 0}).Add(math3d.Pose{math3d.Vector3{0, 0, stepRadius - hyp}, 0, 0, 0}).Position
	v.Y = 0.0
	return v
}
