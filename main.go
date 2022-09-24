package main

import (
	"fmt"
	"math"
	"strings"
	"syscall/js"
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

func (v Vec3) Scale(amount float64) Vec3 {
	return Vec3{X: v.X * amount, Y: v.Y * amount, Z: v.Z * amount}
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
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

func (m Mat3x3) Mul(v Vec3) Vec3 {
	return Vec3{
		X: m.XX*v.X + m.XY*v.Y + m.XZ*v.Z,
		Y: m.YX*v.X + m.YY*v.Y + m.YZ*v.Z,
		Z: m.ZX*v.X + m.ZY*v.Y + m.ZZ*v.Z,
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

func Radians(degrees float64) float64 {
	return math.Pi * degrees / 180.0
}

// ------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------

const (
	N              = 32
	W              = N * 2
	H              = N
	MAX_ITERATIONS = N
	EPS            = 0.001
	// https://en.wikipedia.org/wiki/Phong_reflection_model
	AMBIENT_REFLECTION_CONSTANT = 0.05
	AMBIENT_LIGHT_INTENSITY     = 255.0
	DIFFUSE_REFLECTION_CONSTANT = 0.95
	DIFFUSE_LIGHT_INTENSITY     = 255.0
	// SPECULAR_REFLECTION_CONSTANT = 0.1
	// MATERIAL_SHININESS_CONSTANT = 0.1
)

var (
	ASCII    = []byte{' ', '.', ',', '-', '+', '*', 'o', '0', '@', '#'}
	BUFFER_1 = [H][W]byte{}
)

// func sphere_sd(p Vec3, r float64) float64 {
// 	return p.Length() - r
// }

func torus_sd(p Vec3, tx, ty float64) float64 {
	q := Vec3{X: p.XZ().Length() - tx, Y: p.Y, Z: 0}
	return q.Length() - ty
}

func scene_sd(p Vec3, t float64) float64 {
	p1 := Mat3x3{}.Rot(Vec3{1, 0, 0}, t).Mul(p)
	// p1 := Mat3x3{}.Rot(Vec3{1, 1, math.Sin(7 * t)}.Normalize(), 10*t).Mul(p)
	d := torus_sd(p1, 2., .75)
	// d := sphere_sd(p, 1.)
	return d
}

func normals_sd(p Vec3, t float64) Vec3 {
	dx1 := scene_sd(p.Add(Vec3{EPS, 0, 0}), t)
	dx2 := scene_sd(p.Add(Vec3{-EPS, 0, 0}), t)

	dy1 := scene_sd(p.Add(Vec3{0, EPS, 0}), t)
	dy2 := scene_sd(p.Add(Vec3{0, -EPS, 0}), t)

	dz1 := scene_sd(p.Add(Vec3{0, 0, EPS}), t)
	dz2 := scene_sd(p.Add(Vec3{0, 0, -EPS}), t)

	return Vec3{
		X: dx2 - dx1/2*EPS,
		Y: dy2 - dy1/2*EPS,
		Z: dz2 - dz1/2*EPS,
	}.Normalize()
}

func calc_pixel_color(x, y, t float64) float64 {
	nx := 2*x/W - 1
	ny := 2*y/H - 1

	// fmt.Println("nx", nx, "ny", ny)

	cam_pos := Vec3{0, 0, 4}
	cam_dir := Vec3{0, 0, -1}
	screen_dist := 1.0
	screen_cen := cam_pos.Add(cam_dir.Scale(screen_dist))

	vec_up := Vec3{0, 1, 0}
	cam_right := cam_dir.Cross(vec_up).Normalize()
	cam_up := cam_right.Cross(cam_dir).Normalize()
	screen_pos := screen_cen.Add(cam_right.Scale(nx)).Add(cam_up.Scale(ny))
	light_pos := Vec3{0, 4, 2}

	ray_dir := screen_pos.Sub(cam_pos).Normalize()

	p := cam_pos
	d := 10000.0
	for i := 0; i < MAX_ITERATIONS; i++ {
		d = scene_sd(p, t)
		if d < EPS {
			// fmt.Println("***************************************")
			break
		}
		p = p.Add(ray_dir.Scale(d))
	}

	// fmt.Println("...")

	if d >= EPS {
		return 0.
	}

	// https://en.wikipedia.org/wiki/Phong_reflection_model
	normal := normals_sd(p, t)
	light_dir := light_pos.Sub(p).Normalize()

	dot := light_dir.Dot(normal)
	if dot < 0 {
		dot = 0
	}
	//	dot2 := 0
	Ip := AMBIENT_REFLECTION_CONSTANT*AMBIENT_LIGHT_INTENSITY + DIFFUSE_REFLECTION_CONSTANT*(dot)*DIFFUSE_LIGHT_INTENSITY //+ Ks*pow(dot2,alpha)*Ims
	// fmt.Println("Dot", dot, "Ip", Ip)
	return Ip
}

func step(buffer [H][W]byte, t float64) {
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			color := calc_pixel_color(float64(x), float64(y), t)
			idx := int(10.0*color/256.0) % 10 // len(ASCII) == 10
			// fmt.Println("color", color, "idx", idx)
			BUFFER_1[y][x] = ASCII[idx]
		}
	}
}

func draw(buffer [H][W]byte) {
	for _, row := range buffer {
		fmt.Println(string(row[:]))
	}
}

func CalculateDonut(_ js.Value, args []js.Value) interface{} {
	t := args[0].Float()
	step(BUFFER_1, t)
	// sb := strings.Builder{}
	ss := []string{}
	for _, row := range BUFFER_1 {
		ss = append(ss, string(row[:]))
	}
	return strings.Join(ss, "\n")
}

func main() {
	js.Global().Set("CalculateDonut", js.FuncOf(CalculateDonut))
	select {}
	// t := 0.1
	// for {
	// 	step(BUFFER_1, t)
	// 	draw(BUFFER_1)
	// 	t += .1
	// }
}
