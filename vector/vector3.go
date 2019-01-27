package vector

import (
	"math"
)

type Vec3 struct {
	x float64
	y float64
	z float64
}

func New3(x, y, z float64) *Vec3 {
	return &Vec3{
		x: x,
		y: y,
		z: z,
	}
}

func (v *Vec3) Add(a *Vec3) *Vec3 {
	return &Vec3{
		x: v.x + a.x,
		y: v.y + a.y,
		z: v.z + a.z,
	}
}

func (v Vec3) X() float64 {
	return v.x
}

func (v Vec3) Y() float64 {
	return v.y
}

func (v Vec3) Z() float64 {
	return v.z
}

func (v *Vec3) Length() float64 {
	return math.Sqrt(v.Length2())
}

func (v *Vec3) Length2() float64 {
	return v.Dot(v)
}

func (v *Vec3) Dot(a *Vec3) float64 {
	return v.x*a.x + v.y*a.y + v.z*a.z
}

func (v *Vec3) Diff(a *Vec3) *Vec3 {
	return New3(v.x-a.x, v.y-a.y, v.z-a.z)
}

func (v *Vec3) Scale(factor float64) *Vec3 {
	return New3(v.x*factor, v.y*factor, v.z*factor)
}

func (v *Vec3) Normalize() *Vec3 {
	l := v.Length()
	return New3(v.x/l, v.y/l, v.z/l)
}
