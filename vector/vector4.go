package vector

type Vec4 struct {
	*Vec3
	m float64
}

func (v *Vec4) Add(a *Vec4) *Vec4 {
	new3 := v.Vec3.Add(a.Vec3)
	return &Vec4{
		new3,
		v.m + a.m,
	}
}

func New4(x, y, z, m float64) *Vec4 {
	vec3 := New3(x, y, z)
	return &Vec4{
		vec3,
		m,
	}
}

func (v *Vec4) M() float64 {
	return v.m
}

func (v *Vec4) Dot(o *Vec4) float64 {
	d := v.Vec3.Dot(o.Vec3)
	return d + v.m*o.m
}
