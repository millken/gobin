## Schema

Data structures are defined in .schema files. The format is quite self-explanatory.
```
// Package demo offers a demonstration.
// These comment lines will end up in the generated code.
package demo

// Course is the grounds where the game of golf is played.
type course struct {
	ID    uint64
	name  text
	holes []hole
	image binary
	tags  []text
}

type hole struct {
	// Lat is the latitude of the cup.
	lat float64
	// Lon is the longitude of the cup.
	lon float64
	// Par is the difficulty index.
	par uint8
	// Water marks the presence of water.
	water bool
	// Sand marks the presence of sand.
	sand bool
}
```