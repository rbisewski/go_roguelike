/*
 * File: item.go
 *
 * Description: Class to handle the various small items that can be picked
 *              up or equipped.
 */

package main

import "fmt"

// Item ... Structure to hold the attributes of an item.
type Item struct {

	// Holds the name of the given item.
	name string

	// Holds the type of the given item.
	category string

	// Store the current (x,y) coord of the item.
	Y int
	X int

	// Appearance of the item.
	ch rune

	// Pointer to the level of where the item is.
	area *Area

	// Determines whether an item can be equipped by a creature.
	canEquip bool

	// Variable that holds whether or not an item is broken.
	isBroken bool

	// The current condition of an item. If the max durability and the
	// current are equal, then the item is in perfect condition.
	durabilityCurrent int
	durabilityMaximum int

	// Price of the item in gold coins.
	priceToPurchase int
	priceToSell     int

	// Weight of the item in grams.
	weight int

	// Variables that determine how much an item increases (if positive)
	// or decreases (if negative) the attack / defence of a creature.
	attackIncrease  int
	defenceIncrease int
}

// NewItem ... Item constructor function.
/*
 * @param     string    item name
 * @param     string    category
 * @param     int       Y
 * @param     int       X
 * @param     rune      item appearance as Unicode rune
 * @param     *Area     pointer to area object that the item is located,
 *                      if this is `nil` then the item is held by a creature
 * @param     bool      whether or not the item can be equipped
 * @param     bool      whether or not the item is broken
 * @param     int       current durability
 * @param     int       maximum durability
 * @param     int       purchase price
 * @param     int       sell price
 * @param     int       weight
 * @param     int       attack increase (if > 0) or decrease (if < 0)
 * @param     int       defence increase (if > 0) or decrease (if < 0)
 *
 * @return    Item*    pointer to a newly initialized item.
 */
func NewItem(name string,
	category string,
	Y int,
	X int,
	ch rune,
	area *Area,
	canEquip bool,
	isBroken bool,
	durabilityCurrent int,
	durabilityMaximum int,
	priceToPurchase int,
	priceToSell int,
	weight int,
	attackIncrease int,
	defenceIncrease int) *Item {

	// Return an address to a newly allocated Creature object.
	return &Item{name,
		category,
		Y,
		X,
		ch,
		area,
		canEquip,
		isBroken,
		durabilityCurrent,
		durabilityMaximum,
		priceToPurchase,
		priceToSell,
		weight,
		attackIncrease,
		defenceIncrease}
}

//! Function to handle what occurs when the current durability of an item
//! changes.
//!
//! This function is meant to deal with situations where a given
//! item is either being worn-out, so the current durability decreases or in
//! the process of being repaired, in which case the current durability
//! increases.
//!
//! Note that an item can only be repaired at most to the defined maximum
//! durability value. On the other hand, if an item reaches zero, then the
//! item should become broken and unusable.
/*
 * @caller    Item*    item to adjust durability thereof
 *
 * @param     int      amount to adjust the current durability by
 *
 * @return    none
 */
func (itm *Item) adjustDurability(amount int) {

	// If the durability was zero, then nothing to do.
	if amount == 0 {
		DebugLog(&G, fmt.Sprintf("adjustDurability() --> was given amount "+
			"of value zero, so nothing to be done..."))
		return
	}

	// If debug, state how much the item durability currently is.
	DebugLog(&G, fmt.Sprintf("adjustDurability() --> Item [%s] durability"+
		"before adjustment is: %d / %d",
		itm.name,
		itm.durabilityCurrent,
		itm.durabilityMaximum))

	// Adjust the current durability by the amount specified.
	itm.durabilityCurrent += amount

	// If debug, state how much the item durability has changed by.
	DebugLog(&G, fmt.Sprintf("adjustDurability() --> item [%s] durability"+
		"after adjustment is: %d / %d",
		itm.name,
		itm.durabilityCurrent,
		itm.durabilityMaximum))

	// If the current item durability has reached beyond the max, go ahead
	// and cap it to the maximum durability.
	if itm.durabilityCurrent > itm.durabilityMaximum {

		// Cap the durability of the item.
		itm.durabilityCurrent = itm.durabilityMaximum

		// If debug mode, tell the developer what just happened.
		DebugLog(&G, "adjustDurability() --> item ["+itm.name+"] has "+
			"exceeded max durability, and so has "+
			"been capped")
	}

	// If the item is less than zero, then consider it to be broken, in
	// which event t
	if itm.durabilityCurrent < 1 {

		// If debug mode, tell the developer this is going to break the item.
		DebugLog(&G, fmt.Sprintf("adjustDurability() --> item [%s] is "+
			"preparing to be broken", itm.name))

		// Set the relevant item properties via this function in order to
		// render the item broken.
		itm.eventBroken()

		// Leave this routine, since the item broken event has been handled.
		return
	}
}

//! Function to handle what occurs when an item breaks.
/*
 * @caller    Item*    the given item which will become broken
 *
 * @return    none
 */
func (itm *Item) eventBroken() {

	// Prevent the item from being equipped.
	itm.canEquip = false

	// Set the isBroken flag to true.
	itm.isBroken = true

	// If debug mode, tell the developer that the item is now broken.
	DebugLog(&G, fmt.Sprintf("eventBroken() --> item [%s] is now broken",
		itm.name))
}
