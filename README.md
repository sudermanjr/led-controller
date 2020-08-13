# LED Controller

This is a project to control NeoPixel lights with a PiZeroW and golang

## Usage

`led-controller help`

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
