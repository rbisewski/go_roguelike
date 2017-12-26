/*
 * File: misc.go
 *
 * Description: Simple utility functions for the program.
 */

package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

// TossCoin ... Randomly returns "true" or "false"
/*
 * @returns    bool    whether the coin was heads (true) or tails (false)
 */
func TossCoin() bool {
	return rand.Intn(100) > 50
}

// getRandomNumBetweenZeroAndMax ... Randomly returns "true" or "false"
/*
 * @returns    bool    whether the coin was heads (true) or tails (false)
 */
func getRandomNumBetweenZeroAndMax(maximum int) int {

	// Input validation, make sure this is at least zero.
	if maximum < 1 {

		// As a default, give back a zero.
		return 0
	}

	// Otherwise return a number between 0 and the maximum.
	return rand.Intn(maximum)
}

// Min ... Get the minimum of a list of int values (i.e. the lowest value)
/*
 * @param    int    a given integer value
 * @param    ...
 *
 * @returns  int    lowest value
 */
func Min(a ...int) int {

	// Use bitshifting to get the largest possible value for a
	// golang integer.
	min := int(^uint(0) >> 1)

	// Cycle thru all of the given int values...
	for _, i := range a {

		// Smaller? Set the variable then.
		if i < min {
			min = i
		}
	}

	// Finally return the minimum value
	return min
}

// Max ... Get the max of a list of int values (i.e. the lowest value)
/*
 * @param    int    a given integer value
 * @param    ...
 *
 * @returns  int    lowest value
 */
func Max(a ...int) int {

	// Use bitshifting the smallest possible value for a golang int.
	max := -(int(^uint(0)>>1) - 1)

	// Cycle thru the range of integers.
	for _, i := range a {

		// Bigger? Then use that one.
		if i > max {
			max = i
		}
	}

	// Finally return the max
	return max
}

// Percent ... Return the x-percentage of a given value.
/*
 *  @param    int   current
 *  @param    int   total
 *
 * @return    int   the x% of something
 */
func Percent(percent, of int) int {
	return of * percent / 100
}

// Round a given value or float to an integer.
/*
 * @param   float64    given value to round
 *
 * @return  int        rounded result as int
 */
func Round(x float64) int {

	// Define the precise (e.g. 1 == int)
	prec := 1

	// Variable to store the rounded values.
	var rounder float64

	// Set the X^Y power based on the earlier defined precise.
	pow := math.Pow(10, float64(prec))

	// Determine the intermediate
	intermed := x * pow

	// Grab the modulo
	_, frac := math.Modf(intermed)

	// Increment the intermediate (i.e. between 0 and 1 is 0.5)
	intermed += .5

	// If the fraction is negative, invert it.
	x = .5
	if frac < 0.0 {
		x = -.5
		intermed--
	}

	// Fraction 0.5 or more? Ceiling then.
	if frac >= x {
		rounder = math.Ceil(intermed)

		// Else grab the floor instead since less than 0.5
	} else {
		rounder = math.Floor(intermed)
	}

	// Safety check, if pow is 0 then simply return zero. This probably
	// won't happen, but it's usually a good idea to prevent divide by
	// zero errors.
	if pow == 0 {
		return 0
	}

	// Finally dump to an int, as per the precision.
	return int(rounder / pow)
}

// IsNumeric ... Evaluates if a given ASCII character is a number (0-9)
/*
 * @param    string    Keyboard ASCII character input.
 *
 * @return   bool      whether or not the character is a numeric
 */
func IsNumeric(character string) bool {

	// Input validation, make sure this actually was given a valid string
	// ASCII character of length 1.
	if len(character) != 1 {
		return false
	}

	// Convert the key pressed to a hex string value.
	charAsHex := fmt.Sprintf("%x", character)

	// Sanity check, make sure this actually was able to return non-blank.
	if len(charAsHex) < 1 {
		return false
	}

	// Attempt to convert the hexadecimal value to an uint16 decimal.
	charAsInt, err := strconv.ParseUint(charAsHex, 16, 16)

	// Safety check, if an error occurred, this probably isn't a number,
	// so go ahead and return false.
	if err != nil {
		return false
	}

	// Determine if the character given is 0-9
	//
	// 0 -> 0x30 -> 48
	// 9 -> 0x39 -> 57
	//
	if charAsInt > 47 && charAsInt < 58 {
		return true
	}

	// Otherwise some other sort of character was present here, so then
	// return false here as a default.
	return false
}

// ConvertKeyToNumeric ... converts a keyboard number key to integer value
/*
 * @param    string    Keyboard ASCII character input.
 *
 * @return   uint64    integer output
 * @return   error     error message, if any
 */
