package main

import (
   "image"
   "image/color"
   "image/draw"
   "image/png"
   "log"
   "os"
   "os/exec"
   "./imgtext"
   "github.com/llgcode/draw2d/draw2dimg"
)

var (
	background  color.Color = color.RGBA{151, 161, 178, 255}
	circle      color.Color = color.RGBA{102, 140, 204, 255}
	line        color.Color = color.RGBA{178, 101, 89, 255}

	max_x = 2048
	max_y = 2048
)


type displayShape interface {
	drawShape() *image.RGBA
}


type Rectangle struct {
	p image.Point
	length, width int
}


type Circle struct {
	p image.Point
	r int
}


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


func (r Rectangle) drawShape() *image.RGBA{
	return image.NewRGBA(image.Rect(0, 0, r.length, r.width))
}

func drawlines(pos [4]float64, m *image.RGBA){
	gc := draw2dimg.NewGraphicContext(m)
	gc.MoveTo(pos[0], pos[1])
	gc.LineTo(pos[2], pos[3])
	gc.SetStrokeColor(line)
	gc.SetLineWidth(3)
	gc.Stroke()
}

func buildMap() *image.RGBA{

	surface := Rectangle{ length:max_y, width:max_x } // Surface to draw on
	rpainter := Rectangle{ length:max_y, width:max_x } // Colored Mask Layer

	m  := surface.drawShape()
	cr := rpainter.drawShape()

	draw.Draw(m, m.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)
	draw.Draw(cr, cr.Bounds(), &image.Uniform{circle}, image.ZP, draw.Src)

    drawlines([4]float64{1000, 1000, 30, 30}, m)
	drawlines([4]float64{500, 500, 30, 30}, m)
	drawlines([4]float64{30, 30, 300, 165}, m)
	drawlines([4]float64{30, 30, 200, 405}, m)

	draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{1000, 1000}, 30}, image.ZP, draw.Over)
	draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{500, 500}, 30}, image.ZP, draw.Over)
	draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{30, 30}, 30}, image.ZP, draw.Over)
	draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{300, 165}, 30}, image.ZP, draw.Over)
	draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{200, 405}, 30}, image.ZP, draw.Over)

    imgtext.AddLabel(m, 500, 500, "1")
    imgtext.AddLabel(m, 1000, 1000, "5")
    imgtext.AddLabel(m, 30, 30, "2")
    imgtext.AddLabel(m, 300, 165, "3")
    imgtext.AddLabel(m, 200, 405, "4")


	return m
}

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

	w, '_' := os.Create("blogmap.png")
	defer w.Close()
	png.Encode(w, m) //Encode writes the Image m to w in PNG format.

	Show(w.Name())
}
