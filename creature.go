/*
 * File: creature.go
 *
 * Description: Class to handle player character and monsters.
 *
 * Future feature, add support for unicode char mobs like... 🦂
 */

package main

import (
	"fmt"
	"strconv"

	"github.com/rbisewski/go_roguelike/types"
)

// Creature ... structure to hold characters / monsters
type Creature struct {

	// Holds the name of the given creature.
	name string

	// Holds the type of the given creature.
	species string

	// Store the current (x,y) coord of the creature.
	Y int
	X int

	// Appearance of the creature.
	ch rune

	// Pointer to the level the creature is on.
	area *Area

	// Pointer to the inventory, which consists of objects of class Item.
	inventory []*Item

	// Stats attributes.
	Hp    int
	MaxHp int
	Att   int
	Def   int

	// creature class
	class *types.ClassTypeInfo

	// Four primary creature attributes
	Strength     uint
	Intelligence uint
	Agility      uint
	Wisdom       uint

	// The number of steps required to heal by 1 point of health.
	Healrate uint

	// The number of steps currently walked by the creature in question.
	Healcounter uint

	// Pointer to the creature equipment locations.
	*equipment
}

// Structure to hold the equipment being utilized by certain creatures.
type equipment struct {
	Head      *Item
	Neck      *Item
	Torso     *Item
	RightHand *Item
	LeftHand  *Item
	Pants     *Item
}

//! Creature Equipment Constructor
/*
 * @param    *Item         item equipped in the head location
 * @param    *Item         item equipped in the neck location
 * @param    *Item         item equipped in the torso location
 * @param    *Item         item equipped in the right hand location
 * @param    *Item         item equipped in the left hand location
 * @param    *Item         item equipped in the pants location
 *
 * @return   equipment*    pointer to a newly allocated equipment object
 */
func newEquipment(head, neck, torso, righthand, lefthand,
	pants *Item) *equipment {

	// Return the address of the new equipment object.
	return &equipment{head, neck, torso, righthand, lefthand, pants}
}

// NewCreature ... create a creature w/o Equipment Constructor
/*
 * @param     string       creature name
 * @param     string       creature species (i.e type)
 * @param     int          y-value
 * @param     int          x-value
 * @param     rune         ASCII character graphic
 * @param     Area*        pointer to an Area object
 * @param     int          current hit points
 * @param     int          maximum hit points
 * @param     int          attack
 * @param     int          defence
 * @param     uint         heal rate
 * @param     uint         heal counter
 *
 * @return    Creature*    pointer to a Creature w/ Stats
 */
func NewCreature(name string,
	species string,
	y int,
	x int,
	ch rune,
	area *Area,
	inventory []*Item,
	hp int,
	max int,
	att int,
	def int,
	class *types.ClassTypeInfo,
	stre uint,
	intel uint,
	agl uint,
	wis uint,
	hr uint,
	hc uint) *Creature {

	// Assign memory for a creature object and return the address.
	return &Creature{name,
		species,
		y,
		x,
		ch,
		area,
		inventory,
		hp,
		max,
		att,
		def,
		class,
		stre,
		intel,
		agl,
		wis,
		hr,
		hc,
		nil}
}

// NewCreatureWithEquipment ... creature w/ Equipment Constructor
/*
 * @param     string       creature name
 * @param     string       creature species (i.e type)
 * @param     int          y-value
 * @param     int          x-value
 * @param     rune         ASCII character graphic
 * @param     Area*        pointer to an Area object
 * @param     int          current hit points
 * @param     int          maximum hit points
 * @param     int          attack
 * @param     int          defence
 * @param     uint         heal rate
 * @param     uint         heal counter
 *
 * @return    Creature*    pointer to a Creature w/ Stats
 */
func NewCreatureWithEquipment(name string,
	species string,
	y int,
	x int,
	ch rune,
	area *Area,
	inventory []*Item,
	hp int,
	max int,
	att int,
	def int,
	class *types.ClassTypeInfo,
	stre uint,
	intel uint,
	agl uint,
	wis uint,
	hr uint,
	hc uint) *Creature {

	// Assign memory for a creature object and return the address.
	return &Creature{name,
		species,
		y,
		x,
		ch,
		area,
		inventory,
		hp,
		max,
		att,
		def,
		class,
		stre,
		intel,
		agl,
		wis,
		hr,
		hc,
		newEquipment(nil, nil, nil, nil, nil, nil)}
}

