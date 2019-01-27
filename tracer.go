package main

import (
	"bufio"
	"github.com/evdokimovn/ray-tracer/objects"
	"github.com/evdokimovn/ray-tracer/objects/sphere"
	"github.com/evdokimovn/ray-tracer/vector"
	"math"
	"os"
	"runtime"
	"strconv"
)

const (
	width  = 1024
	height = 768
	fov    = math.Pi / 3
)

func main() {
	var (
		ivory     = &objects.Material{RefractiveIndex: 1.0, Albedo: vector.New4(0.6, 0.3, 0.1, 0), DiffuseColor: vector.New3(0.4, 0.4, 0.3), SpecularComponent: 50}
		glass     = &objects.Material{RefractiveIndex: 1.5, Albedo: vector.New4(0.0, 0.5, 0.1, 0.8), DiffuseColor: vector.New3(0.6, 0.7, 0.8), SpecularComponent: 125}
		redRubber = &objects.Material{RefractiveIndex: 1.0, Albedo: vector.New4(0.9, 0.1, 0.0, 0), DiffuseColor: vector.New3(0.3, 0.1, 0.1), SpecularComponent: 10}
		mirror    = &objects.Material{RefractiveIndex: 1.0, Albedo: vector.New4(0.0, 10.0, 0.8, 0), DiffuseColor: vector.New3(1.0, 1.0, 1.0), SpecularComponent: 1425}
	)

	var spheres []*sphere.Sphere
	spheres = append(spheres, sphere.New(vector.New3(-3, 0, -16), 2, ivory))
	spheres = append(spheres, sphere.New(vector.New3(-1.0, -1.5, -12), 2, glass))
	spheres = append(spheres, sphere.New(vector.New3(1.5, -0.5, -18), 3, redRubber))
	spheres = append(spheres, sphere.New(vector.New3(7, 5, -18), 4, mirror))

	var lights []*objects.Light
	lights = append(lights, &objects.Light{Position: vector.New3(-20, 20, 20), Intensity: 1.5})
	lights = append(lights, &objects.Light{Position: vector.New3(30, 50, -25), Intensity: 1.8})
	lights = append(lights, &objects.Light{Position: vector.New3(30, 20, 30), Intensity: 1.7})

	render(spheres, lights)
}

func castRay(orig, dir *vector.Vec3, spheres []*sphere.Sphere, lights []*objects.Light, depth int) *vector.Vec3 {
	var (
		point    = new(vector.Vec3)
		N        *vector.Vec3
		material = objects.NewMaterial()
	)

	if depth > 4 || !sceneIntersect(orig, dir, spheres, &point, &N, &material) {
		return vector.New3(0.2, 0.7, 0.8)
	}

	reflectDir := objects.Reflect(dir, N).Normalize()
	refractDir := objects.Refract(dir, N, material.RefractiveIndex).Normalize()
	var o *vector.Vec3
	if reflectDir.Dot(N) < 0 {
		o = point.Diff(N.Scale(1e-3))
	} else {
		o = point.Add(N.Scale(1e-3))
	} // offset the original point to avoid occlusion by the object itself

	reflectOrig := o

	if refractDir.Dot(N) < 0 {
		o = point.Diff(N.Scale(1e-3))
	} else {
		o = point.Add(N.Scale(1e-3))
	} // offset the original point to avoid occlusion by the object itself

	refractOrig := o

	reflectColor := castRay(reflectOrig, reflectDir, spheres, lights, depth+1)
	refractColor := castRay(refractOrig, refractDir, spheres, lights, depth+1)

	var diffuseLightIntensity, specularLightIntensity float64
	for i := 0; i < len(lights); i++ {
		light := lights[i].Position.Diff(point)
		lightDir := light.Normalize()
		lightDistance := light.Length()

		var shadowOrigin *vector.Vec3
		if lightDir.Dot(N) < 0 {
			shadowOrigin = point.Diff(N.Scale(1e-3))
		} else {
			shadowOrigin = point.Add(N.Scale(1e-3))
		}

		var shadowPT, shadowN *vector.Vec3
		var tmpMaterial = objects.NewMaterial()
		if sceneIntersect(shadowOrigin, lightDir, spheres, &shadowPT, &shadowN, &tmpMaterial) &&
			(shadowPT.Diff(shadowOrigin)).Length() < lightDistance {
			continue
		}

		diffuseLightIntensity += lights[i].Intensity * math.Max(0, lightDir.Dot(N))
		specularLightIntensity += math.Pow(
			math.Max(0, -objects.Reflect(lightDir.Scale(-1), N).Dot(dir)),
			material.SpecularComponent) * lights[i].Intensity
	}

	return material.DiffuseColor.Scale(diffuseLightIntensity).Scale(material.Albedo.X()).
		Add(vector.New3(1., 1., 1.).Scale(specularLightIntensity).Scale(material.Albedo.Y())).
		Add(reflectColor.Scale(material.Albedo.Z())).
		Add(refractColor.Scale(material.Albedo.M()))
}

