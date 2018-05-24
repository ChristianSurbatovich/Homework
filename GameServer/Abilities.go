package main

import (
	"time"
	"bytes"
	"encoding/binary"
)

type abilityData struct{
	abilityID int16
	statMods []statMod
	duration time.Duration
}


func loadAbilityData() map[int16]abilityData{
	abilityList := make(map[int16]abilityData)
	abilityList[300] = abilityData{abilityID:300,duration:10*time.Second,statMods:[]statMod{{statID:SPEED,value:0.5,percent:true,resource:true,resourceMode:MAX}}}

	return abilityList
}

type ability interface {
	onStart(time time.Time)
	onEnd()
	onTick(time time.Time) bool
}

type statBuff struct{
	abilityID int16
	playerID int16
	statMods []statMod
	startTime time.Time
	duration time.Duration
}

func NewStatBuff(data abilityData, player int16) statBuff{
	buff := statBuff{abilityID:data.abilityID,playerID:player,statMods:data.statMods,duration:data.duration}
	return buff
}

func (statBuff *statBuff) onStart(time time.Time){
	statBuff.startTime = time
	tempStat := clientStats[statBuff.playerID]
	for _,mod := range statBuff.statMods{
		if mod.percent {
			tempStat.multMods[mod.statID] += mod.value

		}else{
			tempStat.flatMods[mod.statID] += mod.value
		}
		var percentOfTotal float32
		if mod.resource {
			messageBuffer := new(bytes.Buffer)
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

			messageBuffer.WriteByte(STAT)
			binary.Write(messageBuffer,binary.LittleEndian,statBuff.playerID)
			binary.Write(messageBuffer,binary.LittleEndian,mod.statID + int16(1000))
			binary.Write(messageBuffer,binary.LittleEndian,clientStats[statBuff.playerID].baseStats[mod.statID + 1000])
			actions.Push(messageBuffer.Bytes())
		}else {
			tempStat.totalStats[mod.statID] = (tempStat.baseStats[mod.statID] + tempStat.flatMods[mod.statID]) * tempStat.multMods[mod.statID]
		}
		messageBuffer := new(bytes.Buffer)
		messageBuffer.WriteByte(STAT)
		binary.Write(messageBuffer,binary.LittleEndian,statBuff.playerID)
		binary.Write(messageBuffer,binary.LittleEndian,mod.statID)
		binary.Write(messageBuffer,binary.LittleEndian,clientStats[statBuff.playerID].totalStats[mod.statID])
		actions.Push(messageBuffer.Bytes())
	}
}

func (statBuff *statBuff ) onTick(time time.Time) bool{
	if time.Sub(statBuff.startTime) > statBuff.duration{
		return true
	}
	return false
}

func (statBuff *statBuff ) onEnd(){
	tempStat := clientStats[statBuff.playerID]
	for _,mod := range statBuff.statMods{
		if mod.percent {
			tempStat.multMods[mod.statID] -= mod.value

		}else{
			tempStat.flatMods[mod.statID] -= mod.value
		}
		var percentOfTotal float32
		if mod.resource {
			messageBuffer := new(bytes.Buffer)
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

			messageBuffer.WriteByte(STAT)
			binary.Write(messageBuffer,binary.LittleEndian,statBuff.playerID)
			binary.Write(messageBuffer,binary.LittleEndian,mod.statID + int16(1000))
			binary.Write(messageBuffer,binary.LittleEndian,clientStats[statBuff.playerID].baseStats[mod.statID + 1000])
			actions.Push(messageBuffer.Bytes())
		}else {
			tempStat.totalStats[mod.statID] = (tempStat.baseStats[mod.statID] + tempStat.flatMods[mod.statID]) * tempStat.multMods[mod.statID]
		}
		messageBuffer := new(bytes.Buffer)
		messageBuffer.WriteByte(STAT)
		binary.Write(messageBuffer,binary.LittleEndian,statBuff.playerID)
		binary.Write(messageBuffer,binary.LittleEndian,mod.statID)
		binary.Write(messageBuffer,binary.LittleEndian,clientStats[statBuff.playerID].totalStats[mod.statID])
		actions.Push(messageBuffer.Bytes())
	}
}
