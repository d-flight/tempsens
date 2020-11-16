package data

import "fmt"

// Color ...
type Color []byte

// NewColor ...
func NewColor(r byte, g byte, b byte) *Color {
	return &Color{r, g, b}
}

// Red ...
func (c Color) Red() byte {
	return c[0]
}

// Green ...
func (c Color) Green() byte {
	return c[1]
}

// Blue ...
func (c Color) Blue() byte {
	return c[2]
}

func (c *Color) String() string {
	return fmt.Sprintf("%v|%v|%v", c.Red(), c.Green(), c.Blue())
}

func Blue() *Color   { return NewColor(0, 0, 0xff) }
func Orange() *Color { return NewColor(0xff, 0xa5, 0) }
func Red() *Color    { return NewColor(0xff, 0, 0) }
func None() *Color   { return NewColor(0, 0, 0) }
