package main

import (
	"net"
	"time"
)

const (
	ALIVE = iota
	SUNK = iota
	WAITFORSPAWN = iota
	OBSERVE = iota
)

type vector struct{
	x float32
	y float32
	z float32
}

type playerInfo struct{
	client net.Conn
	transform *playerTransform
	agentID int16
}

type playerTransform struct{
	playerPosition vector
	playerRotation vector
	agentID int16
	locationTime float32
}



type event struct{
	eventPosition vector
	eventRange float32
}

type agent struct{
	ID int16
	name string
}

type shipState struct{
	ID int16
	doorsOpen bool
	destructionState []bool
	name string
	state int16
	startTime time.Time
	canSpawn bool
	alive bool
}
