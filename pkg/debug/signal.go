package debug

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/jnafolayan/sip/internal/imageutils"
	"github.com/jnafolayan/sip/pkg/signal"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func DrawSignal2D(s signal.Signal2D, region image.Rectangle, fileName string) error {
	cellSize := 30

	img := image.NewRGBA(image.Rect(0, 0, region.Dx()*cellSize, region.Dy()*cellSize))
	fg, bg := image.Black, image.White
	d := &font.Drawer{
		Dst:  img,
		Src:  fg,
		Face: basicfont.Face7x13,
	}

	// Draw background
	draw.Draw(img, img.Bounds(), bg, image.Point{}, draw.Src)

	halfW := d.MeasureString("255").Ceil() / 2
	marginX := cellSize/2 - halfW
	marginY := cellSize / 2

	for y := 0; y < region.Dy(); y++ {
		for x := 0; x < region.Dx(); x++ {
			v := int(s[y+region.Min.Y][x+region.Min.X])
			d.Dot = fixed.P(x*cellSize+marginX, y*cellSize+marginY)
			d.DrawString(fmt.Sprintf("%3d", v))
		}
	}

	dest := fmt.Sprintf("dump/%s", fileName)
	err := imageutils.SaveImage(dest, img)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote signal image to %s\n", dest)
	return nil
}
