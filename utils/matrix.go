package utils

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func (v1 Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{X: v1.X + v2.X, Y: v1.Y + v2.Y, Z: v1.Z + v2.Z}
}

func (v1 Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{X: v1.X - v2.X, Y: v1.Y - v2.Y, Z: v1.Z - v2.Z}
}

func (v1 Vec3) Mul(v2 Vec3) Vec3 {
	return Vec3{X: v1.X * v2.X, Y: v1.Y * v2.Y, Z: v1.Z * v2.Z}
}

func (v Vec3) Scale(amount float64) Vec3 {
	return Vec3{X: v.X * amount, Y: v.Y * amount, Z: v.Z * amount}
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) LengthSq() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3) Normalize() Vec3 {
	l := v.Length()
	return Vec3{X: v.X / l, Y: v.Y / l, Z: v.Z / l}
}

func (v1 Vec3) Dot(v2 Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{X: v1.Y*v2.Z - v1.Z*v2.Y, Y: v1.Z*v2.X - v1.X*v2.Z, Z: v1.X*v2.Y - v1.Y*v2.X}
}

func (v Vec3) XZ() Vec3 {
	return Vec3{X: v.X, Y: 0, Z: v.Z}
}

// ------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------

type Mat3x3 struct {
	XX, XY, XZ float64
	YX, YY, YZ float64
	ZX, ZY, ZZ float64
}

func (m Mat3x3) MulV(v Vec3) Vec3 {
	return Vec3{
		X: m.XX*v.X + m.XY*v.Y + m.XZ*v.Z,
		Y: m.YX*v.X + m.YY*v.Y + m.YZ*v.Z,
		Z: m.ZX*v.X + m.ZY*v.Y + m.ZZ*v.Z,
	}
}

// RightMulV is 1x3 = 1x3 3x3
// It's can also be thought of as transpose and multiply.
func (m Mat3x3) RightMulV(v Vec3) Vec3 {
	return Vec3{
		X: m.XX*v.X + m.YX*v.Y + m.ZX*v.Z,
		Y: m.XY*v.X + m.YY*v.Y + m.ZY*v.Z,
		Z: m.XZ*v.X + m.YZ*v.Y + m.ZZ*v.Z,
	}
}

func (m1 Mat3x3) Mul(m2 Mat3x3) Mat3x3 {
	return Mat3x3{
		XX: m1.XX*m2.XX + m1.XY*m2.YX + m1.XZ*m2.ZX,
		XY: m1.XX*m2.XY + m1.XY*m2.YY + m1.XZ*m2.ZY,
		XZ: m1.XX*m2.XZ + m1.XY*m2.YZ + m1.XZ*m2.ZZ,

		YX: m1.YX*m2.XX + m1.YY*m2.YX + m1.YZ*m2.ZX,
		YY: m1.YX*m2.XY + m1.YY*m2.YY + m1.YZ*m2.ZY,
		YZ: m1.YX*m2.XZ + m1.YY*m2.YZ + m1.YZ*m2.ZZ,

		ZX: m1.ZX*m2.XX + m1.ZY*m2.YX + m1.ZZ*m2.ZX,
		ZY: m1.ZX*m2.XY + m1.ZY*m2.YY + m1.ZZ*m2.ZY,
		ZZ: m1.ZX*m2.XZ + m1.ZY*m2.YZ + m1.ZZ*m2.ZZ,
	}
}

func (Mat3x3) Rot(axis Vec3, angle float64) Mat3x3 {
	return Mat3x3{
		XX: math.Cos(angle) + axis.X*axis.X*(1-math.Cos(angle)),
		XY: axis.X*axis.Y*(1-math.Cos(angle)) - axis.Z*math.Sin(angle),
		XZ: axis.X*axis.Z*(1-math.Cos(angle)) + axis.Y*math.Sin(angle),

		YX: axis.Y*axis.X*(1-math.Cos(angle)) + axis.Z*math.Sin(angle),
		YY: math.Cos(angle) + axis.Y*axis.Y*(1-math.Cos(angle)),
		YZ: axis.Y*axis.Z*(1-math.Cos(angle)) - axis.X*math.Sin(angle),

		ZX: axis.Z*axis.X*(1-math.Cos(angle)) - axis.Y*math.Sin(angle),
		ZY: axis.Z*axis.Y*(1-math.Cos(angle)) + axis.X*math.Sin(angle),
		ZZ: math.Cos(angle) + axis.Z*axis.Z*(1-math.Cos(angle)),
	}
}
