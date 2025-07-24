package util

import (
	"log/slog"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Get all the debug logs out!
func TestMain(m *testing.M) {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(h))
	os.Exit(m.Run())
}

func TestGetVertices(t *testing.T) {
	// unrotated
	r := NewRect(&Point{0, 0}, 2, 2, 0)
	v := r.GetVertices()
	assert.Equal(t, [4]Point{
		{1, 1},
		{-1, 1},
		{-1, -1},
		{1, -1},
	}, v)

	// rotate by 90 degrees
	r = NewRect(&Point{0, 0}, 2, 2, 90.0*math.Pi/180)
	v = r.GetVertices()
	compareVertices(t, [4]Point{
		{-1, 1},
		{-1, -1},
		{1, -1},
		{1, 1},
	}, v)

	// rotate by 45 degrees
	r = NewRect(&Point{0, 0}, 2, 2, 45.0*math.Pi/180)
	v = r.GetVertices()
	sqrt2 := math.Sqrt(2)
	compareVertices(t, [4]Point{
		{0, sqrt2},
		{-sqrt2, 0},
		{0, -sqrt2},
		{sqrt2, 0},
	}, v)

}

// compareVertices is a helper to accommodate losses of precision caused by integer <-> float conversion
func compareVertices(t *testing.T, expected, actual [4]Point) {
	t.Helper()

	for i := range 4 {
		xDiff := math.Abs(expected[i].X - actual[i].X)
		assert.LessOrEqual(t, xDiff, 0.01)

		yDiff := math.Abs(expected[i].Y - actual[i].Y)
		assert.LessOrEqual(t, yDiff, 0.01)
	}
}

func TestIntersectRectAndCircle(t *testing.T) {
	// test non-collision
	r := NewRect(&Point{0, 0}, 4, 2, 0)
	c := NewCircle(&Point{4, 3}, 1.9*math.Sqrt(2))

	_, yes := r.IntersectRectAndCircle(c)
	assert.False(t, yes)

	r = NewRect(&Point{0, 0}, 4, 2, 0)
	c = NewCircle(&Point{0, 3}, 2)

	_, yes = r.IntersectRectAndCircle(c)
	assert.False(t, yes)

	// test non-collision with rotated rects
	r = NewRect(&Point{0, 0}, math.Sqrt(2), math.Sqrt(2), math.Pi/4)
	c = NewCircle(&Point{2, 2}, 1.4*math.Sqrt(2))

	_, yes = r.IntersectRectAndCircle(c)
	assert.False(t, yes)

	// test collision
	r = NewRect(&Point{0, 0}, 4, 2, 0)
	c = NewCircle(&Point{0, 3}, 2.1)

	v, yes := r.IntersectRectAndCircle(c)
	assert.True(t, yes)

	c.Center.Add(v)
	_, yes = r.IntersectRectAndCircle(c)
	assert.False(t, yes)

	// test collision with rotated rects
	r = NewRect(&Point{0, 0}, math.Sqrt(2), math.Sqrt(2), math.Pi/4)
	c = NewCircle(&Point{2, 2}, 2*math.Sqrt(2))

	v, yes = r.IntersectRectAndCircle(c)
	assert.True(t, yes)

	c.Center.Add(v)
	_, yes = r.IntersectRectAndCircle(c)
	assert.False(t, yes)
}

func TestIntersectRectAndRect(t *testing.T) {
	// test non-collision
	r1 := NewRect(&Point{0, 0}, 10, 5, 0)
	r2 := NewRect(&Point{7.5, 7.5}, 5, 10, 0)

	_, yes := r1.IntersectRectAndRect(r2)
	assert.False(t, yes)

	// test non-collision with rotated rects
	r1 = NewRect(&Point{0, 0}, 2, 2, 0)
	r2 = NewRect(&Point{3, 3}, 1, 2*math.Sqrt(2), -math.Pi/4)

	_, yes = r1.IntersectRectAndRect(r2)
	assert.False(t, yes)

	// test collision
	r1 = NewRect(&Point{0, 0}, 10, 5, 0)
	r2 = NewRect(&Point{6.5, 6.5}, 5, 10, 0)

	v, yes := r1.IntersectRectAndRect(r2)
	assert.True(t, yes)

	r2.Center.Add(v)
	_, yes = r1.IntersectRectAndRect(r2)
	assert.False(t, yes)

	// test collision with rotated rects
	r1 = NewRect(&Point{0, 0}, 2, 2, 0)
	r2 = NewRect(&Point{1.4, 1.4}, 1, 2*math.Sqrt(2), -math.Pi/4)

	v, yes = r1.IntersectRectAndRect(r2)
	assert.True(t, yes)

	r2.Center.Add(v)
	_, yes = r1.IntersectRectAndRect(r2)
	assert.False(t, yes)
}

func TestIntersectCircleAndCircle(t *testing.T) {
	// test non-collision
	c1 := NewCircle(&Point{10, 10}, 2.49)
	c2 := NewCircle(&Point{13, 14}, 2.49)

	_, yes := c1.IntersectCircleAndCircle(c2)
	assert.False(t, yes)

	// test collision
	c1 = NewCircle(&Point{10, 10}, 2.6)
	c2 = NewCircle(&Point{13, 14}, 2.6)

	v, yes := c1.IntersectCircleAndCircle(c2)
	assert.True(t, yes)

	c2.Center.Add(v)
	_, yes = c1.IntersectCircleAndCircle(c2)
	assert.False(t, yes)
}
