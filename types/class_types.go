/*
 * File: types/class_types.go
 *
 * Description: Hold type information about creature classes
 */

package types

// Structure to hold creature information
type ClassTypeInfo struct {

	// Holds the name of the given item.
	Name string

	// The innate abilities of the given class
	//
	// "warrior" => has warrior abilities
	// "thief" => has warrior abilities
	// "cleric" => has cleric abilities
	// "wizard" => has wizard abilities
	// "unknown" => default null value
	//
	HasAbilities string

	// Which ability score defines a given class
	//
	// "strength" => class requires 14 strength
	// "intelligence" => class requires 14 intelligence
	// "agility" => class requires 14 agility
	// "wisdom" => class requires 14 wisdom
	// "unknown" => default null value
	//
	EssentialAttribute string
}

//! Function to populate details about various class types
/*
 * @return    none
 */
func GenClassTypes(clstype map[string]ClassTypeInfo) bool {

	// Input validation
	if clstype == nil {
		return false
	}

	//
	// Unknown
	//
	clstype["0"] = ClassTypeInfo{"Unknown", "unknown", "unknown"}

	//
	// Warrior
	//
	clstype["1"] = ClassTypeInfo{"Warrior", "warrior", "strength"}

	//
	// Wizard
	//
	clstype["2"] = ClassTypeInfo{"Wizard", "wizard", "intelligence"}

	//
	// Thief
	//
	clstype["3"] = ClassTypeInfo{"Thief", "thief", "agility"}

	//
	// Cleric
	//
	clstype["4"] = ClassTypeInfo{"Cleric", "cleric", "wisdom"}

	// All of the classes have been populated successfully, so return true.
	return true
}
