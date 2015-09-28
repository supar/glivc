FROM debian:jessie
MAINTAINER supar "paulrez@gmail.com"

# Continuous integration image for the plate-go project
#

# Get rid of the debconf messages
ENV DEBIAN_FRONTEND noninteractive

# Update packages and install the ones that are needed
RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get install -y curl ca-certificates gcc libc6-dev \
                       bzr git-core mercurial locales \
                       build-essential debhelper devscripts \
                       fakeroot pkgconf \
                       make cmake \
                       --no-install-recommends \
                       && rm -rf /var/lib/apt/lists/*

RUN echo "Europe/Moscow" > /etc/timezone && dpkg-reconfigure --frontend noninteractive tzdata

# Set locale
RUN echo "LC_ALL=\"ru_RU.UTF-8\"" >> /etc/default/locale && \
    dpkg-reconfigure --frontend noninteractive locales

# Go install section
RUN curl -sSL https://storage.googleapis.com/golang/go1.4.3.linux-amd64.tar.gz \
    | tar -v -C /usr/local -xz

ENV PATH /usr/local/go/bin:$PATH
RUN mkdir -p /go/src
ENV GOPATH /go
ENV PATH /go/bin:$PATH

# Install libgit2.v0.22.3
RUN cd /opt && git clone https://github.com/libgit2/libgit2.git
RUN cd /opt/libgit2 && git fetch && git checkout tags/v0.22.3
RUN mkdir /opt/libgit2/build && cd /opt/libgit2/build && \
    cmake ../ -DLIB_INSTALL_DIR=/usr/lib/x86_64-linux-gnu && \
    cmake --build ./ --target install

# Install beegoo dependency
RUN go get github.com/astaxie/beego/logs && \
    go get gopkg.in/libgit2/git2go.v22 && \
    go get github.com/go-martini/martini

ADD run.sh /run.sh
ENTRYPOINT [ "/run.sh" ]

