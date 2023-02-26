#!/bin/bash

sops -d cloud-init.yaml > cloud-init.decrypted.yaml
~/repos/github.com/hypriot/flash/flash \
  --userdata cloud-init.decrypted.yaml \
  --bootconf no-uart-config.txt \
  https://github.com/hypriot/image-builder-rpi/releases/download/v1.12.3/hypriotos-rpi-v1.12.3.img.zip
