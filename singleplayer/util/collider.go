package util

import (
	"log/slog"
	"math"
)

type Collider interface {
	Collide(other Collider) (Vector, bool)
}

type ColliderType int

const (
	RectCollider ColliderType = iota
	CircleCollider
)

var dirs = [4][2]float64{
	{1, 1},
	{1, -1},
	{-1, -1},
	{-1, 1},
}

// There are currently two Collider objects: Rect and Circle
// Detecting collision with SAT: https://www.sevenson.com.au/programming/sat/

// Theorectically, this can be done for any convex polygons, but it is pretty hard
// to calculate all the vertices with rotations, so I only have Rect and Circle for now.

type Rect struct {
	Center   *Point
	DimX     float64
	DimY     float64
	Rotation float64
}

func NewRect(center *Point, dimX, dimY, rotation float64) Rect {
	if dimX <= 0 || dimY <= 0 {
		slog.Error("failed to create rect", "dimX", dimX, "dimY", dimY)
		return Rect{}
	}

	return Rect{
		Center:   center,
		DimX:     dimX,
		DimY:     dimY,
		Rotation: rotation,
	}
}

func (r Rect) Collide(other Collider) (Vector, bool) {
	switch other := other.(type) {
	case Rect:
		return r.IntersectRectAndRect(other)
	case Circle:
		return r.IntersectRectAndCircle(other)
	default:
		slog.Error("unrecognized collider type")
		return Vector{}, false
	}
}

// r.GetVertices returns the top right vertice followed by others in a clockwise fashion
func (r Rect) GetVertices() [4]Point {
	var res [4]Point
	c := math.Cos(r.Rotation)
	s := math.Sin(r.Rotation)

	for i := range 4 {
		x := r.Center.X + dirs[i][0]*r.DimX/2
		y := r.Center.Y + dirs[i][1]*r.DimY/2

		res[i] = Point{
			X: x*c - y*s,
			Y: x*s + y*c,
		}
	}

	return res
}

type Circle struct {
	Center *Point
	Radius float64
}

func NewCircle(center *Point, radius float64) Circle {
	return Circle{
		Center: center,
		Radius: radius,
	}
}

func (c Circle) Collide(other Collider) (Vector, bool) {
	switch other := other.(type) {
	case Circle:
		return c.IntersectCircleAndCircle(other)
	case Rect:
		v, yes := other.IntersectRectAndCircle(c)
		return v.ReverseDirection(), yes
	default:
		slog.Error("unrecognized collider type")
		return Vector{}, false
	}
}

// Note for colliders below:
// They only check for collisions that "just" happened, i.e., objects running into each other.
// I am quite sure if an object is fully embedded in another, the returned Vector will be incorrect.
// So I also need to make sure objects' colliders are big enough and their moving speeds are slow enough.

// Use of the returned Vector to separate the two shapes: add it to c, or subtract it from r
func (r Rect) IntersectRectAndCircle(c Circle) (Vector, bool) {
	vertices := r.GetVertices()

	var dist = math.Inf(1)
	var axis Vector

	for _, v := range vertices {
		curr := v.Distance(*c.Center)
		if dist > curr {
			dist = curr
			axis = c.Center.Vector(v)
		}
	}

	var smallestOverlap = math.Inf(1)
	var offsetVector Vector

	for i := 0; i < 4; i++ {
		// get r's 2 projected points onto the axis
		max1 := axis.InnerProduct(Vector(vertices[0]))
		min1 := max1

		for j := 1; j < 4; j++ {
			max1 = math.Max(max1, axis.InnerProduct(Vector(vertices[j])))
			min1 = math.Min(min1, axis.InnerProduct(Vector(vertices[j])))
		}

		// get c's 2 projected points onto the axis
		cent := axis.InnerProduct(Vector(*c.Center))
		max2 := cent + c.Radius
		min2 := cent - c.Radius

		slog.Debug("rect circle projected points", "axis", axis, "min1", min1, "max1", max1, "min2", min2, "max2", max2)

		// check if they overlap
		if min1 < min2 && max1 > min2 { // r is left of other && they overlap
			if smallestOverlap > max1-min2 {
				smallestOverlap = max1 - min2
				offsetVector = axis
			}
		} else if min1 > min2 && max2 > min1 { // r is right of other && they overlap
			if smallestOverlap > max2-min1 {
				smallestOverlap = max2 - min1
				offsetVector = axis
			}
		} else { // found a gap, so they don't overlap
			return offsetVector, false
		}
	}

	offsetVector = offsetVector.Normalize().Scale(smallestOverlap)

	// TODO: there may be a way to deterministically test the direction of the offset vector
	var test Circle = c
	test.Center.Add(offsetVector)
	if _, yes := r.IntersectRectAndCircle(test); yes {
		offsetVector = offsetVector.ReverseDirection()
	}

	slog.Info("rect circle collision", "r", r.Center, "c", c.Center, "offset", smallestOverlap, "vector", offsetVector)

	return offsetVector, true
}

