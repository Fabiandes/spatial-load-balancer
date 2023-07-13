package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"math"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/fabiandes/slb/demo/client"
	"github.com/fabiandes/slb/demo/simulation/entity"
)

const EntitySize = 2

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a window.
	w := app.NewWindow(app.Size(1000, 1000))

	// Create Scene
	s := &client.Scene{
		Entities: []*entity.Entity{},
		Window:   w,
	}

	// Setup client
	c := client.New(s)
	go func() {
		if err := c.Listen(ctx, "localhost:8080"); err != nil {
			fmt.Printf("An error occurred while listening for updates: %v\n", err)
			cancel()
		}
	}()

	// Setup window and render loop
	go func() {
		if err := run(ctx, w, s); err != nil {
			fmt.Printf("An error occurred while running window: %v\n", err)
			cancel()
		}
	}()

	app.Main()
}

func run(ctx context.Context, w *app.Window, s *client.Scene) error {
	var ops op.Ops
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancellation")
		default:
			e := <-w.Events()
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				for _, e := range s.Entities {
					DrawEntity(gtx, image.Point{X: int(e.Transform.Position.X), Y: int(e.Transform.Position.Y)})
				}

				e.Frame(gtx.Ops)
			}
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
