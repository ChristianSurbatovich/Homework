package main


const(
	CURRENT_HEALTH int16 = 1101
	MAX_HEALTH int16 = 101
	SPEED int16 = 102
	CURRENT_SPEED int16 = 1102
	MOBILITY int16 = 103
	COOLDOWN int16 = 104
	DAMAGE_REDUCTION int16 = 105
	VISION_RANGE int16 = 106
	CANNON_DAMAGE int16 = 107
	HARPOON_DAMAGE int16 = 108
	LEVEL int16 = 109
	EXP int16 = 110
	CANNON_COOLDOWN int16 = 111
	HARPOON_COOLDOWN int16 = 112
	DAMAGE_INCREASE int16 = 113
	ACCELERATION int16 = 114
	TURN_SPEED int16 = 115
	SLOW_SPEED int16 = 116
	WIND_DRAG int16 = 117
	WEAPON_COOLDOWN int16 = 118
	WEAPON_POWER int16 = 119
	WEAPON_DAMAGE int16 = 120
	WEAPON_VERTICAL int16 = 121
	WEAPON_HORIZONTAL int16 = 122
	WEAPON_RANGE int16 = 123
	WEAPON_GROUPING int16 = 124
	WEAPON_SPREAD int16 = 125
	WEAPON_ACTIVE int16 = 126
)

type statMod struct{
	statID int16
	value float32
	percent bool
	resource bool
	resourceMode int16
}

const(
	MAX = iota
	MIN = iota
	PERCENT = iota
	ABSOLUTE = iota
	NONE = iota
)

func initializePlayerBaseStats() map[int16]float32{
	stats := make(map[int16]float32)
	stats[MAX_HEALTH] = 150
	stats[CURRENT_HEALTH] = 150
	stats[CANNON_DAMAGE] = 10
	stats[HARPOON_DAMAGE] = 5
	stats[SPEED] = 22
	stats[ACCELERATION] = 5
	stats[MOBILITY] = 15
	return stats
}

func initializePlayerFlatMod() map[int16]float32{
	stats := make(map[int16]float32)
	stats[MAX_HEALTH] = 0
	stats[CURRENT_HEALTH] = 0
	stats[CANNON_DAMAGE] = 0
	stats[HARPOON_DAMAGE] = 0
	stats[SPEED] = 0
	stats[ACCELERATION] = 0
	stats[MOBILITY] = 0
	return stats
}

func initializePlayerMultMod() map[int16]float32{
	stats := make(map[int16]float32)
	stats[MAX_HEALTH] = 1
	stats[CURRENT_HEALTH] = 1
	stats[CANNON_DAMAGE] = 1
	stats[HARPOON_DAMAGE] = 1
	stats[SPEED] = 1
	stats[ACCELERATION] = 1
	stats[MOBILITY] = 1
	return stats
}


type playerStats struct{
	baseStats map[int16]float32
	flatMods map[int16]float32
	multMods map[int16]float32
	totalStats map[int16]float32
}
