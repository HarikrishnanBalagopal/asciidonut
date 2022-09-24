package utils

import "math"

func torus_sd(p Vec3, tx, ty float64) float64 {
	q := Vec3{X: p.XZ().Length() - tx, Y: p.Y, Z: 0}
	return q.Length() - ty
}
func torus_sdg(p Vec3, ra, rb float64) (float64, Vec3) {
	h := p.XZ().Length()
	return Vec3{h - ra, p.Y, 0}.Length() - rb, p.Mul(Vec3{h - ra, h, h - ra}).Normalize()
}

func Scene_sd(p Vec3, m Mat3x3) float64 {
	p1 := m.MulV(p)
	d := torus_sd(p1, 2., .75)
	return d
}

func Normals_sd(p Vec3, m Mat3x3) Vec3 {
	p1 := m.MulV(p)
	_, n := torus_sdg(p1, 2., .75)
	n1 := m.RightMulV(n)
	return n1
	// dx1 := Scene_sd(p.Add(Vec3{EPS, 0, 0}), m)
	// dx2 := Scene_sd(p.Add(Vec3{-EPS, 0, 0}), m)

	// dy1 := Scene_sd(p.Add(Vec3{0, EPS, 0}), m)
	// dy2 := Scene_sd(p.Add(Vec3{0, -EPS, 0}), m)

	// dz1 := Scene_sd(p.Add(Vec3{0, 0, EPS}), m)
	// dz2 := Scene_sd(p.Add(Vec3{0, 0, -EPS}), m)

	// return Vec3{
	// 	X: dx2 - dx1/2*EPS,
	// 	Y: dy2 - dy1/2*EPS,
	// 	Z: dz2 - dz1/2*EPS,
	// }.Normalize()
}

func GetRotMatAtTime(t float64) Mat3x3 {
	r1 := Mat3x3{}.Rot(Vec3{1, 0, 0}, -.1*.7*t)
	r2 := Mat3x3{}.Rot(Vec3{0, 0, 1}, .1*math.Pi*t)
	return r2.Mul(r1)
}
