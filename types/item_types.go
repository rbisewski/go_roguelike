/*
 * File: types/item_types.go
 *
 * Description: Hold type information about items.
 */

package types

// Structure to hold creature information
type ItemTypeInfo struct {

    // Holds the name of the given item.
    Name string

    // Holds the type of the given item.
    Category string

    // Appearance of the item.
    Ch rune

    // Determines whether an item can be equipped by a creature.
    Can_equip bool

    // Variable that holds whether or not an item is broken.
    Is_broken bool

    // The current condition of an item. If the max durability and the
    // current are equal, then the item is in perfect condition. 
    Durability_current int
    Durability_maximum int

    // Price of the item in gold coins.
    Price_to_purchase int
    Price_to_sell int

    // Weight of the item in grams.
    Weight int

    // Variables that determine how much an item increases (if positive)
    // or decreases (if negative) the attack / defence of a creature. 
    Attack_increase int
    Defence_increase int
}

//
// TODO: add item functions here
//
