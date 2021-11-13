package screen

import (
	"embed"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"time"

	"github.com/nfnt/resize"
	"github.com/sudermanjr/led-controller/pkg/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"k8s.io/klog"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/devices/ssd1306/image1bit"
	"periph.io/x/periph/host"
)

//go:embed gifs/*
var gifs embed.FS

// Display is a screen that you can display info on
type Display struct {
	LCD *ssd1306.Dev
}

// NewDisplay returns a *Display
func NewDisplay() (*Display, error) {

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		klog.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	bus, err := i2creg.Open("")
	if err != nil {
		klog.Fatal(err)
	}

	opts := &ssd1306.Opts{
		W:             128,
		H:             64,
		Rotated:       false,
		Sequential:    false,
		SwapTopBottom: false,
	}

	dev, err := ssd1306.NewI2C(bus, opts)
	if err != nil {
		klog.Fatalf("failed to initialize ssd1306: %v", err)
	}

	obj := &Display{LCD: dev}
	return obj, nil
}

// Demo runs a screen demo
func (display *Display) Demo() error {
	err := display.displayGif("ballerine.gif")
	if err != nil {
		return err
	}

	err = display.clear()
	if err != nil {
		return err
	}

	err = display.InfoDisplay()
	if err != nil {
		return err
	}

	err = display.ScrollText(2, 16, "This is a test.")
	if err != nil {
		return err
	}
	err = display.ScrollText(2, 32, "This is a test.")
	if err != nil {
		return err
	}
	err = display.ScrollText(2, 48, "This is a test.")
	if err != nil {
		return err
	}
	return nil
}

func (display *Display) clear() error {
	err := display.displayGif("black.gif")
	if err != nil {
		return err
	}
	return nil
}

func (display *Display) displayGif(gifName string) error {
	g, err := openGif(gifName)
	if err != nil {
		return err
	}

	// Converts every frame to image.Gray and resize them:
	imgs := make([]*image.Gray, len(g.Image))
	for i := range g.Image {
		imgs[i] = convertAndResizeAndCenter(128, 64, g.Image[i])
	}

	// Display the frames in a loop:
	for i := 0; i < len(imgs)*2; i++ {
		index := i % len(imgs)
		c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
		img := imgs[index]
		err := display.LCD.Draw(img.Bounds(), img, image.Point{})
		if err != nil {
			return err
		}
		<-c
	}
	return nil
}

// Text displays text at a position
func (display *Display) Text(x int, y int, message string) error {
	img := display.stringImage(x, y, message)

	err := display.LCD.Draw(img.Rect, img, image.Point{x, y})
	if err != nil {
		return err
	}
	return nil
}

// ScrollText displays text and scrolls it
func (display *Display) ScrollText(x int, y int, message string) error {
	img := display.stringImage(x, y, message)
	err := display.LCD.Draw(img.Rect, img, image.Point{x, y})
	if err != nil {
		return err
	}
	scrollMin := nearestMultiple(y, 8, false)
	if scrollMin < display.LCD.Bounds().Min.X {
		scrollMin = display.LCD.Bounds().Min.Y
	}
	scrollMax := nearestMultiple(y+16, 8, true)
	if scrollMax > display.LCD.Bounds().Max.Y {
		scrollMax = display.LCD.Bounds().Max.Y
	}

	klog.V(7).Infof("scroll min: %d, max: %d", scrollMin, scrollMax)
	err = display.LCD.Scroll(ssd1306.Left, ssd1306.FrameRate5, scrollMin, scrollMax)
	if err != nil {
		return err
	}
	return nil
}

// stringImage returns an img that is a string of text
func (display *Display) stringImage(x int, y int, message string) *image1bit.VerticalLSB {
	f := basicfont.Face7x13
	bounds := image.Rect(x, y, x+(display.LCD.Bounds().Dx()-x), y+f.Height)

	klog.V(5).Infof("image bounds: %v", bounds)
	klog.V(5).Infof("screen bounds: %v", display.LCD.Bounds())
	img := image1bit.NewVerticalLSB(bounds)

	drawer := font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{image1bit.On},
		Face: f,
		Dot:  fixed.P(x, y+f.Height),
	}
	drawer.DrawString(message)
	return img
}

// InfoDisplay shows current information
func (display *Display) InfoDisplay() error {
	ip, err := utils.IPAddress()
	if err != nil {
		klog.Error(err)
		return err
	}
	err = display.Text(0, 0, fmt.Sprintf(" IP: %s", ip))
	if err != nil {
		return err
	}
	err = display.LCD.Scroll(ssd1306.Left, ssd1306.FrameRate5, 0, 16)
	if err != nil {
		return err
	}
	return nil
}

// convertAndResizeAndCenter takes an image, resizes and centers it on a
// image.Gray of size w*h.
func convertAndResizeAndCenter(w, h int, src image.Image) *image.Gray {
	src = resize.Thumbnail(uint(w), uint(h), src, resize.Bicubic)
	img := image.NewGray(image.Rect(0, 0, w, h))
	r := src.Bounds()
	r = r.Add(image.Point{(w - r.Max.X) / 2, (h - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	return img
}

func openGif(fileName string) (*gif.GIF, error) {
	f, err := gifs.Open(fileName)
	if err != nil {
		return nil, err
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	if err != nil {
		return nil, err
	}
	return g, nil
}

// nearestMultiple returns the clostest multiple of target
// that is above or below the input
func nearestMultiple(input int, target int, up bool) int {
	if input%target == 0 {
		return input
	}
	a := (input / target) * target
	b := a + target

	if up {
		if a < input {
			return a + target
		}
		return a
	}
	if !up {
		if b > input {
			return b - target
		}
		return b
	}
	return -1
}
