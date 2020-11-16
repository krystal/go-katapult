package namegenerator

/*
 * Top-level convenience functions
 */

var globalGenerator = New(DefaultAdjectives, DefaultNouns)

func RandomHostname() string {
	return globalGenerator.RandomHostname()
}

func RandomName(prefixes ...string) string {
	return globalGenerator.RandomName(prefixes...)
}
