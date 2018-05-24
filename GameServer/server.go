package main

import (
	"net"
	"log"

	"sync"
	"time"

	"bytes"
	"encoding/binary"
)
var clients map[int16]*playerInfo
var playerStates map[int16]*playerTransform
var shipStates map[int16]*shipState
var clientStats map[int16]playerStats
var actions *Queue
var stateLock sync.RWMutex
var currentPlayerID int16 = 0
var IDlock sync.Mutex
var ActionLock sync.Mutex
var serverUpdateRate chan int
var serverTickRate int16 = 20
var serverTime = time.Now()
var abilityList map[int16]abilityData
var spawnLocations = []vector{vector{0,0,0},vector{400,0,-400},vector{650,0,-180},vector{0,0,-400}}
var respawnTime = 15 * time.Second
var nextSpawn chan int
var startingHealth float32 = 110.0
var playerEffects []ability
var effectLock sync.Mutex
var itemList map[int16][]baseItem

var resetItemsOnDeath bool

func spawnLoc(){
	next := 0
	for {
		nextSpawn <- next
		next++
		if next >= len(spawnLocations){
			next = 0
		}
	}
}


func checkError(err error) bool{
	if err != nil{
		log.Println(err)
		return true
	}
	return false
}




func processState(){
	ticker := time.NewTicker(time.Millisecond * time.Duration(1000/serverTickRate))
	for{
		_ = <- ticker.C
		for _, ps := range shipStates{
			if clientStats[ps.ID].baseStats[CURRENT_HEALTH] <= 0 && ps.state == ALIVE{
				action1 := new(bytes.Buffer)
				action1.WriteByte(EXPLODE)
				binary.Write(action1,binary.LittleEndian,clients[ps.ID].transform.playerPosition.x)
				binary.Write(action1,binary.LittleEndian,clients[ps.ID].transform.playerPosition.y)
				binary.Write(action1,binary.LittleEndian,clients[ps.ID].transform.playerPosition.z)
				actions.Push(action1.Bytes())
				action2 := new(bytes.Buffer)
				action2.WriteByte(SINK)
				binary.Write(action2,binary.LittleEndian,ps.ID)
				binary.Write(action2,binary.LittleEndian,float32(0))
				binary.Write(action2,binary.LittleEndian,float32(0))
				binary.Write(action2,binary.LittleEndian,float32(1))
				binary.Write(action2,binary.LittleEndian,float32(60))
				binary.Write(action2,binary.LittleEndian,float32(60))
				binary.Write(action2,binary.LittleEndian,float32(0.1))
				actions.Push(action2.Bytes())
				ps.state = SUNK
				ps.startTime = time.Now()
				action3 := new(bytes.Buffer)
				binary.Write(action3,binary.LittleEndian,int16(9))
				binary.Write(action3,binary.LittleEndian,int16(0))
				binary.Write(action3,binary.LittleEndian,int16(1))
				action3.WriteByte(RESPAWN)
				binary.Write(action3,binary.LittleEndian,float32(respawnTime.Seconds()))

				clients[ps.ID].client.Write(action3.Bytes())
			}
			if ps.state == SUNK && time.Now().Sub(ps.startTime) > respawnTime{
				action := new(bytes.Buffer)
				action.WriteByte(REMOVE)
				binary.Write(action,binary.LittleEndian,ps.ID)
				actions.Push(action.Bytes())
				ps.canSpawn = true
				ps.state = WAITFORSPAWN

			}
		}
		actionSize := actions.Size()
		positionSize := int16(len(clients))

		effectLock.Lock()
		time := time.Now()
		var finished []int
		for i, effect := range playerEffects{
			if effect.onTick(time) {
				effect.onEnd()
				finished = append(finished,i)
			}
		}
		for i := len(finished) - 1; i >= 0; i--{
			playerEffects[finished[i]], playerEffects = playerEffects[len(playerEffects)-1],playerEffects[:len(playerEffects)-1]
		}
		effectLock.Unlock()
		if actionSize > 0 || positionSize > 0 {

			message := new(bytes.Buffer)
			binary.Write(message, binary.LittleEndian, int16(0))
			binary.Write(message, binary.LittleEndian, positionSize)
			binary.Write(message, binary.LittleEndian, int16(actionSize))
			stateLock.Lock()

			for _, pt := range clients {
				binary.Write(message, binary.LittleEndian, pt.agentID)
				binary.Write(message, binary.LittleEndian, pt.transform.playerPosition.x)
				binary.Write(message, binary.LittleEndian, pt.transform.playerPosition.y)
				binary.Write(message, binary.LittleEndian, pt.transform.playerPosition.z)
				binary.Write(message, binary.LittleEndian, pt.transform.playerRotation.x)
				binary.Write(message, binary.LittleEndian, pt.transform.playerRotation.y)
				binary.Write(message, binary.LittleEndian, pt.transform.playerRotation.z)
				binary.Write(message, binary.LittleEndian, pt.transform.locationTime)
			}

			for i := 0; i < actionSize; i++ {
				message.Write(actions.Pop())
			}

			stateLock.Unlock()
			byteMessage := message.Bytes()
			message.Reset()
			binary.Write(message, binary.LittleEndian, int16(len(byteMessage)-2))
			for _, connection := range clients {
				connection.client.Write(byteMessage)
			}
		}
	}

}

func listenFunc(){
	listener,err := net.Listen("tcp",":7423")
	if err != nil{
		log.Println(err)
	}
	for{
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}else{
			go handleClient(conn)
		}

	}
}

func main(){
	log.Println("Started server")
	ticker := time.NewTicker(time.Minute * 10)
	clients = make(map[int16]*playerInfo)
	playerStates = make(map[int16]*playerTransform)
	shipStates = make(map[int16]*shipState)
	clientStats = make(map[int16]playerStats)
	nextSpawn = make(chan int)
	actions = NewQueue()
	abilityList = loadAbilityData()
	itemList = populateLoot()
	go listenFunc()
	go processState()
	go spawnLoc()
	for{
		_ = <-ticker.C
		log.Println("Server is up")
	}
}