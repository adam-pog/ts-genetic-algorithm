package main

import (
   "image"
   "image/color"
   "image/draw"
   "image/png"
   "github.com/llgcode/draw2d/draw2dimg"

   "log"
   "os"
   "os/exec"
)

// Sets Colors for easy use
var (
	teal  color.Color = color.RGBA{0, 200, 200, 255}
	red   color.Color = color.RGBA{200, 30, 30, 255}

	max_x = 1024
	max_y = 1024
)

// Interface for displaying Shapes

type displayShape interface {
	drawShape() *image.RGBA
}

// Struct for rectangle

type Rectangle struct {
	p image.Point
	length, width int
}

// Struct for circle

type Circle struct {
	p image.Point
	r int
}


// The following is taken from:
// http://blog.golang.org/go-imagedraw-package
// Excellent use of masking!!!

func (c *Circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *Circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *Circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}



// Rectangle Draw shape function

func (r Rectangle) drawShape() *image.RGBA{
	return image.NewRGBA(image.Rect(0, 0, r.length, r.width))
}

// Draws lines between a source and destination,
// this is a bit ugly, but should be fairly straight forward

func drawlines(pos [4]float64, m *image.RGBA){
	gc := draw2dimg.NewGraphicContext(m)
	gc.MoveTo(pos[0], pos[1])
	gc.LineTo(pos[2], pos[3])
	gc.SetStrokeColor(red)
	gc.SetLineWidth(3)
	gc.Stroke()
}

// Currently a static map

func buildMap() *image.RGBA{

	surface := Rectangle{ length:max_y, width:max_x } // Surface to draw on
	rpainter := Rectangle{ length:max_y, width:max_x } // Colored Mask Layer

	m  := surface.drawShape()
	cr := rpainter.drawShape()

	draw.Draw(m, m.Bounds(), &image.Uniform{teal}, image.ZP, draw.Src)
	draw.Draw(cr, cr.Bounds(), &image.Uniform{red}, image.ZP, draw.Src)

	/** Generates two circles **/
	draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{500, 500}, 20}, image.ZP, draw.Over)
	draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{0, 0}, 35}, image.ZP, draw.Over)
	draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{300, 165}, 55}, image.ZP, draw.Over)
	draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{200, 405}, 47}, image.ZP, draw.Over)

	/** Draws Line Between Circles **/
	drawlines([4]float64{500, 500, 0, 0}, m)
	drawlines([4]float64{0, 0, 300, 165}, m)
	drawlines([4]float64{0, 0, 200, 405}, m)

	return m
}

// for OS X(darwin)

func Show(name string) {
	command := "open"
	arg1 := "-a"
	arg2 := "/Applications/Preview.app"
	cmd := exec.Command(command, arg1, arg2, name)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	m := buildMap();

	w, _ := os.Create("blogmap.png")
	defer w.Close()
	png.Encode(w, m) //Encode writes the Image m to w in PNG format.

	Show(w.Name())
}
