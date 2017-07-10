# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/antonve/kareebot

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN curl https://glide.sh/get | sh
RUN cd /go/src/github.com/antonve/kareebot && glide install
RUN go install github.com/antonve/kareebot

# Move dictionaries to where the binary expects them
RUN mv /go/src/github.com/antonve/kareebot/dicts /go/bin/dicts

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/kareebot
