package utils

import "fmt"

func calc_pixel_color(x, y, t float64) float64 {
	nx := 2*x/W - 1
	ny := 2*y/H - 1

	cam_pos := Vec3{0, 0, 4}
	cam_dir := Vec3{0, 0, -1}
	screen_dist := 1.0
	screen_cen := cam_pos.Add(cam_dir.Scale(screen_dist))

	vec_up := Vec3{0, 1, 0}
	cam_right := cam_dir.Cross(vec_up).Normalize()
	cam_up := cam_right.Cross(cam_dir).Normalize()
	screen_pos := screen_cen.Add(cam_right.Scale(nx)).Add(cam_up.Scale(ny))
	light_pos := cam_pos

	ray_dir := screen_pos.Sub(cam_pos).Normalize()

	p := cam_pos
	d := 10000.0 // max distance
	r := GetRotMatAtTime(t)
	for i := 0; i < MAX_ITERATIONS; i++ {
		d = Scene_sd(p, r)
		if d < EPS {
			break
		}
		p = p.Add(ray_dir.Scale(d))
	}

	if d >= EPS {
		return 0.
	}

	// https://en.wikipedia.org/wiki/Phong_reflection_model
	normal := Normals_sd(p, r)
	light_dir := light_pos.Sub(p).Normalize()

	dot := light_dir.Dot(normal)
	if dot < 0 {
		dot = 0
	}
	//	dot2 := 0 // specular component
	Ip := AMBIENT_REFLECTION_CONSTANT*AMBIENT_LIGHT_INTENSITY + DIFFUSE_REFLECTION_CONSTANT*(dot)*DIFFUSE_LIGHT_INTENSITY //+ Ks*pow(dot2,alpha)*Ims
	return Ip
}

//go:export Step
func Step(t float64) {

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			color := calc_pixel_color(float64(x), float64(y), t)
			idx := int(ASCII_LEN*color/256.0) % ASCII_LEN
			// fmt.Println("color", color, "idx", idx)
			BUFFER_1[y][x] = ASCII[idx]
		}
	}
}

func Draw() {
	out := [32 * 65]byte{}
	for i, row := range BUFFER_1 {
		copy(out[i*65:i*65+65], row[:])
	}
	fmt.Println(string(out[:]))
}

//go:export GetBufferAddress
func GetBufferAddress() *[H][W + 1]byte {
	return &BUFFER_1
}

func InitializeBuffer() {
	for i := range BUFFER_1 {
		BUFFER_1[i][W] = '\n'
	}
}
