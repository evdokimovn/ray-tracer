package objects

import (
	"github.com/evdokimovn/ray-tracer/vector"
	"math"
)

type Material struct {
	DiffuseColor      *vector.Vec3
	SpecularComponent float64
	Albedo            *vector.Vec4
	RefractiveIndex   float64
}

func NewMaterial() *Material {
	return &Material{
		DiffuseColor:    vector.New3(0, 0, 0),
		Albedo:          vector.New4(1, 0, 0, 0),
		RefractiveIndex: 1,
	}
}

func Reflect(L, N *vector.Vec3) *vector.Vec3 {
	return L.Diff(N.Scale(2).Scale(L.Dot(N)))
}

func refract(L, N *vector.Vec3, refractiveIndexMedia, refractiveIndexAir float64) *vector.Vec3 {
	cosi := -math.Max(-1, math.Min(1, L.Dot(N)))
	// If the ray comes from the inside the object, swap the air and the media
	if cosi < 0 {
		return refract(L, N.Scale(-1), refractiveIndexAir, refractiveIndexMedia)
	}
	eta := refractiveIndexAir / refractiveIndexMedia
	k := 1 - eta*eta*(1-cosi*cosi)
	if k < 0 {
		return vector.New3(1, 0, 0)
	}
	return L.Scale(eta).Add(N.Scale(eta*cosi - math.Sqrt(k)))
}

// Refract finds direction of refracted light using Snell's law
func Refract(L, N *vector.Vec3, refractiveIndex float64) *vector.Vec3 {
	return refract(L, N, refractiveIndex, 1)
}
