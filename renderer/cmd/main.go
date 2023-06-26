package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/fabiandes/spatial-load-balancer/renderer"
)

const (
	EntitySize = 10
)

func main() {
	w := app.NewWindow(app.Size(1000, 1000))

	// Create and initialize the scene.
	s := renderer.NewScene(w)
	go func() {
		if err := s.Initialize(); err != nil {
			log.Fatalf("An error occurred while managing the scene: %v", err)
		}
	}()

	// Setup window and render loop
	go func() {
		err := run(w, s)
		if err != nil {
			log.Fatalf("An error occurred while running window: %v", err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func run(w *app.Window, s *renderer.Scene) error {
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			fmt.Printf("Drawing %d entities.\n", len(s.Entities))
			gtx := layout.NewContext(&ops, e)
			for _, e := range s.Entities {
				DrawEntity(gtx, image.Point{X: int(e.Transform.Position.X), Y: int(e.Transform.Position.Y)})
			}

			e.Frame(gtx.Ops)
		}
	}
}

func DrawEntity(gtx layout.Context, position image.Point) {
	x := int(math.Round(float64(float32(position.X) * gtx.Metric.PxPerDp)))
	y := int(math.Round(float64(float32(position.Y) * gtx.Metric.PxPerDp)))

	entityRadius := int(math.Round(float64(EntitySize / 2 * gtx.Metric.PxPerDp)))
	defer clip.Ellipse{
		Min: image.Point{
			X: x - entityRadius,
			Y: y - entityRadius,
		},
		Max: image.Point{
			X: x + entityRadius,
			Y: y + entityRadius,
		},
	}.Push(gtx.Ops).Pop()

	paint.ColorOp{Color: color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: 0xFF}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

}
