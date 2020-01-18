# LED Controller

This is a project to control NeoPixel lights with a PiZeroW and golang

## pi-builder

This is the cloud-init for building my pi image using Hypriot. I have encrypted the cloud-init.yaml file using sops and pgp since it contains secrets. There's not much to this file, so it should be easy to re-create on your own if you like.

## References

- golang library: https://github.com/rpi-ws281x/rpi-ws281x-go
- Pi Image Builder: https://github.com/hypriot/image-builder-rpi
