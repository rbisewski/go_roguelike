/*
 * File: misc.go
 *
 * Description: Various functions that don't fit elsewhere.
 */

package main

import "math"
import "math/rand"

//! Randomly returns "true" or "false"
/*
 * @returns    bool    whether the coin was heads (true) or tails (false)
 */
func TossCoin() bool {
    return rand.Intn(100) > 50;
}

//! Randomly returns "true" or "false"
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
    return rand.Intn(maximum);
}

//! Get the minimum of a list of int values (i.e. the lowest value)
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

//! Get the max of a list of int values (i.e. the lowest value)
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

//! Return the x-percentage of a given value.
/*
 *  @param    int   current
 *  @param    int   total
 *
 * @return    int   the x% of something
 */
func Percent(percent, of int) int {
    return of * percent / 100
}

//! Round a given value or float to an integer.
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
        intermed -= 1
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