// Use of the returned Vector to separate the two shapes: add it to other, or subtract it from r
func (r Rect) IntersectRectAndRect(other Rect) (Vector, bool) {
	vertices1 := r.GetVertices()
	vertices2 := other.GetVertices()

	var smallestOverlap = math.Inf(1)
	var offsetVector Vector

	for i := 0; i < 4; i++ {
		var axis Vector
		if i != 3 {
			axis = vertices1[i].Vector(vertices1[i+1]).GetPerpendicularVector()
		} else {
			axis = vertices1[i].Vector(vertices1[0]).GetPerpendicularVector()
		}

		// get r's 2 projected points onto the axis
		max1 := axis.InnerProduct(Vector(vertices1[0]))
		min1 := max1

		for j := 1; j < 4; j++ {
			max1 = math.Max(max1, axis.InnerProduct(Vector(vertices1[j])))
			min1 = math.Min(min1, axis.InnerProduct(Vector(vertices1[j])))
		}

		// get other's 2 projected points onto the axis
		max2 := axis.InnerProduct(Vector(vertices2[0]))
		min2 := max2

		for j := 1; j < 4; j++ {
			max2 = math.Max(max2, axis.InnerProduct(Vector(vertices2[j])))
			min2 = math.Min(min2, axis.InnerProduct(Vector(vertices2[j])))
		}

		slog.Debug("rect rect projected points", "axis", axis, "min1", min1, "max1", max1, "min2", min2, "max2", max2)

		// check if they overlap
		if min1 < min2 && max1 > min2 { // r is left of other && they overlap
			if smallestOverlap > max1-min2 {
				smallestOverlap = max1 - min2
				offsetVector = axis
			}
		} else if min1 > min2 && max2 > min1 { // r is right of other && they overlap
			if smallestOverlap > max2-min1 {
				smallestOverlap = max2 - min1
				offsetVector = axis
			}
		} else { // found a gap, so they don't overlap
			return offsetVector, false
		}
	}

	// no gaps found, so the two shapes overlap

	offsetVector = offsetVector.Normalize().Scale(smallestOverlap)

	// TODO: there may be a way to deterministically test the direction of the offset vector
	var test Rect = other
	test.Center.Add(offsetVector)
	if _, yes := r.IntersectRectAndRect(test); yes {
		offsetVector = offsetVector.ReverseDirection()
	}

	slog.Info("rect rect collision", "r", r.Center, "other", other.Center, "offset", smallestOverlap, "vector", offsetVector)

	return offsetVector, true
}

// Use of the returned Vector to separate the two shapes: add it to other, or subtract it from c
func (c Circle) IntersectCircleAndCircle(other Circle) (Vector, bool) {
	dist := other.Center.Distance(*c.Center)
	if dist >= c.Radius+other.Radius {
		return Vector{}, false
	}

	offset := c.Radius + other.Radius - dist

	offsetVector := c.Center.Vector(*other.Center).Normalize().Scale(offset)
	slog.Info("circle circle collision", "c", c.Center, "other", other.Center, "offset", offset, "vector", offsetVector)

	return offsetVector, true
}
