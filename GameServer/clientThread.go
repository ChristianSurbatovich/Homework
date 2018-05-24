package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"
	"bufio"
	"fmt"
	"io"
	"math/rand"
)

type player struct{
	conn net.Conn
	playerID int16
	accountID uint64
}


func (playerData *player) SetStat(statID int16, statValue float32){

}

func handleClient(conn net.Conn){
	defer conn.Close()
	var clientTransform playerTransform
	var clientID int16
	var messageLength int16
	var baseStats map[int16]float32
	var flatStatModifiers map[int16]float32
	var multStatModifiers map[int16]float32
	var totalStats map[int16]float32
	var inventoryItems map[int16]baseItem
	var equippedItems map[int16]baseItem
	var localIDs int16 = 0
	lootableAreas := populateAreas()
	lengthBuffer := make([]byte,2)
	baseStats = initializePlayerBaseStats()
	flatStatModifiers = initializePlayerFlatMod()
	multStatModifiers = initializePlayerMultMod()
	inventoryItems = make(map[int16]baseItem)
	equippedItems = make(map[int16]baseItem)
	messageBuffer := new(bytes.Buffer)
	messageBuffer.Grow(256)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	currentPlayerState := 0
	reader := bufio.NewReader(conn)
	fmt.Println("Received a connection from: " + conn.RemoteAddr().String())
	log.Println("Received a connection from: " + conn.RemoteAddr().String())


	elem, ok := clients[-1]
	if ok{
		elem.client = conn
		newMessage(REGISTER_LENGTH,0,1,messageBuffer)
		addMessage(REGISTER,0,shipStates[clientID],nil,messageBuffer)
		conn.Write(messageBuffer.Bytes())
	}else{
		IDlock.Lock()
		clientID = currentPlayerID
		currentPlayerID++
		IDlock.Unlock()
		sL := <-nextSpawn
		clientTransform = playerTransform{vector{spawnLocations[sL].x,spawnLocations[sL].y,spawnLocations[sL].z},vector{0,float32(rnd.Intn(360)),0},clientID,0}
		playerStates[clientID] = &clientTransform
		client := playerInfo{
			conn,
			&clientTransform,
			clientID,
		}
		totalStats = make(map[int16]float32)
		for key,value := range baseStats{
			totalStats[key] = (value + flatStatModifiers[key]) * multStatModifiers[key]
		}
		clientStats[clientID] = playerStats{baseStats,flatStatModifiers,multStatModifiers,totalStats}
		clients[clientID] = &client
		shipStates[clientID] = &shipState{clientID,false,[]bool{true},"",ALIVE,time.Now(),true,true}

		newMessage(REGISTER_LENGTH,0,1,messageBuffer)
		addMessage(REGISTER,0, shipStates[clientID],nil,messageBuffer)
		conn.Write(messageBuffer.Bytes())
	}

	actions.Push(createAction(SPAWN,0,shipStates[clientID],&clientTransform))
	messageBuffer.Reset()
	newMessage(0,0,int16(len(clients) - 1),messageBuffer)
	stateLock.Lock()
	for _, pt := range clients{
		if pt.agentID == clientID || shipStates[pt.agentID].alive == false{
			continue
		}
		addMessage(SPAWN,0,shipStates[pt.agentID],pt.transform,messageBuffer)
	}
	stateLock.Unlock()
	byteMessage := messageBuffer.Bytes()
	messageBuffer.Reset()
	binary.Write(messageBuffer,binary.LittleEndian,int16(len(byteMessage) - 2))
	conn.Write(byteMessage)
	for{
		n, err := io.ReadFull(reader,lengthBuffer)

		if err != nil {
			log.Printf("Bytes read: %d\n", n)
			log.Println(err)
			break
		}
		binary.Read(bytes.NewReader(lengthBuffer),binary.LittleEndian,&messageLength)
		message := make([]byte,messageLength)
		n, err = reader.Read(message)
		if err != nil{
			log.Println(err)
			if err == io.EOF{
				continue
			}else{
				break
			}
		}
		if len(message) > 1{
			switch message[0]{
			case POSITION:
				if currentPlayerState == 1{
					continue
				}
				messageReader := bytes.NewReader(message[3:])
				stateLock.RLock()
				binary.Read(messageReader,binary.LittleEndian,&clientTransform.playerPosition.x)
				binary.Read(messageReader,binary.LittleEndian,&clientTransform.playerPosition.y)
				binary.Read(messageReader,binary.LittleEndian,&clientTransform.playerPosition.z)
				binary.Read(messageReader,binary.LittleEndian,&clientTransform.playerRotation.x)
				binary.Read(messageReader,binary.LittleEndian,&clientTransform.playerRotation.y)
				binary.Read(messageReader,binary.LittleEndian,&clientTransform.playerRotation.z)
				binary.Read(messageReader,binary.LittleEndian,&clientTransform.locationTime)
				stateLock.RUnlock()
				// fire id numshots x1 y1 z1 x2 y2 z2 ...
			case FIRE:
				actions.Push(message)
			case OPEN:
				actions.Push(message)
				shipStates[clientID].doorsOpen = !shipStates[clientID].doorsOpen
			case HIT:
				var playerHit bool
				var id int16
				var weaponID int16
				messageReader := bytes.NewReader(message[1:])
				binary.Read(messageReader,binary.LittleEndian,&playerHit)
				binary.Read(messageReader,binary.LittleEndian,&id)
				binary.Read(messageReader,binary.LittleEndian,&weaponID)
				if playerHit {
					temp := shipStates[id]
					if temp != nil {
						clientStats[id].baseStats[CURRENT_HEALTH] -= clientStats[clientID].totalStats[weaponID]
						if clientStats[id].baseStats[CURRENT_HEALTH] <= 0 && shipStates[id].alive {
							shipStates[id].alive = false
							action := new(bytes.Buffer)
							action.WriteByte(FEED)
							binary.Write(action, binary.LittleEndian, clientID)
							binary.Write(action, binary.LittleEndian, id)
							actions.Push(action.Bytes())
						}
					}
					if id == clientID && clientStats[id].baseStats[CURRENT_HEALTH] <= 0 {
						currentPlayerState = 1
					}
					action := new(bytes.Buffer)
					action.WriteByte(HEALTH)
					binary.Write(action,binary.LittleEndian,id)
					binary.Write(action,binary.LittleEndian,clientStats[id].baseStats[CURRENT_HEALTH])
					actions.Push(action.Bytes())
				}
				actions.Push(message)
			case NAME:
				var id, nameLength int16
				messageReader := bytes.NewReader(message[1:])
				binary.Read(messageReader,binary.LittleEndian,&id)
				binary.Read(messageReader,binary.LittleEndian,&nameLength)
				shipStates[id].name = string(message[5:5+nameLength])
				actions.Push(message)
			case SPAWN:
				if shipStates[clientID].canSpawn{
					sL := <-nextSpawn
					clientTransform.playerPosition = vector{spawnLocations[sL].x,spawnLocations[sL].y,spawnLocations[sL].z}
					clientTransform.playerRotation = vector{0,float32(rnd.Intn(360)),0}
					clientStats[clientID].baseStats[CURRENT_HEALTH] = clientStats[clientID].totalStats[MAX_HEALTH]
					shipStates[clientID].doorsOpen = false
					shipStates[clientID].state = ALIVE
					shipStates[clientID].alive = true
					actions.Push(createAction(SPAWN,0,shipStates[clientID],&clientTransform))
				}
			case CHAT:
				var chatLength, nameLength int16
				nameLength = int16(len(shipStates[clientID].name))
				chatLength = int16(len(message) - 3)
				action := new(bytes.Buffer)
				action.WriteByte(CHAT)
				binary.Write(action,binary.LittleEndian,nameLength)
				binary.Write(action,binary.LittleEndian,chatLength)
				action.Write([]byte(shipStates[clientID].name))
				action.Write(message[3:])
				actions.Push(action.Bytes())
			case ABILITY:
				var abilityID int16
				binary.Read(bytes.NewReader(message[1:]),binary.LittleEndian,&abilityID)
				newAbilityData := abilityList[abilityID]
				newAbility := NewStatBuff(newAbilityData,clientID)
				newAbility.onStart(time.Now())
				playerEffects = append(playerEffects,&newAbility)
			case LOOT_AREA:
				var areaID int16
				binary.Read(bytes.NewReader(message[1:]),binary.LittleEndian,&areaID)
				if !lootableAreas[areaID].looted{
					temp := lootableAreas[areaID]
					delete(lootableAreas, areaID)
					temp.generateLoot(itemList, &localIDs)
					lootableAreas[areaID] = temp
				}
				messageBuffer.Reset()
				binary.Write(messageBuffer,binary.LittleEndian,int16(len(lootableAreas[areaID].lootList) * 4 + 7))
				binary.Write(messageBuffer,binary.LittleEndian,int16(0))
				binary.Write(messageBuffer,binary.LittleEndian,int16(1))
				messageBuffer.WriteByte(ADD_LOOT_ITEM)
				binary.Write(messageBuffer,binary.LittleEndian,int16(len(lootableAreas[areaID].lootList)))
				for localID, item := range lootableAreas[areaID].lootList{
					binary.Write(messageBuffer,binary.LittleEndian,item.id())
					binary.Write(messageBuffer,binary.LittleEndian,localID)
				}
				conn.Write(messageBuffer.Bytes())
			case PICKUP_ITEM:
				var localID int16
				var areaID int16
				binary.Read(bytes.NewReader(message[1:]),binary.LittleEndian,&localID)
				binary.Read(bytes.NewReader(message[3:]),binary.LittleEndian,&areaID)
				item, exists := lootableAreas[areaID].lootList[localID]
				if exists{
					delete(lootableAreas[areaID].lootList,localID)
					inventoryItems[localID] = item
					messageBuffer.Reset()
					binary.Write(messageBuffer,binary.LittleEndian,int16(9))
					binary.Write(messageBuffer,binary.LittleEndian,int16(0))
					binary.Write(messageBuffer,binary.LittleEndian,int16(1))
					messageBuffer.WriteByte(LOOT_ITEM)
					binary.Write(messageBuffer,binary.LittleEndian,localID)
					binary.Write(messageBuffer,binary.LittleEndian,areaID)
					conn.Write(messageBuffer.Bytes())
				}
			case EQUIP:
				var localID int16
				binary.Read(bytes.NewReader(message[1:]),binary.LittleEndian,&localID)
				if item, exists := inventoryItems[localID]; exists{
					equippedItems[localID] = item
					delete(inventoryItems,localID)
					item.onEquip(clientID)
				}
			case UNEQUIP:
				var localID int16
				binary.Read(bytes.NewReader(message[1:]),binary.LittleEndian,&localID)
				if item, exists := equippedItems[localID]; exists{
					inventoryItems[localID] = item
					delete(equippedItems,localID)
					item.onUnequip(clientID)
				}
			}
		}
	}
	log.Printf("Player %s disconnected",shipStates[clientID].name)
	stateLock.Lock()
	delete(shipStates, clientID)
	delete(playerStates,clientID)
	delete(clients,clientID)
	action := new(bytes.Buffer)
	action.WriteByte(REMOVE)
	binary.Write(action,binary.LittleEndian,clientID)
	actions.Push(action.Bytes())
	stateLock.Unlock()
}