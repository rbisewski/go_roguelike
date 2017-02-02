/*
 * File: item.go
 *
 * Description: Class to handle the various small items that can be picked
 *              up or equipped.
 */

package main

import "fmt"

// Structure to hold the attributes of an item.
type Item struct {

    // Holds the name of the given item.
    name string

    // Holds the type of the given item.
    category string

    // Store the current (x,y) coord of the item.
    Y  Coord
    X  Coord

    // Appearance of the item.
    ch rune

    // Pointer to the level of where the item is.
    area *Area

    // Determines whether an item can be equipped by a creature.
    can_equip bool

    // Variable that holds whether or not an item is broken.
    is_broken bool

    // The current condition of an item. If the max durability and the
    // current are equal, then the item is in perfect condition. 
    durability_current int
    durability_maximum int

    // Price of the item in gold coins.
    price_to_purchase int
    price_to_sell int

    // Weight of the item in grams.
    weight int

    // Variables that determine how much an item increases (if positive)
    // or decreases (if negative) the attack / defence of a creature. 
    attack_increase int
    defence_increase int
}

//! Item constructor function.
/*
 * @param     string    item name
 * @param     string    category
 * @param     Coord     Y
 * @param     Coord     X
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
             Y Coord,
             X Coord,
             ch rune,
             area *Area,
             can_equip bool,
             is_broken bool,
             durability_current int,
             durability_maximum int,
             price_to_purchase int,
             price_to_sell int,
             weight int,
             attack_increase int,
             defence_increase int) *Item {

    // Return an address to a newly allocated Creature object.
    return &Item{name,
                 category,
                 Y,
                 X,
                 ch,
                 area,
                 can_equip,
                 is_broken,
                 durability_current,
                 durability_maximum,
                 price_to_purchase,
                 price_to_sell,
                 weight,
                 attack_increase,
                 defence_increase}
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
func (itm *Item) adjust_durability(amount int) {

    // If debug, state how much the item durability currently is.
    DebugLog(&G, fmt.Sprintf("adjust_durability() --> Item [%s] durability" +
                             "before adjustment is: %d / %d",
                             itm.name,
                             itm.durability_current,
                             itm.durability_maximum))

    // Adjust the current durability by the amount specified.
    itm.durability_current += amount

    // If debug, state how much the item durability has changed by.
    DebugLog(&G, fmt.Sprintf("adjust_durability() --> item [%s] durability" +
                             "after adjustment is: %d / %d",
                             itm.name,
                             itm.durability_current,
                             itm.durability_maximum))

    // If the current item durability has reached beyond the max, go ahead
    // and cap it to the maximum durability.
    if itm.durability_current > itm.durability_maximum {

        // Cap the durability of the item.
        itm.durability_current = itm.durability_maximum

        // If debug mode, tell the developer what just happened.
        DebugLog(&G, fmt.Sprintf("adjust_durability() --> item [%s] has " +
                                 "exceeded max durability, and so has " +
                                 "been capped"))
    }

    // If the item is less than zero, then consider it to be broken, in
    // which event t
    if itm.durability_current < 1 {

        // If debug mode, tell the developer this is going to break the item.
        DebugLog(&G, fmt.Sprintf("adjust_durability() --> item [%s] is " +
                                 "preparing to be broken", itm.name))

        // Set the relevant item properties via this function in order to
        // render the item broken.
        itm.event_broken()

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
func (itm *Item) event_broken() {

    // Prevent the item from being equipped.
    itm.can_equip = false

    // Set the is_broken flag to true.
    itm.is_broken = true

    // If debug mode, tell the developer that the item is now broken.
    DebugLog(&G, fmt.Sprintf("event_broken() --> item [%s] is now broken",
                             itm.name))
}
