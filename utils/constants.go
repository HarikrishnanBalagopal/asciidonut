package utils

const (
	N              = 32
	W              = N * 2
	H              = N
	MAX_ITERATIONS = N
	EPS            = 0.001
	// https://en.wikipedia.org/wiki/Phong_reflection_model
	AMBIENT_REFLECTION_CONSTANT = 0.1
	AMBIENT_LIGHT_INTENSITY     = 255.0
	DIFFUSE_REFLECTION_CONSTANT = 0.9
	DIFFUSE_LIGHT_INTENSITY     = 255.0
	// SPECULAR_REFLECTION_CONSTANT = 0.1
	// MATERIAL_SHININESS_CONSTANT = 0.1
)

var (
	// ASCII (old set of characters), these are my own
	// ASCII    = []byte{' ', '.', ',', '-', '+', '*', 'o', '0', '@', '#'}
	// ASCII (new set of characters) taken from https://www.a1k0n.net/2011/07/20/donut-math.html
	//           output[xp, yp] = ".,-~:;=!*#$@"[luminance_index];
	ASCII    = []byte{' ', '.', ',', '-', '~', ':', ';', '=', '!', '*', '#', '$', '@'}
	BUFFER_1 = [H][W + 1]byte{}
)

const ASCII_LEN = 13 // (old set has length 10)
