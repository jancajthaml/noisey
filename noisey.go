/*
Package noisey is a library that implements coherent noise algorithms.

The selection is currently:

	* 2D/3D Perlin noise (64bit)
	* 2D/3D OpenSimplex noise (64bit)

The sources above can be combined with different generators and modifiers
like the following:

	* FBMGenerator2D - fractal Brownian Motion
	* Select2D - choose from source A or B depending on control source
	* Scale2D - modify output by multiplying by a scale and adding a bias constant


Once the noise generators have been set up, a Builder2D object can be created
to map a region of noise into a float64 array.

An interface called 'RandomSource' is also exported so that a client can implement
a different random number generator and pass it to the noise generators.

Sample programs can be found in the 'examples' directory.

*/
package noisey

// RandomSource is a generic interface for a random number generator
// allowing the user to use the built-in RNG or a custom one that implements
// this interface.
type RandomSource interface {
	Float64() float64
	Perm(int) []int
}

// NoiseyGet2D is an interface defining how the modules types get noise from a source.
type NoiseyGet2D interface {
	Get2D(float64, float64) float64
}

// NoiseyGet3D is an interface defining how the modules types get noise from a source.
type NoiseyGet3D interface {
	Get3D(float64, float64, float64) float64
}

// Vec2f is a simple 2D vector of 64 bit floats
type Vec2f struct {
	X, Y float64
}

// Vec3f is a simple 3D vector of 64 bit floats
type Vec3f struct {
	X, Y, Z float64
}

// Vec4f is a simple 4D vector of 64 bit floats
type Vec4f struct {
	X, Y, Z, W float64
}

// Vec2i is a simple 3D vector of ints
type Vec2i struct {
	X, Y int
}

// Vec3i is a simple 3D vector of ints
type Vec3i struct {
	X, Y, Z int
}

func calcCubicSCurve(v float64) float64 {
	return v * v * (3 - 2*v)
}

func calcQuinticSCurve(v float64) float64 {
	v3 := v * v * v
	v4 := v3 * v
	v5 := v4 * v
	return (6.0 * v5) - (15.0 * v4) + (10.0 * v3)
}

func lerp(a, b, v float64) float64 {
	return a*(1-v) + b*v
}
