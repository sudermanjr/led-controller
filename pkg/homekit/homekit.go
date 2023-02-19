package homekit

import (
	"github.com/brutella/hc" // TODO: move to the brutella/hap library since this is deprecated
	"github.com/brutella/hc/accessory"
	"github.com/lucasb-eyer/go-colorful"
	"go.uber.org/zap"

	"github.com/sudermanjr/led-controller/pkg/neopixel"
	"github.com/sudermanjr/led-controller/pkg/utils"
)

// Start starts the homekit server
func Start(homekitPin string, led *neopixel.LEDArray, logger *zap.SugaredLogger) {
	// create an accessory
	info := accessory.Info{
		Name:         "LED",
		Manufacturer: "sudermanjr",
		Model:        "led-controller",
	}
	ac := accessory.NewLightbulb(info)

	// configure the ip transport
	config := hc.Config{Pin: homekitPin}
	t, err := hc.NewIPTransport(config, ac.Accessory)
	if err != nil {
		logger.Fatalw("failed to start homekit", "error", err)
	}

	ac.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on {
			logger.Debugw("switch on")
			err = led.Fade(led.MaxBrightness)
			if err != nil {
				logger.Errorw("error changing brightness", "error", err)
			}
		} else {
			logger.Debugw("switch off")
			err = led.Fade(led.MinBrightness)
			if err != nil {
				logger.Errorw("error changing brightness", "error", err)
			}
		}
	})

	ac.Lightbulb.Hue.OnValueRemoteUpdate(func(value float64) {
		logger.Debugw("homekit hue set", "value", value)
		led.Color = modifyHue(led.Color, value, logger)
		err = led.Display(0)
		if err != nil {
			logger.Errorw("error changing color", "error", err)
		}
	})

	ac.Lightbulb.Saturation.OnValueRemoteUpdate(func(value float64) {
		logger.Debugw("homekit saturation set", "value", value)
		led.Color = modifySaturation(led.Color, value, logger)
		err = led.Display(0)
		if err != nil {
			logger.Errorw("error changing saturation", "error", err)
		}
	})

	ac.Lightbulb.Brightness.OnValueRemoteUpdate(func(value int) {
		logger.Debugw("homekit brightness set", "value", value)
		err = led.Fade(utils.ScaleBrightness(value, led.MinBrightness, led.MaxBrightness))
		if err != nil {
			logger.Errorw("could not set brightness", "error", err)
		}
	})

	hc.OnTermination(func() {
		logger.Infow("terminated - turning off lights")
		err = led.Fade(led.MinBrightness)
		if err != nil {
			logger.Errorw("error turning off lights", "error", err)
		}
		<-t.Stop()
	})

	logger.Infow("starting homekit server",
		"max brightness", led.MaxBrightness,
		"min brigntness", led.MinBrightness,
		"fade duration", led.FadeDuration,
	)

	t.Start()
}

// modifySaturation changes the saturation and returns a new color
func modifySaturation(oldColor colorful.Color, saturation float64, logger *zap.SugaredLogger) colorful.Color {
	h, s, v := oldColor.Hsv()
	logger.Debugw("old color", "hue", h, "saturation", s, "value", v)
	s = saturation * .1 // hc sends this 1-100, but colorful uses 0-1
	newColor := colorful.Hsv(h, s, v)
	logger.Debugw("new color", "hue", h, "saturation", s, "value", v)
	return newColor
}

// modifyHue changes the hue and returns a new color
func modifyHue(oldColor colorful.Color, hue float64, logger *zap.SugaredLogger) colorful.Color {
	h, s, v := oldColor.Hsv()
	logger.Debugw("new color", "hue", h, "saturation", s, "value", v)
	h = hue
	newColor := colorful.Hsv(h, s, v)
	logger.Debugw("new color", "hue", h, "saturation", s, "value", v)
	return newColor
}
