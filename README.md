# LED Controller

This is a project to control NeoPixel lights with a PiZeroW and golang

## Usage

```
A cli for running neopixels

Usage:
  led-controller [flags]
  led-controller [command]

Available Commands:
  dashboard   Run a dashboard
  demo        Run a demo.
  help        Help about any command
  homekit     Run the lights as a homekit accessory.
  off         Turn off the lights.
  on          Turn on the lights.
  version     Prints the current version of the tool.

Flags:
      --add_dir_header                   If true, adds the file directory to the header
      --alsologtostderr                  log to standard error as well as files
  -f, --fade-duration int                The duration of fade-ins and fade-outs in ms. [FADE_DURTION] (default 100)
  -h, --help                             help for led-controller
  -l, --led-count int                    The number of LEDs in the array. [LED_COUNT] (default 12)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --log_file string                  If non-empty, use this log file
      --log_file_max_size uint           Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
      --logtostderr                      log to standard error instead of files (default true)
      --max-brightness int               The maximum brightness that will work within the 0-250 range. [MAX_BRIGHTNESS] (default 200)
      --min-brightness int               The minimum brightness that will work within the 0-250 range. [MIN_BRIGHTNESS] (default 25)
      --skip_headers                     If true, avoid header prefixes in the log messages
      --skip_log_headers                 If true, avoid headers when opening log files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          number for the log level verbosity
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

## Homekit

`led-controller homekit` will start this as a homekit device. Check the help for options, specifically the homekit pin. The homekit device will work with color and brightness controls once registered.

## Dashboard

Currently under heavy development. The dashboard will allow viewing and controlling the neopixel strip.

## Building your Pi Controller

In the [pi-builder](pi-builder) directory is the cloud-init for building my pi image using Hypriot. I have encrypted the cloud-init.yaml file using sops and pgp since it contains secrets. There's not much to this file, so it should be easy to re-create on your own if you like.

## Compiling

This is cross-compiled for the Raspberry Pi using the instructions in the rpi-ws281x repository. This utilizes a build container that is based on the [Balena Golang Image](https://registry.hub.docker.com/r/balenalib/raspberry-pi-golang). If you are using a different Pi, you will want to change the base image in the [Dockerfile](Dockerfile)

The build commands are in the [Makefile](Makefile). You can `make build` to build the Docker container and then use that to build the binary for the Pi. At the end it will show the output of `file led-controller` to verify the type of the binary.

The build will also create a local `.tmp` directory for storing build cache so that subsequent builds are much much faster.

Another word of caution: This container build will utilize your local GOPATH and GOCACHE so that it doesn't have to download every package every time.

## References

- golang library: https://github.com/rpi-ws281x/rpi-ws281x-go
- Pi Image Builder: https://github.com/hypriot/image-builder-rpi