// Move ... translates a mob to a new (x,y) location.
/*
 * @param     int    y-value
 * @param     int    x-value
 *
 * @return    none
 */
func (m *Creature) Move(y, x int) {

	if x > 32767 || y > 32767 || x < -32767 || y < -32767 {
		DebugLog(&G, "Error: Invalid (x,y) coord detected.")
		return
	}

	if m == nil {
		DebugLog(&G, "Error: Null pointer or invalid creature detected.")
		return
	}

	// Since this has a creature and it appears to have valid coords, then
	// go ahead and test it again the tile the creature in question wishes
	// to move to.
	tileRune, blocks, hasCreature, hasItems := m.area.GetTileInfo(m.Y+y, m.X+x)

	if tileRune == 0 {
		DebugLog(&G, "Error: Null or invalid Unicode rune.")
		return
	}

	// Increment the monster's healing counter.
	m.Healcounter++

	// If the heal counter has surpassed the healing rate value, then...
	if m.Healcounter >= m.Healrate {

		// Set the counter back to zero
		m.Healcounter = 0
	}

	// If the healing counter is zero and creature is not fully healed...
	if m.Healcounter == 0 && m.MaxHp > m.Hp {

		// Increase the current hitpoints of the creature by 1
		m.Hp++
	}

	// If the player attempts to move to a blocking tile, and it is a wall,
	// go ahead and print a short message and then leave function.
	if blocks && m.species == "player" && tileRune == '#' {
		MessageLog.log("The wall is solid and damp, and you cannot move past.")
		return
	}

	// Catch-all message for when the player moves into a blocking tile.
	if blocks && m.species == "player" {
		MessageLog.log("Something here is blocking, and you cannot move past.")
		return
	}

	// If some other creature attempts to move, simply return here since
	// there is no need to print a message, except in debug mode.
	if blocks && m.species != "player" {
		DebugLog(&G, fmt.Sprintf(
			"The %s attempted to move to location (%d,%d), but it "+
				"was blocked.",
			m.name,
			m.Y+y,
			m.X+x))
		return
	}

	// If the tile is non-blocking, but a creature is here, go ahead and
	// switch to combat mode via the attack() function.
	if hasCreature != nil && m != hasCreature {
		DebugLog(&G, fmt.Sprintf(
			"The %s is attacking %s at location (%d,%d).",
			m.name,
			hasCreature.name,
			m.Y+y,
			m.X+x))

		m.attack(hasCreature)
		return
	}

	// If debug mode, tell the developer where the creature has moved to.
	DebugLog(&G, fmt.Sprintf(
		"The %s moved to location (%d,%d).",
		m.name,
		m.Y+y,
		m.X+x))

	// Since the tile is non-blocking, and no creature is present, then
	// go ahead and move there.
	m.Y += y
	m.X += x

	// If there are items laying on the ground, give the player some
	// indicator of what is there.
	if m.species == "player" && len(hasItems) == 1 {
		MessageLog.log(fmt.Sprintf("On the ground lies a %s.",
			hasItems[0].category))

		// Else if the player has moved to a tile that contains more than 1 item,
		// print the following message.
	} else if m.species == "player" && len(hasItems) > 1 {
		MessageLog.log("There are items here on the ground.")
	}
}

//! Function to handle what occurs if a monster attacks.
/*
 * @param     Creature*    defending creature / PC
 *
 * @return    none
 */
