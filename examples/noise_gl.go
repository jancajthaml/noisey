/* Copyright 2015, Timothy Bogdala <tdb@animal-machine.com>
   See the LICENSE file for more details. */

package main

/*

This is a test module that does the following:

1) Creates an OpenGL window
2) Creates an RGB texture from perlin noise
3) Displays the noise as a texture on a plane in the window

It requires the GLFW3 and GLEW libraries as well as the Go wrappers
for them: go-gl/gl and go-gl/glfw3.

Basic build instructions are:

	go get github.com/go-gl/gl/v3.3-core/gl
	go get github.com/go-gl/glfw/v3.1/glfw
	go get github.com/tbogdala/noisey
	cd $GOHOME/src/github.com/tbogdala/noisey/examples
	./bulid.sh
	./noise_gl

Hit `esc` to quit the program.

*/

import (
	gl "github.com/go-gl/gl/v3.3-core/gl"
	glfw "github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/noisey"
	"math"
	"math/rand"
)

var (
	app   *ExampleApp
	plane *Renderable
)

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

// createTextureFromRGB makes an OpenGL texture and buffers the RGB data into it
func createTextureFromRGB(rgb []byte, imageSize int32) (tex uint32) {
	gl.GenTextures(1, &tex)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, imageSize, imageSize, 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(rgb))

	return tex
}

func generateNoiseImage(imageSize int, r noisey.RandomSource) []byte {
	// create a new perlin noise generator
	perlin := noisey.NewPerlinGenerator(r)

	// create the fractal Brownian motion generator based on perlin
	fbmPerlin := noisey.NewFBMGenerator2D(&perlin, 1, 0.5, 2.0, 1.0)

	// make an pixel image by calculating random noise and creating
	// an RGB byte triplet array based off the scaled noise value
	colors := make([]byte, imageSize*imageSize*3)
	for y := 0; y < imageSize; y++ {
		for x := 0; x < imageSize; x++ {
			v := fbmPerlin.Get2D(float64(x)*0.1, float64(y)*0.1)
			b := byte(math.Floor((v*0.5 + 0.5) * 255)) // normalize 0..1 then scale by 255
			colorIndex := y*imageSize*3 + x*3
			colors[colorIndex] = b
			colors[colorIndex+1] = b
			colors[colorIndex+2] = b
		}
	}

	return colors
}

func renderCallback(delta float64) {
	gl.Viewport(0, 0, int32(app.Width), int32(app.Height))
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// make the projection and view matrixes
	projection := mgl.Ident4()
	view := mgl.Ident4()

	plane.Draw(projection, view)
}

func main() {
	app = NewApp()
	app.InitGraphics("Noisey Perlin Test", 768, 768)
	app.SetKeyCallback(keyCallback)
	app.OnRender = renderCallback

	// compile our shader
	var err error
	textureShader, err := LoadShaderProgram(UnlitTextureVertShader, UnlitTextureFragShader)
	if err != nil {
		panic("Failed to compile the shader! " + err.Error())
	}

	// create the plane to draw as a test
	plane = CreatePlaneXY(-0.75, -0.75, 0.75, 0.75, 1.0)
	plane.Shader = textureShader

	// generate the noise and make a image
	r := rand.New(rand.NewSource(int64(1)))
	randomPixels := generateNoiseImage(512, r)
	noiseTex := createTextureFromRGB(randomPixels, 512)
	plane.Tex0 = noiseTex

	app.RenderLoop()
}
