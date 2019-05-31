package geom

import (
	"fmt"
	"math"
	"strconv"
)

//////////////////////////////////////////////////////////////
type Point struct {
	X float64 `json:"x" msgpack:"x"`
	Y float64 `json:"y" msgpack:"y"`
}

func (p Point) String() string {
	return "(" + strconv.FormatFloat(p.X, 'f', 2, 64) + "," + strconv.FormatFloat(p.Y, 'f', 2, 64) + ")"
}

// Add returns the vector p+q.
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p Point) Mul(k float64) Point {
	return Point{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p Point) Div(k float64) Point {
	return Point{p.X / k, p.Y / k}
}

// In reports whether p is in r.
func (p Point) In(r Rectangle) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

// Point in Section.
func (p Point) Sec4() int {
	if p.X < 50 && p.Y < 50 {
		return 1
	} else if p.X < 50 && p.Y >= 50 {
		return 2
	} else if p.X >= 50 && p.Y >= 50 {
		return 3
	} else {
		return 4
	}
}

func (p Point) Angle(q Point) float64 {
	return math.Atan2(q.Y-p.Y, q.X-p.X) / math.Pi * 180
}

//////////////////////////////////////////////////////////////////////////////////
type Line struct {
	Start Point `json:"start" msgpack:"start"` // start point
	End   Point `json:"end" msgpack:"end"`     // end point
}

func (l Line) String() string {
	return "[" + l.Start.String() + "," + l.End.String() + "]"
}

// 线段是否相交
func (l Line) IntersectLine(m Line) bool {
	if math.Min(l.Start.X, l.End.X) <= math.Max(m.Start.X, m.End.X) && math.Min(m.Start.X, m.End.X) <= math.Max(l.Start.X, l.End.X) &&
		math.Min(l.Start.Y, l.End.Y) <= math.Max(m.Start.Y, m.End.Y) && math.Min(m.Start.Y, m.End.Y) <= math.Max(l.Start.Y, l.End.Y) {
		u := (m.Start.X-l.Start.X)*(l.End.Y-l.Start.Y) - (l.End.X-l.Start.X)*(m.Start.Y-m.Start.Y)
		v := (m.End.X-l.Start.X)*(l.End.Y-l.Start.Y) - (l.End.X-l.Start.X)*(m.End.Y-l.Start.Y)
		w := (l.Start.X-m.Start.X)*(m.End.Y-m.Start.Y) - (m.End.X-m.Start.X)*(l.Start.Y-m.Start.Y)
		z := (l.End.X-m.Start.X)*(m.End.Y-m.Start.Y) - (m.End.X-m.Start.X)*(l.End.Y-m.Start.Y)
		return u*v <= 0.00000001 && w*z <= 0.00000001
	}
	return false
}

////////////////////////////////////////////////////////////////////////////
// 矩形
type Rectangle struct {
	Min Point `json:"min" msgpack:"min"`
	Max Point `json:"max" msgpack:"max"`
}

func (r Rectangle) String() string {
	return "{" + r.Min.String() + "," + r.Max.String() + "}"
}

// 获取矩形中点坐标
func (r Rectangle) CenterPoint() Point {
	x := (r.Max.X-r.Min.X)/2 + r.Min.X
	y := (r.Max.Y-r.Min.Y)/2 + r.Min.Y
	x, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", x), 64)
	y, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", y), 64)
	return Point{x, y}
}

// 获取底部线段
func (r Rectangle) BottomLine() Line {
	return Line{
		Start: Point{r.Min.X, r.Max.Y},
		End:   Point{r.Max.X, r.Max.Y},
	}
}

// 判断点在矩形内
func (r Rectangle) PointAt(point Point) bool {
	if point.X >= r.Min.X && point.X <= r.Max.X && point.Y >= r.Min.Y && point.Y <= r.Max.Y {
		return true
	}
	return false
}

// 判断线段在矩形内
func (r Rectangle) LineIn(line Line) bool {
	if r.PointAt(line.Start) && r.PointAt(line.End) {
		return true
	}
	return false
}

//////////////////////////////////////////////////////////////////////////////
// 多边形
type Polygon struct {
	Points []Point `json:"points" msgpack:"points"`
}

func (p Polygon) String() string {
	var str string
	for _, v := range p.Points {
		str = str + v.String() + ","
	}
	return "{" + str + "}"
}

// 判断点在多边形内
func (p Polygon) PointAt(point Point) bool {
	x := point.X
	y := point.Y
	sz := len(p.Points)
	isIn := false

	for i := 0; i < sz; i++ {
		j := i - 1
		if i == 0 {
			j = sz - 1
		}
		vi := p.Points[i]
		vj := p.Points[j]

		xMin := vi.X
		xMax := vj.X
		if xMin > xMax {
			xMin, xMax = xMax, xMin
		}
		yMin := vi.Y
		yMax := vj.Y
		if yMin > yMax {
			yMin, yMax = yMax, yMin
		}

		if eq(vj.Y, vi.Y) {
			if eq(y, vi.Y) && lte(xMin, x) && lte(x, xMax) {
				return true
			}
			continue
		}

		xt := (vj.X-vi.X)*(y-vi.Y)/(vj.Y-vi.Y) + vi.X
		if eq(xt, x) && lte(yMin, y) && lte(y, yMax) {
			return true
		}
		if lt(x, xt) && lte(yMin, y) && lt(y, yMax) {
			isIn = !isIn
		}

	}
	return isIn
}

// 判断线段是否与图形相交
func (p Polygon) LineIntersect(line Line) bool {
	if p.PointAt(line.Start) || p.PointAt(line.End) {
		return true
	}
	size := len(p.Points)
	for i := 0; i < size; i++ {
		var l Line
		if i == 0 {
			l = Line{Start: p.Points[size-1], End: p.Points[i]}
		} else if i == size-1 {
			l = Line{Start: p.Points[i], End: p.Points[0]}
		} else {
			l = Line{Start: p.Points[i], End: p.Points[i+1]}
		}
		if line.IntersectLine(l) {
			return true
		}
	}
	return false
}

///////////////////////////////////////////////////////////////
func eq(x float64, y float64) bool {
	v := x - y
	const delta float64 = 1e-6
	if v < delta && v > -delta {
		return true
	}
	return false

}

func lt(x float64, y float64) bool {
	if eq(x, y) {
		return false
	}
	return x < y
}

func lte(x float64, y float64) bool {
	if eq(x, y) {
		return true
	}
	return x < y
}