func ConvertKeyToNumeric(character string) (uint64, error) {

	// Input validation, make sure this actually was given a valid string
	// ASCII character of length 1.
	if len(character) != 1 {
		return 0, fmt.Errorf("ConvertKeyToNumeric() --> invalid input")
	}

	// Variable declaration
	var num uint64

	// Convert the key pressed to a hex string value.
	charAsHex := fmt.Sprintf("%x", character)

	// Sanity check, make sure this actually was able to return non-blank.
	if len(charAsHex) < 1 {
		return 0, fmt.Errorf("ConvertKeyToNumeric() --> improper hexidecimal")
	}

	// Attempt to convert the hexadecimal value to an uint16 decimal.
	charAsInt, err := strconv.ParseUint(charAsHex, 16, 16)

	// Safety check, if an error occurred, this probably isn't a number,
	// so go ahead and return false.
	if err != nil {
		return 0, err
	}

	// Determine if the character given is 0-9
	//
	// 0 -> 0x30 -> 48
	// 9 -> 0x39 -> 57
	//
	if charAsInt < 48 || charAsInt > 57 {
		return 0, fmt.Errorf("ConvertKeyToNumeric() --> non-ASCII value")
	}

	// Convert the value into a plain ol' unsigned int.
	num = charAsInt - 48

	// Otherwise some other sort of character was present here, so then
	// return false here as a default.
	return num, nil
}

// IsAlphaCharacter ... evaluates if ASCII char is alphabetical (a-zA-Z)
/*
 * @param    string    Keyboard ASCII character input.
 *
 * @return   bool      whether or not the character is an alphabetical letter
 */
func IsAlphaCharacter(character string) bool {

	// Input validation, make sure this actually was given a valid string
	// ASCII character of length 1.
	if len(character) != 1 {
		return false
	}

	// Convert the key pressed to a hex string value.
	charAsHex := fmt.Sprintf("%x", character)

	// Sanity check, make sure this actually was able to return non-blank.
	if len(charAsHex) < 1 {
		return false
	}

	// Attempt to convert the hexadecimal value to an uint16 decimal.
	charAsInt, err := strconv.ParseUint(charAsHex, 16, 16)

	// Safety check, if an error occurred, this probably isn't alphabetical,
	// so go ahead and return false.
	if err != nil {
		return false
	}

	// Determine if the character given is a-zA-Z
	//
	// A -> 0x41 -> 65
	// Z -> 0x5a -> 90
	// a -> 0x61 -> 97
	// z -> 0x7a -> 122
	//
	if (charAsInt > 64 && charAsInt < 91) || (charAsInt > 96 && charAsInt < 123) {
		return true
	}

	// Otherwise some other sort of character was present here, so then
	// return false here as a default.
	return false
}

// IsDeleteOrBackspace ... check if key is ASCII backspace or delete char
/*
 * @param    string    Keyboard ASCII character input.
 *
 * @return   bool      whether or not the character is backspace or delete
 */
func IsDeleteOrBackspace(character string) bool {

	// Input validation, make sure this actually was given a valid string
	// ASCII character of length 1.
	if len(character) != 1 {
		return false
	}

	// Convert the key pressed to a hex string value.
	charAsHex := fmt.Sprintf("%x", character)

	// Sanity check, make sure this actually was able to return non-blank.
	if len(charAsHex) < 1 {
		return false
	}

	// Determine if the character given is the "Backspace" or "Delete" key.
	if charAsHex == "7f" || charAsHex == "c58a" {
		return true
	}

	// Otherwise some other sort of character was present here, so then
	// return false here as a default.
	return false
}

// WasEnterPressed ... check if the enter key was pressed.
/*
 * @param    string    Keyboard ASCII character input.
 *
 * @return   bool      whether or not the character is the enter key.
 */
func WasEnterPressed(character string) bool {

	// Input validation, make sure this actually was given a valid string
	// ASCII character of length 1.
	if len(character) != 1 {
		return false
	}

	// Convert the key pressed to a hex string value.
	charAsHex := fmt.Sprintf("%x", character)

	// Sanity check, make sure this actually was able to return non-blank.
	if len(charAsHex) < 1 {
		return false
	}

	// Determine if the character given is the "Enter" key.
	if charAsHex == "0a" {
		return true
	}

	// Otherwise some other sort of character was present here, so then
	// return false here as a default.
	return false
}

// AlignAndSpaceString ... align and space a given string
/*
 * @return    none
 */
func AlignAndSpaceString(phrase string, alignment string, length int) string {

	// Ensure that the phrase has length of at least 1
	if len(phrase) < 1 {
		return ""
	}

	// Ensure that the
	if len(alignment) < 1 {
		return ""
	}

	// Ensure that the length is greater that zero
	if length < 1 {
		return ""
	}

	// Handle each of the different cases.
	switch alignment {

	//
	// Right alignment
	//
	case "right":
		for i := len(phrase); i < length; i++ {
			phrase += " "
		}

	//
	// Centre alignment
	//
	case "middle":
		fallthrough
	case "center":
		fallthrough
	case "centre":
		for i := 0; i < length; i++ {
			if (i % 2) == 1 {
				phrase += " "
			} else {
				phrase = " " + phrase
			}
		}

	//
	// Left alignment
	//
	case "left":
		for i := len(phrase); i < length; i++ {
			phrase = " " + phrase
		}

	//
	// By default, do nothing...
	//
	default:
		break
	}

	// Return the aligned and spaced string.
	return phrase
}
