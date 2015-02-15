package noisey

/* Copyright 2015, Timothy Bogdala <tdb@animal-machine.com>
See the LICENSE file for more details. */

import "math"

// Scale2D is a module that uses gets the noise from Source, scales
// it and then adds a bias.
type Scale2D struct {
  // the noise that the select module uses
  Source  NoiseyGet2D

  // what to scale the noise value from Source by
  Scale float64

  // the const value to add to the scaled noise value
  Bias float64

  // the minimum value to return
  Min float64

  // the maximum value to return
  Max float64
}

// Scale2D creates a new scale 2d module.
func NewScale2D(src NoiseyGet2D, scale float64, bias float64, min float64, max float64) (scales Scale2D) {
  scales.Source = src
  scales.Scale = scale
  scales.Bias = bias
  scales.Min = min
  scales.Max = max
  return
}

// Get2D calculates the noise value scaling it by Scale and adding Bias
func (scales *Scale2D) Get2D(x float64, y float64) (v float64) {
  v = scales.Source.Get2D(x, y)
  v *= scales.Scale
  v += scales.Bias
  v = math.Max(scales.Min, v)
  v = math.Min(scales.Max, v)
  return v
}