func sceneIntersect(orig, dir *vector.Vec3, spheres []*sphere.Sphere, hit, N **vector.Vec3, material **objects.Material) bool {
	sphereDist := math.MaxFloat64
	for i := 0; i < len(spheres); i += 1 {
		var distI *float64
		if spheres[i].IntersectsRay(orig, dir, &distI) && *distI < sphereDist {
			d := *distI
			sphereDist = d
			*hit = orig.Add(dir.Scale(d))
			*N = ((*hit).Diff(spheres[i].Center)).Normalize()
			*material = spheres[i].Material
		}
	}

	checkerboardDist := math.MaxFloat64
	if math.Abs(dir.Y()) > 1e-3 {
		d := -(orig.Y() + 4) / dir.Y() // the checkerboard plane has equation y = -4
		pt := orig.Add(dir.Scale(d))
		if d > 0 && math.Abs(pt.X()) < 10 && pt.Z() < -10 && pt.Z() > -30 && d < sphereDist {
			checkerboardDist = d
			*hit = pt
			*N = vector.New3(0, 1, 0)
			if ((int(.5*(*hit).X()+1000) + int(.5*(*hit).Z())) & 1) == 1 {
				(*material).DiffuseColor = vector.New3(1, 1, 1).Scale(0.3)
			} else {
				(*material).DiffuseColor = vector.New3(1, 0.7, 0.3).Scale(0.3)
			}
		}
	}
	return math.Min(sphereDist, checkerboardDist) < 1000
}

func render(spheres []*sphere.Sphere, lights []*objects.Light) {
	var frameBuffer = make([]*vector.Vec3, height*width)

	l := runtime.NumCPU()
	var sem = make(chan struct{}, l)
	z := -height / (2. * math.Tan(fov/2.))
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			sem <- struct{}{}
			go func(i, j int) {
				x := (float64(i) + 0.5) - width/2.
				y := -(float64(j) + 0.5) + height/2.
				dir := vector.New3(x, y, z).Normalize()
				frameBuffer[i+j*width] = castRay(vector.New3(0, 0, 0), dir, spheres, lights, 0)
				<-sem
			}(i, j)
		}
	}

	for len(sem) != 0 {
	}

	f, err := os.Create("./out.ppm")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString("P6\n")
	w.WriteString(strconv.Itoa(width))
	w.WriteString(" ")
	w.WriteString(strconv.Itoa(height))
	w.WriteString("\n255\n")

	for i := 0; i < height*width; i += 1 {
		c := frameBuffer[i]
		max := math.Max(c.X(), math.Max(c.Y(), c.Z()))
		if max > 1 {
			frameBuffer[i] = c.Scale(1. / max)
		}
		w.WriteByte(byte(255.99 * frameBuffer[i].X()))
		w.WriteByte(byte(255.99 * frameBuffer[i].Y()))
		w.WriteByte(byte(255.99 * frameBuffer[i].Z()))
	}
	w.Flush()
}
