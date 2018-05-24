package main

import (
	"math/rand"
)

type lootArea struct{
	id int16
	tier int16
	randomSeed float32
	looted bool
	lootList map[int16]baseItem
}


func (area *lootArea) generateLoot(itemList map[int16][]baseItem, IDcounter *int16){
	area.looted = true
	area.lootList = make(map[int16]baseItem)
	listLength := len(itemList[area.tier])
	numItems := rand.Intn(4) + 1
	for i := 0; i < numItems; i++{
		area.lootList[*IDcounter] = itemList[area.tier][rand.Intn(listLength)]
		*IDcounter++
	}

}

func populateAreas()map[int16]lootArea{
	list := make(map[int16]lootArea)
	list[2000] = lootArea{2000,1,1,false,nil}
	list[2001] = lootArea{2001,1,1,false,nil}
	list[2002] = lootArea{2002,1,1,false,nil}
	list[2003] = lootArea{2003,2,1,false,nil}
	list[2004] = lootArea{2004,2,1,false,nil}
	list[2005] = lootArea{2005,2,1,false,nil}
	list[2006] = lootArea{2006,2,1,false,nil}
	list[2007] = lootArea{2007,3,1,false,nil}
	list[2008] = lootArea{2008,3,1,false,nil}
	list[2009] = lootArea{2009,3,1,false,nil}
	return list
}

func populateLoot() map[int16][]baseItem{
	list := make(map[int16][]baseItem)
	list[1] = make([]baseItem,3)
	list[1][0] = &genericStatItem{itemID:1100,statMods:[]statMod{{statID:CANNON_DAMAGE,value:0.5,percent:true,resource:false}}}
	list[1][1] = &genericStatItem{itemID:1101,statMods:[]statMod{{statID:MAX_HEALTH,value:50,percent:false,resource:true,resourceMode:PERCENT}}}
	list[1][2] = &genericStatItem{itemID:1102,statMods:[]statMod{{statID:SPEED,value:7,percent:false,resource:true,resourceMode:PERCENT}}}
	list[2] = make([]baseItem,3)
	list[2][0] = &genericStatItem{itemID:1100,statMods:[]statMod{{statID:CANNON_DAMAGE,value:0.75,percent:true,resource:false}}}
	list[2][1] = &genericStatItem{itemID:1101,statMods:[]statMod{{statID:MAX_HEALTH,value:100,percent:false,resource:true,resourceMode:PERCENT}}}
	list[2][2] = &genericStatItem{itemID:1102,statMods:[]statMod{{statID:SPEED,value:11,percent:false,resource:true,resourceMode:PERCENT}}}
	list[3] = make([]baseItem,3)
	list[3][0] = &genericStatItem{itemID:1100,statMods:[]statMod{{statID:CANNON_DAMAGE,value:1.0,percent:true,resource:false}}}
	list[3][1] = &genericStatItem{itemID:1101,statMods:[]statMod{{statID:MAX_HEALTH,value:150,percent:false,resource:true,resourceMode:PERCENT}}}
	list[3][2] = &genericStatItem{itemID:1102,statMods:[]statMod{{statID:SPEED,value:15,percent:false,resource:true,resourceMode:PERCENT}}}
	return list
}