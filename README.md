# LED Controller

This is a project to control NeoPixel lights with a PiZeroW and golang

## pi-builder

This is the cloud-init for building my pi image using Hypriot. I have encrypted the cloud-init.yaml file using sops and pgp since it contains secrets. There's not much to this file, so it should be easy to re-create on your own if you like.

## Cross Compiling

This is cross-compiled for the Raspberry Pi using the instructions in the rpi-ws281x repository. This utilizes a build container that is based on the [Balena Golang Image](https://registry.hub.docker.com/r/balenalib/raspberry-pi-golang). If you are using a different Pi, you will want to change the base image in the [Dockerfile](Dockerfile)

The build commands are in the [Makefile](Makefile). You can `make build` to build the Docker container and then use that to build the binary for the Pi. At the end it will show the output of `file led-controller` to verify the type of the binary.

Another word of caution: This container build will utilized your local GOPATH so that it doesn't have to download every package every time.

## References

- golang library: https://github.com/rpi-ws281x/rpi-ws281x-go
- Pi Image Builder: https://github.com/hypriot/image-builder-rpi
