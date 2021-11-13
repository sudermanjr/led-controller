FROM balenalib/raspberry-pi-golang:1.17-build AS builder
RUN [ "cross-build-start" ]
WORKDIR /tmp
ENV GO111MODULE=on
RUN apt-get update -y && apt-get install -y scons
RUN git clone https://github.com/jgarff/rpi_ws281x.git && \
  cd rpi_ws281x && \
  scons
RUN [ "cross-build-end" ]

FROM balenalib/raspberry-pi-golang:1.17
RUN [ "cross-build-start" ]
ENV GO111MODULE=on
COPY --from=builder /tmp/rpi_ws281x/*.a /usr/local/lib/
COPY --from=builder /tmp/rpi_ws281x/*.h /usr/local/include/
RUN [ "cross-build-end" ]
