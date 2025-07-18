//go:generate mockgen -source=color.go -destination=color_mock.go -package=proxy

package proxy

import (
	"github.com/fatih/color"
)

// Color provides a proxy interface for color package functionality
type Color interface {
	New(value ...color.Attribute) ColorFunc
	Red(format string, a ...interface{}) string
	Green(format string, a ...interface{}) string
	Yellow(format string, a ...interface{}) string
	Blue(format string, a ...interface{}) string
	Magenta(format string, a ...interface{}) string
	Cyan(format string, a ...interface{}) string
	White(format string, a ...interface{}) string
	Black(format string, a ...interface{}) string
}

// ColorFunc provides a proxy interface for color.Color
type ColorFunc interface {
	Sprint(a ...interface{}) string
	Sprintf(format string, a ...interface{}) string
	Sprintln(a ...interface{}) string
	Print(a ...interface{}) (n int, err error)
	Printf(format string, a ...interface{}) (n int, err error)
	Println(a ...interface{}) (n int, err error)
}

// ColorImpl implements the Color interface using the color package
type ColorImpl struct{}

// ColorFuncImpl implements the ColorFunc interface wrapping color.Color
type ColorFuncImpl struct {
	colorFunc *color.Color
}

// NewColor creates a new Color implementation
func NewColor() Color {
	return &ColorImpl{}
}

func (c *ColorImpl) New(value ...color.Attribute) ColorFunc {
	return &ColorFuncImpl{colorFunc: color.New(value...)}
}

func (c *ColorImpl) Red(format string, a ...interface{}) string {
	return color.RedString(format, a...)
}

func (c *ColorImpl) Green(format string, a ...interface{}) string {
	return color.GreenString(format, a...)
}

func (c *ColorImpl) Yellow(format string, a ...interface{}) string {
	return color.YellowString(format, a...)
}

func (c *ColorImpl) Blue(format string, a ...interface{}) string {
	return color.BlueString(format, a...)
}

func (c *ColorImpl) Magenta(format string, a ...interface{}) string {
	return color.MagentaString(format, a...)
}

func (c *ColorImpl) Cyan(format string, a ...interface{}) string {
	return color.CyanString(format, a...)
}

func (c *ColorImpl) White(format string, a ...interface{}) string {
	return color.WhiteString(format, a...)
}

func (c *ColorImpl) Black(format string, a ...interface{}) string {
	return color.BlackString(format, a...)
}

func (cf *ColorFuncImpl) Sprint(a ...interface{}) string {
	return cf.colorFunc.Sprint(a...)
}

func (cf *ColorFuncImpl) Sprintf(format string, a ...interface{}) string {
	return cf.colorFunc.Sprintf(format, a...)
}

func (cf *ColorFuncImpl) Sprintln(a ...interface{}) string {
	return cf.colorFunc.Sprintln(a...)
}

func (cf *ColorFuncImpl) Print(a ...interface{}) (n int, err error) {
	return cf.colorFunc.Print(a...)
}

func (cf *ColorFuncImpl) Printf(format string, a ...interface{}) (n int, err error) {
	return cf.colorFunc.Printf(format, a...)
}

func (cf *ColorFuncImpl) Println(a ...interface{}) (n int, err error) {
	return cf.colorFunc.Println(a...)
}
