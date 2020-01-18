#!/bin/bash

set -e

sops -d cloud-init.yaml > init.yaml
flash \
    --bootconf no-uart-config.txt \
    --userdata init.yaml \
    https://github.com/hypriot/image-builder-rpi/releases/download/v1.12.0/hypriotos-rpi-v1.12.0.img.zip

rm init.yaml
