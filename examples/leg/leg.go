package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/d0pam1n/dynamixel/network"
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

	var l *legs.Leg

	// Get origins from hexapod legs definitions
	// Path: components/legs/hexapod.go
	if *legBaseId == 10 {
		l = legs.NewLeg(network, legs.OriginFrontLeft.BaseId, legs.OriginFrontLeft.Name, legs.OriginFrontLeft.Vector, legs.OriginFrontLeft.Angle)
	}
	if *legBaseId == 20 {
		l = legs.NewLeg(network, legs.OriginMidLeft.BaseId, legs.OriginMidLeft.Name, legs.OriginMidLeft.Vector, legs.OriginMidLeft.Angle)
	}
	if *legBaseId == 30 {
		l = legs.NewLeg(network, legs.OriginBackLeft.BaseId, legs.OriginBackLeft.Name, legs.OriginBackLeft.Vector, legs.OriginBackLeft.Angle)
	}
	if *legBaseId == 40 {
		l = legs.NewLeg(network, legs.OriginBackRight.BaseId, legs.OriginBackRight.Name, legs.OriginBackRight.Vector, legs.OriginBackRight.Angle)
	}
	if *legBaseId == 50 {
		l = legs.NewLeg(network, legs.OriginMidRight.BaseId, legs.OriginMidRight.Name, legs.OriginMidRight.Vector, legs.OriginMidRight.Angle)
	}
	if *legBaseId == 60 {
		l = legs.NewLeg(network, legs.OriginFrontRight.BaseId, legs.OriginFrontRight.Name, legs.OriginFrontRight.Vector, legs.OriginFrontRight.Angle)
	}

	if l == nil {
		fmt.Printf("unknown leg base ID: %d\n", *legBaseId)
		os.Exit(1)
	}

}

// COPIED FROM components/legs/hexapod.go
// Did not have time to refactor properly.
var stepRadius = 120.0

// homeFootPosition returns a vector in the WORLD coordinate space for the home
// position of the given leg.
func homeFootPosition(offset *math3d.Vector3, leg *legs.Leg, pose math3d.Pose) math3d.Vector3 {
	hyp := math.Sqrt((leg.Origin.X * leg.Origin.X) + (leg.Origin.Z * leg.Origin.Z))
	v := pose.Add(math3d.Pose{*offset, 0, 0, 0}).Add(math3d.Pose{math3d.Vector3{0, 0, 10}, 0, 0, 0}).Add(math3d.Pose{*leg.Origin, leg.Angle, 0, 0}).Add(math3d.Pose{math3d.Vector3{0, 0, stepRadius - hyp}, 0, 0, 0}).Position
	v.Y = 0.0
	return v
}