func (m *Creature) attack(defender *Creature) {

	if m == nil || defender == nil {
		DebugLog(&G, "attack() --> invalid creature input.")
		return
	}

	var damageDealt int = m.Att - defender.Def

	// Cap the damage dealt at zero, this is to prevent the enemies from
	// accidently healing other creatures when they attack.
	if damageDealt < 0 {
		damageDealt = 0
	}

	// Adjust the defender's HP based on the damage dealt.
	defender.Hp -= damageDealt
	if defender.Hp <= 0 {
		defender.die()
	}

	// If two monsters are attacking each other, there is no need to
	// print battle messages, ergo this function is complete.
	if defender.species != "player" && m.species != "player" {
		return
	}

	// Print a message telling the end-user they have been injured
	// during the attack.
	if defender.species == "player" {
		MessageLog.log(fmt.Sprintf("The %s injures you for %d hit points.",
			m.name, damageDealt))
		return
	}

	// Otherwise the player is doing the attack, so explain how much damage
	// was done to the creature being attacked.
	MessageLog.log(fmt.Sprintf("You strike the %s for %d hit points "+
		"of damage.", defender.name, damageDealt))

	// If creature being attacked has reached zero hit points, go ahead and
	// print a message stating that the creature has died.
	if defender.Hp < 1 {
		MessageLog.log(fmt.Sprintf("The %s has died.", defender.name))
		return
	}

	// Otherwise if the creature has not yet died, go ahead and give a
	// description of the current state of the attacked creature in the
	// lower-left message screen.
	//
	// The adjectives to be used are as follows:
	//
	// 100% --> unscathed
	//  75% --> slightly injured
	//  50% --> injured
	//  25% --> severely injured
	//
	if defender.Hp == defender.MaxHp {
		MessageLog.log(fmt.Sprintf("The %s still looks unscathed.",
			defender.name))

	} else if defender.Hp > int(float64(defender.MaxHp)*0.50) {
		MessageLog.log(fmt.Sprintf("The %s looks slightly injured.",
			defender.name))

	} else if defender.Hp > int(float64(defender.MaxHp)*0.25) {
		MessageLog.log(fmt.Sprintf("The %s looks injured.",
			defender.name))
	} else {
		MessageLog.log(fmt.Sprintf("The %s looks severely injured.",
			defender.name))
	}
}

//! Function to handle what occurs when a monster dies.
/*
 * @param     Creature*    monster who is currently dying
 *
 * @return    none
 */
func (m *Creature) die() {

	if m == nil {
		DebugLog(&G, "die() --> invalid creature input.")
		return
	}

	// If the "monster" who died is the player, then call that routine.
	if m == G.Player {

		// Call the player death functionality, which should end the current
		// instance of the game.
		G.Death()

		// Leave the function.
		return
	}

	// Adjust the array of monsters to account for the newly dead monster.
	for i, monster := range m.area.Creatures {

		// Sanity check, make sure this actually got a valid monster.
		if monster == nil {

			// Otherwise tell the developer something odd was appended here.
			numStr := strconv.Itoa(i)
			DebugLog(&G, "die() --> nil mob at index ["+numStr+"]")

			// Move on to the next monster.
			continue
		}

		// Ignore other monsters since they are still alive.
		if m != monster {
			continue
		}

		// If the monster who died is still in this level, then it needs
		// to be deleted from the array.
		m.area.Creatures = append(m.area.Creatures[:i], m.area.Creatures[i+1:]...)
	}

	// Check if the monster in question has no inventory, or an empty
	// inventory, then this is done.
	if m.inventory == nil || len(m.inventory) < 1 {

		// Variable declaration.
		var corpse *Item

		// Create an item that consists of the monster corpse.
		corpse = NewItem(fmt.Sprintf("corpse of %s", m.name), "corpse",
			m.Y, m.X, '%', m.area, false, false, 0, 0, 0, 0, 10, 0, 0)

		// Leave a creature corpse item in the shape of a % at the given
		// vertex (x,y) location of the formerly alive monster.
		m.area.Items = append(m.area.Items, corpse)

		// All is now done.
		return
	}

	// Otherwise, attempt to add all of the inventory items to the ground.
	for i, item := range m.inventory {

		// Sanity check, make sure this actually got a valid item.
		if item == nil {
			numStr := strconv.Itoa(i)
			DebugLog(&G, "die() --> nil item at index ["+numStr+"]")
			continue
		}

		// Give the item a % rune shape, and set it to the (x,y) coord of
		// dead monster.
		item.ch = '%'
		item.X = m.X
		item.Y = m.Y
		m.area.Items = append(m.area.Items, item)
	}
}
