package main

import (
	"encoding/binary"
	"bytes"
)

type baseItem interface{
	id() int16
	onLoot(playerID int16)
	onPickup(playerID int16)
	onEquip(playerID int16)
	onUnequip(playerID int16)
	onDestroy(playerID int16)
}

type genericStatItem struct{
	itemID int16
	statMods []statMod
}


func (item *genericStatItem) id() int16{
	return item.itemID
}


func (item *genericStatItem) onLoot(playerID int16){

}

func (item *genericStatItem) onPickup(playerID int16){

}

func (item *genericStatItem) onEquip(playerID int16){
	tempStat := clientStats[playerID]
	for _,mod := range item.statMods{
		if mod.percent {
			tempStat.multMods[mod.statID] += mod.value

		}else{
			tempStat.flatMods[mod.statID] += mod.value
		}
		var percentOfTotal float32
		if mod.resource {
			percentOfTotal = tempStat.baseStats[mod.statID + 1000] / tempStat.totalStats[mod.statID]
			tempStat.totalStats[mod.statID] = (tempStat.baseStats[mod.statID] + tempStat.flatMods[mod.statID]) * tempStat.multMods[mod.statID]
			switch mod.resourceMode{
			case MAX:
				tempStat.baseStats[mod.statID + 1000] = tempStat.totalStats[mod.statID]
			case MIN:
				tempStat.baseStats[mod.statID + 1000] = tempStat.baseStats[mod.statID]
			case PERCENT:
				tempStat.baseStats[mod.statID + 1000] = percentOfTotal * tempStat.totalStats[mod.statID]
			case NONE:
			}
			messageBuffer := new(bytes.Buffer)
			messageBuffer.WriteByte(STAT)
			binary.Write(messageBuffer,binary.LittleEndian,playerID)
			binary.Write(messageBuffer,binary.LittleEndian,mod.statID + int16(1000))
			binary.Write(messageBuffer,binary.LittleEndian,clientStats[playerID].baseStats[mod.statID + 1000])
			actions.Push(messageBuffer.Bytes())
		}else {
			tempStat.totalStats[mod.statID] = (tempStat.baseStats[mod.statID] + tempStat.flatMods[mod.statID]) * tempStat.multMods[mod.statID]
		}
		messageBuffer := new(bytes.Buffer)
		messageBuffer.WriteByte(STAT)
		binary.Write(messageBuffer,binary.LittleEndian,playerID)
		binary.Write(messageBuffer,binary.LittleEndian,mod.statID)
		binary.Write(messageBuffer,binary.LittleEndian,clientStats[playerID].totalStats[mod.statID])
		actions.Push(messageBuffer.Bytes())
	}
}

func (item *genericStatItem) onUnequip(playerID int16){
	tempStat := clientStats[playerID]
	for _,mod := range item.statMods{
		if mod.percent {
			tempStat.multMods[mod.statID] -= mod.value

		}else{
			tempStat.flatMods[mod.statID] -= mod.value
		}
		var percentOfTotal float32
		if mod.resource {
			percentOfTotal = tempStat.baseStats[mod.statID + 1000] / tempStat.totalStats[mod.statID]
			tempStat.totalStats[mod.statID] = (tempStat.baseStats[mod.statID] + tempStat.flatMods[mod.statID]) * tempStat.multMods[mod.statID]
			switch mod.resourceMode{
			case MAX:
				tempStat.baseStats[mod.statID + 1000] = tempStat.totalStats[mod.statID]
			case MIN:
				tempStat.baseStats[mod.statID + 1000] = tempStat.baseStats[mod.statID]
			case PERCENT:
				tempStat.baseStats[mod.statID + 1000] = percentOfTotal * tempStat.totalStats[mod.statID]
			case NONE:
			}
			messageBuffer := new(bytes.Buffer)
			messageBuffer.WriteByte(STAT)
			binary.Write(messageBuffer,binary.LittleEndian,playerID)
			binary.Write(messageBuffer,binary.LittleEndian,mod.statID + int16(1000))
			binary.Write(messageBuffer,binary.LittleEndian,clientStats[playerID].baseStats[mod.statID + 1000])
			actions.Push(messageBuffer.Bytes())
		}else {
			tempStat.totalStats[mod.statID] = (tempStat.baseStats[mod.statID] + tempStat.flatMods[mod.statID]) * tempStat.multMods[mod.statID]
		}
		messageBuffer := new(bytes.Buffer)
		messageBuffer.WriteByte(STAT)
		binary.Write(messageBuffer,binary.LittleEndian,playerID)
		binary.Write(messageBuffer,binary.LittleEndian,mod.statID)
		binary.Write(messageBuffer,binary.LittleEndian,clientStats[playerID].totalStats[mod.statID])
		actions.Push(messageBuffer.Bytes())
	}
}

func (item *genericStatItem) onDestroy(playerID int16){

}
