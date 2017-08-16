package drawmap

import (
   "image"
   "image/color"
   "image/draw"
   "image/png"
   "os"
   "github.com/llgcode/draw2d/draw2dimg"
   "./imgtext"
   "strconv"
   . "../structs"
)

var (
	background  color.Color = color.RGBA{151, 161, 178, 255}
	circle      color.Color = color.RGBA{102, 140, 204, 255}
	line        color.Color = color.RGBA{178, 101, 89, 255}

	max_x = 4096
	max_y = 4096
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

func buildMap(path []int, coords []Coord) *image.RGBA{

	surface := Rectangle{ length:max_y, width:max_x } // Surface to draw on
	rpainter := Rectangle{ length:max_y, width:max_x } // Colored Mask Layer

	m  := surface.drawShape()
	cr := rpainter.drawShape()

	draw.Draw(m, m.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)
	draw.Draw(cr, cr.Bounds(), &image.Uniform{circle}, image.ZP, draw.Src)



    for i := 0; i < len(path); i++{
        pointA := path[i]
        pointB := path[(i+1) % len(path)]

        x1 := coords[pointA].X + 30
        y1 := coords[pointA].Y + 30

        x2 := coords[pointB].X + 30
        y2 := coords[pointB].Y + 30

        drawlines([4]float64{x1, y1, x2, y2}, m)
    }


    for i := 0; i < len(coords); i++{
        x := int(coords[i].X) + 30
        y := int(coords[i].Y) + 30
        draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{x, y}, 30}, image.ZP, draw.Over)
        if(false){
            imgtext.AddLabel(m, x, y, strconv.Itoa(i))
        }
    }

	return m
}

func DrawMap(path []int, coords []Coord, num int) {

	m := buildMap(path, coords);

	w, _ := os.Create("graphs/tourGraph" + strconv.Itoa(num) + ".png")
	defer w.Close()
	png.Encode(w, m) //Encode writes the Image m to w in PNG format.
}
