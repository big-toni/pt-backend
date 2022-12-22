# simple working example
# FROM golang:1.16.3
# RUN mkdir /app
# ADD . /app
# WORKDIR /app
# RUN go build -o main .
# CMD ["/app/main"]

# FROM golang:1.16.3 as debug

FROM golang:1.19

# RUN apk update && apk upgrade && \
#     apk add --no-cache git \
#         dpkg \
#         gcc \
#         git \
#         musl-dev

ENV gopath /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# install google chrome
# ENV CHROME_VERSION "google-chrome-stable"
# RUN sed -i -- 's&deb http://deb.debian.org/debian jessie-updates main&#deb http://deb.debian.org/debian jessie-updates main&g' /etc/apt/sources.list \
#   && apt-get update && apt-get install wget -y
# ENV CHROME_VERSION "google-chrome-stable"
# RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
#   && echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list \
#   && apt-get update && apt-get -qqy install ${CHROME_VERSION:-google-chrome-stable}
# CMD /bin/bash


# Install Chromium
# Yes, including the Google API Keys sucks but even debian does the same: https://packages.debian.org/stretch/amd64/chromium/filelist
RUN apt-get update && apt-get install -y \
      chromium \
      chromium-l10n \
      fonts-liberation \
      fonts-roboto \
      hicolor-icon-theme \
      libcanberra-gtk-module \
      libexif-dev \
      libgl1-mesa-dri \
      libgl1-mesa-glx \
      libpangox-1.0-0 \
      libv4l-0 \
      fonts-symbola \
      --no-install-recommends \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /etc/chromium.d/ \
    && /bin/echo -e 'export GOOGLE_API_KEY="AIzaSyCkfPOPZXDKNn8hhgu3JrA62wIgC93d44k"\nexport GOOGLE_DEFAULT_CLIENT_ID="811574891467.apps.googleusercontent.com"\nexport GOOGLE_DEFAULT_CLIENT_SECRET="kdloedMFGdGla2P1zacGjAQh"' > /etc/chromium.d/googleapikeys

# Add chromium user
RUN groupadd -r chromium && useradd -r -g chromium -G audio,video chromium \
    && mkdir -p /home/chromium/Downloads && chown -R chromium:chromium /home/chromium

# Run as non privileged user
# USER chromium

# RUN go get github.com/sirupsen/logrus
# RUN go get github.com/buaazp/fasthttprouter
# RUN go get github.com/valyala/fasthttp
# RUN go get github.com/go-delve/delve/cmd/delve

# RUN go get github.com/go-delve/delve/cmd/dlv@latest


### Run the Delve debugger ###
# COPY ./dlv.sh /
# RUN chmod +x /dlv.sh
# ENTRYPOINT ["/dlv.sh"]

# Run delve
# CMD ["/dlv", "--listen=:2345", "--headless=true", "--api-version=2", "exec", "/app"]


### Start new image ###

# FROM golang:1.16.3
# COPY --from=debug /app /
# CMD ./app



# WORKDIR /app
# COPY . /app/

# RUN go mod download

# RUN go build -o app

# RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add -

# RUN sh -c 'echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list'
# RUN apt-get update -qqy --no-install-recommends && apt-get install -qqy --no-install-recommends google-chrome-stable


WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o pt-backend
EXPOSE 3000
CMD [ "./pt-backend" ]

