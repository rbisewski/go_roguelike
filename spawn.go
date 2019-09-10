/*
 * File: spawn.go
 *
 * Description: File to handle spawning new monsters into a given area.
 */

package main

import "fmt"

//! Function to spawn a creature in a given area.
/*
 * @param     string    name of the creature to add
 * @param     int       x-coord as int
 * @param     int       y-coord as int
 * @param     Area*     pointer to the intended area
 *
 * @return    bool      whether or not the creature was added
 */
func spawnCreatureToArray(name string, x int, y int, a *Area) bool {

	if len(name) < 1 || x < 0 || y < 0 || a == nil {
		DebugLog(&G, fmt.Sprintf("spawnCreatureToArray() --> invalid input"))
		return false
	}

	if !GlobalCreatureTypeInfoMapIsPopulated {
		return false
	}

	_, IsCreatureTypeDefined := GlobalCreatureTypeInfoMap[name]
	if !IsCreatureTypeDefined {
		DebugLog(&G, fmt.Sprintf("spawnCreatureToArray() --> improper "+
			"monster string given: %s", name))
		return false
	}

	// Grab the creature's name, species, rune-graphic, health, max-health,
	// attack, and defence attributes from the global creature type map.
	SpawnedCreatureName := GlobalCreatureTypeInfoMap[name].Name
	SpawnedCreatureSpecies := GlobalCreatureTypeInfoMap[name].Species
	SpawnedCreatureGfx := GlobalCreatureTypeInfoMap[name].Ch
	SpawnedCreatureHp := GlobalCreatureTypeInfoMap[name].Hp
	SpawnedCreatureMaxHp := GlobalCreatureTypeInfoMap[name].MaxHp
	SpawnedCreatureAttack := GlobalCreatureTypeInfoMap[name].Att
	SpawnedCreatureDefence := GlobalCreatureTypeInfoMap[name].Def
	SpawnedCreatureClass := GlobalCreatureTypeInfoMap[name].Class
	SpawnedCreatureStrength := GlobalCreatureTypeInfoMap[name].Strength
	SpawnedCreatureIntelligence := GlobalCreatureTypeInfoMap[name].Intelligence
	SpawnedCreatureAgility := GlobalCreatureTypeInfoMap[name].Agility
	SpawnedCreatureWisdom := GlobalCreatureTypeInfoMap[name].Wisdom
	SpawnedCreatureHealrate := GlobalCreatureTypeInfoMap[name].Healrate
	SpawnedCreatureHealcounter := GlobalCreatureTypeInfoMap[name].Healcounter

	// Append it to the array.
	a.Creatures = append(a.Creatures, NewCreature(SpawnedCreatureName,
		SpawnedCreatureSpecies, y, x, SpawnedCreatureGfx, a, nil,
		SpawnedCreatureHp, SpawnedCreatureMaxHp, SpawnedCreatureAttack,
		SpawnedCreatureDefence, SpawnedCreatureClass,
		SpawnedCreatureStrength, SpawnedCreatureIntelligence,
		SpawnedCreatureAgility, SpawnedCreatureWisdom, SpawnedCreatureHealrate,
		SpawnedCreatureHealcounter))

	return true
}
