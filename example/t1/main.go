package main

import (
	"log"
	"math"
)

func main() {
	var x1, x2, y1, y2 = 600.0, 489.624, 200.000000, 114.766276

	sx, lx, sy, ly := math.Min(x1, x2), math.Max(x1, x2), math.Min(y1, y2), math.Max(y1, y2)
	log.Printf("sx:%f, lx:%f, sy:%f, ly:%f, y:%f, x:%f", sx, lx, sy, ly, sy-ly, lx-sx)
}
