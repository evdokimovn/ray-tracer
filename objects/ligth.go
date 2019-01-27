package objects

import "github.com/evdokimovn/ray-tracer/vector"

type Light struct {
	Position  *vector.Vec3
	Intensity float64
}
