package util

import (
	"log/slog"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// get all the debug logs out
func TestMain(m *testing.M) {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(h))
	os.Exit(m.Run())
}

func TestIntersectRectAndCircle(t *testing.T) {
	// test non-collision

	// test non-collision with rotated rects

	// test collision

	// test collision with rotated rects
}

func TestIntersectRectAndRect(t *testing.T) {
	// test non-collision
	r1 := NewRect(Point{0, 0}, 10, 5, 0)
	r2 := NewRect(Point{7.5, 7.5}, 5, 10, 0)

	_, yes := r1.IntersectRectAndRect(r2)
	assert.False(t, yes)

	// test non-collision with rotated rects
	r1 = NewRect(Point{0, 0}, 2, 2, 0)
	r2 = NewRect(Point{1.7, 1.7}, 1, 2*math.Sqrt(2), -math.Pi/4)

	_, yes = r1.IntersectRectAndRect(r2)
	assert.False(t, yes)

	// test collision
	r1 = NewRect(Point{0, 0}, 10, 5, 0)
	r2 = NewRect(Point{6.5, 6.5}, 5, 10, 0)

	v, yes := r1.IntersectRectAndRect(r2)
	assert.True(t, yes)

	r2.Center.Add(v)
	_, yes = r1.IntersectRectAndRect(r2)
	assert.False(t, yes)

	// test collision with rotated rects
	r1 = NewRect(Point{0, 0}, 2, 2, 0)
	r2 = NewRect(Point{1.4, 1.4}, 1, 2*math.Sqrt(2), -math.Pi/4)

	_, yes = r1.IntersectRectAndRect(r2)
	assert.True(t, yes)

	r2.Center.Add(v)
	_, yes = r1.IntersectRectAndRect(r2)
	assert.False(t, yes)
}

func TestIntersectCircleAndCircle(t *testing.T) {
	// test non-collision
	c1 := NewCircle(Point{10, 10}, 2.49)
	c2 := NewCircle(Point{13, 14}, 2.49)

	_, yes := c1.IntersectCircleAndCircle(c2)
	assert.False(t, yes)

	// test collision
	c1 = NewCircle(Point{10, 10}, 2.6)
	c2 = NewCircle(Point{13, 14}, 2.6)

	v, yes := c1.IntersectCircleAndCircle(c2)
	assert.True(t, yes)

	c2.Center.Add(v)
	_, yes = c1.IntersectCircleAndCircle(c2)
	assert.False(t, yes)
}
