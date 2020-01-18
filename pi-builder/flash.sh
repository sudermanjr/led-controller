#!/bin/bash

flash \
    --bootconf no-uart-config.txt \
    --userdata cloud-init.yaml \
    https://github.com/hypriot/image-builder-rpi/releases/download/v1.12.0/hypriotos-rpi-v1.12.0.img.zip
