package main

const(
	HIT byte = 0x01
	HEALTH byte = 0x02
	SINK byte = 0x03
	FIRE byte = 0x04
	OPEN byte = 0x05
	NAME byte = 0x06
	SPAWN byte = 0x07
	REGISTER byte = 0x08
	RESPAWN byte = 0x09
	REMOVE byte = 0x0A
	DESPAWN byte = 0x0B
	CHAT byte = 0x0C
	FEED byte = 0x0D
	EXPLODE byte = 0x0E
	POSITION byte = 0x0F
	STAT byte = 0x10
	ABILITY byte = 0x11
	ADD_ACTION byte = 0x12
	EQUIP byte = 0x13
	ITEM byte = 0x14
	LOOT_AREA byte = 0x15
	LOOT_ITEM byte = 0x16
	ADD_LOOT_ITEM byte = 0x17
	PICKUP_ITEM byte = 0x18
	UNEQUIP byte = 0x19
	TEST byte = 0xFF

	HIT_LENGTH int16 = 18
	HEALTH_LENGTH int16 = 9
	SINK_LENGTH int16 = 29
	REGISTER_LENGTH int16 = 9
	SPAWN_LENGTH int16 = 40
)
