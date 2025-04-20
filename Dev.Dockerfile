FROM golang:alpine AS build-proxy

WORKDIR /build

COPY cmd/ cmd/
COPY pkg/ pkg/
COPY go.mod .
COPY go.sum .

RUN CGO_ENABLED=0 go build -o proxy cmd/proxy/main.go

# ----------------------------------------------

FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y openssh-server && \
    apt-get clean

RUN mkdir /var/run/sshd \
    && chmod 0755 /var/run/sshd

COPY --from=build-proxy /build/proxy /usr/sbin/http-subsystem-proxy

# Permit root login via SSH
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# Enable password authentication
RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config

RUN echo "Subsystem http /usr/sbin/http-subsystem-proxy /etc/sshttproxy.env" >> /etc/ssh/sshd_config

RUN echo 'SSHTTPROXY_TARGET="http://server"\nSSHTTPROXY_LOGFILE="/var/log/proxy.log"\nSSHTTPROXY_LOGLEVEL="debug"'> /etc/sshttproxy.env

RUN echo 'root:root' | chpasswd

EXPOSE 22

CMD ["/usr/sbin/sshd", "-D"]