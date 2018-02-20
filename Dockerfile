FROM golang:1.10 as builder

# Install Glide
ENV GLIDE_VERSION 0.13.1
ENV GLIDE_DOWNLOAD_URL https://github.com/Masterminds/glide/releases/download/v$GLIDE_VERSION/glide-v$GLIDE_VERSION-linux-amd64.tar.gz

RUN curl -fsSL "$GLIDE_DOWNLOAD_URL" -o glide.tar.gz \
	&& tar xvzf glide.tar.gz  linux-amd64/glide \
	&& mv linux-amd64/glide /usr/local/bin \
	&& rm -rf linux-amd64 \
	&& rm glide.tar.gz

# Install dependencies
WORKDIR /go/src/github.com/YotpoLtd/libra/
COPY glide.lock glide.lock
COPY glide.yaml glide.yaml
RUN glide install

# Copy sources into the container (see .dockerignore for excluded files)
COPY . /go/src/github.com/YotpoLtd/libra/

# Build the service app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/libra .

# The 'runtime' container only contains ssl cert chain, and is otherwise an empty base image.
# If you find this too bare, you can instead switch to using a normal base image.
FROM alpine

RUN apk --no-cache add ca-certificates

# Add binary
COPY --from=builder /go/bin/libra /bin/libra

# Expose service app ports
EXPOSE 8646

# Start the service app. Note we have to use the array style because this container does not include /bin/sh
ENTRYPOINT ["/bin/libra"]