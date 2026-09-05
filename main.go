package main

import (
	"image"
	"log"
	"math"
	"math/cmplx"
	"math/rand"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ojrac/opensimplex-go"
)

const (
	screenWidth  = 640
	screenHeight = 640
	maxIt        = 32
)

var (
	palette [maxIt]byte
)

func init() {
	for i := range palette {
		palette[i] = byte(math.Sqrt(float64(i)/float64(len(palette))) * 0x80)
	}
}

func color(it int) (r, g, b byte) {
	if it == maxIt {
		return 0xff, 0xff, 0xff
	}
	c := palette[it]
	return c, c, c
}

type Game struct {
	debugui    debugui.DebugUI
	screen     *ebiten.Image
	pixels     []byte
	dropdown   []string
	current    int
	parametr   int
	noisescale float64
}

func NewGame() *Game {
	g := &Game{
		current:    0,
		dropdown:   []string{"Mandelbrot", "Sierpinski triangle", "Perlin noise"},
		screen:     ebiten.NewImage(screenWidth, screenHeight),
		pixels:     make([]byte, screenWidth*screenHeight*4),
		parametr:   2,
		noisescale: 0.05,
	}
	g.mandelbrot(-0.75, 0.0, 2.7, g.parametr)
	return g
}

func (gm *Game) mandelbrot(centerX, centerY, size float64, pf int) {
	for j := range screenHeight {
		for i := range screenWidth {
			x := float64(i)*size/screenWidth - size/2 + centerX
			y := (screenHeight-float64(j))*size/screenHeight - size/2 + centerY

			c := complex(x, y)
			z := complex(0, 0)
			it := 0

			for ; it < maxIt; it++ {
				z = cmplx.Pow(z, complex(float64(pf), 0)) + c
				if cmplx.Abs(z) > 2 {
					break
				}
			}

			r, g, b := color(it)
			p := 4 * (i + j*screenWidth)
			gm.pixels[p] = r
			gm.pixels[p+1] = g
			gm.pixels[p+2] = b
			gm.pixels[p+3] = 0xff
		}
	}
	gm.screen.WritePixels(gm.pixels)
}
func (gm *Game) sierpinski_triangle(centerX, centerY, size float64) {
	for j := range screenHeight {
		for i := range screenWidth {
			x := float64(i)*size/screenWidth - size/2 + centerX
			y := (screenHeight-float64(j))*size/screenHeight - size/2 + centerY

			const scale = 512.0

			iy := int((size/2 - y) * (scale / size))
			ix := int((x+size/2)*(scale/size)) - (int(scale)-iy)/2

			r, g, b := byte(0), byte(0), byte(0)
			if iy >= 0 && iy < int(scale) && ix >= 0 && ix <= iy {
				if (ix & (iy - ix)) == 0 {
					r, g, b = 0xff, 0xff, 0xff
				}
			}
			p := 4 * (i + j*screenWidth)
			gm.pixels[p] = r
			gm.pixels[p+1] = g
			gm.pixels[p+2] = b
			gm.pixels[p+3] = 0xff
		}
	}
	gm.screen.WritePixels(gm.pixels)
}

func (gm *Game) perlin_noise(scale float64) {
	noise := opensimplex.New(rand.Int63())

	for j := range screenHeight {
		for i := range screenWidth {
			nx := float64(i) * scale
			ny := float64(j) * scale

			val := noise.Eval2(nx, ny)

			c := uint8((val + 1.0) * 127.5)

			p := 4 * (i + j*screenWidth)
			gm.pixels[p] = c
			gm.pixels[p+1] = c
			gm.pixels[p+2] = c
			gm.pixels[p+3] = 0xff
		}
	}
	gm.screen.WritePixels(gm.pixels)
}

func (g *Game) Update() error {
	if _, err := g.debugui.Update(func(ctx *debugui.Context) error {
		ctx.Window("Settings", image.Rect(10, 10, 320, 240), func(layout debugui.ContainerLayout) {

			ctx.Text("Select fractal/noise")
			ctx.Dropdown(&g.current, g.dropdown)

			switch g.current {
			case 0:
				ctx.Text("Fractal pow:")
				ctx.NumberField(&g.parametr, 1)
			case 1:
			case 2:
				ctx.Text("Noise scale")
				ctx.NumberFieldF(&g.noisescale, 0.01, 3)
			}
			ctx.Button("Update").On(func() {
				switch g.current {
				case 0:
					g.screen.Clear()
					if g.parametr > 2 {
						g.mandelbrot(0, 0.0, 2.7, g.parametr)
					} else {
						g.mandelbrot(-0.75, 0.0, 2.7, g.parametr)
					}
				case 1:
					g.screen.Clear()
					g.sierpinski_triangle(0.0, 0.0, 2.7)
				case 2:
					g.screen.Clear()
					g.perlin_noise(g.noisescale)
				}
			})

		})
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.screen, nil)
	g.debugui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Fractalator")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
