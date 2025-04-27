FROM ghcr.io/vanilla-os/pico:main

WORKDIR /app
COPY . .

RUN apt update && apt install -y libudev-dev wget
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
RUN tar -C /usr -xzf go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/go/bin:${PATH}"

CMD ["go", "test", "-v", "./..."]