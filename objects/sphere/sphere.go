package sphere

import (
	"github.com/evdokimovn/ray-tracer/objects"
	"github.com/evdokimovn/ray-tracer/vector"
	"math"
)

type Sphere struct {
	Center   *vector.Vec3
	radius   float64
	Material *objects.Material
}

func New(center *vector.Vec3, radius float64, material *objects.Material) *Sphere {
	return &Sphere{
		Center:   center,
		radius:   radius,
		Material: material,
	}
}

func (s *Sphere) IntersectsRay(orig, dir *vector.Vec3, t0 **float64) bool {
	L := s.Center.Diff(orig)
	tca := L.Dot(dir)
	d2 := L.Length2() - tca*tca
	radiusSquared := s.radius * s.radius
	if d2 > radiusSquared {
		return false
	}
	thc := math.Sqrt(radiusSquared - d2)
	tmp := tca - thc
	t1 := tca + thc
	if tmp < 0 {
		tmp = t1
	}
	*t0 = &tmp
	if tmp < 0 {
		return false
	}
	return true
}
