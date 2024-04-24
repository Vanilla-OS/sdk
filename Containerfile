FROM ghcr.io/vanilla-os/pico:main

WORKDIR /app
COPY . .

RUN apt update && apt install -y libudev-dev golang-go

CMD ["go", "test", "-v", "./..."]