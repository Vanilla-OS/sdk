FROM ghcr.io/vanilla-os/pico:main

WORKDIR /app
COPY . .

RUN apt update && apt install -y libudev-dev wget
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
RUN tar -C /usr -xzf go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/go/bin:${PATH}"
# Something got wrong in latest Pico image, os-release looks incomplete, here
# we are fixing it to make it work with the rest of the system
RUN rm -f /etc/os-release && \
    echo 'PRETTY_NAME="Vanilla OS 2.0"' >> /etc/os-release && \
    echo 'NAME="Vanilla OS"' >> /etc/os-release && \
    echo 'VERSION_ID="2.0"' >> /etc/os-release && \
    echo 'VERSION="2.0 Orchid"' >> /etc/os-release && \
    echo 'VERSION_CODENAME="orchid"' >> /etc/os-release && \
    echo 'ID=vanilla' >> /etc/os-release && \
    echo 'ID_LIKE=debian' >> /etc/os-release && \
    echo 'HOME_URL="https://vanillaos.org"' >> /etc/os-release && \
    echo 'SUPPORT_URL="https://vanillaos.org/help"' >> /etc/os-release && \
    echo 'BUG_REPORT_URL="https://github.com/vanilla-os"' >> /etc/os-release && \
    echo 'PRIVACY_POLICY_URL="https://vanillaos.org/os-privacy-policy"' >> /etc/os-release

CMD ["go", "test", "-v", "./..."]