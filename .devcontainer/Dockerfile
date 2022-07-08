# Go binaries are standalone, so use a multi-stage build to produce smaller images.

# Use base golang image from Docker Hub
FROM golang:1.18 as build




LABEL name="splunk" \
      maintainer="support@splunk.com" \
      vendor="splunk" \
      version="1.1.0" \
      release="1" \
      summary="Simplify the Deployment & Management of Splunk Products on Kubernetes" \
      description="The Splunk Operator for Kubernetes (SOK) makes it easy for Splunk Administrators to deploy and operate Enterprise deployments in a Kubernetes infrastructure. Packaged as a container, it uses the operator pattern to manage Splunk-specific custom resources, following best practices to manage all the underlying Kubernetes objects for you."


WORKDIR /

# Install dependencies in go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the application source code
COPY . ./

# Compile the application to /manager.
# Skaffold passes in debug-oriented compiler flags
ARG SKAFFOLD_GO_GCFLAGS
RUN echo "Go gcflags: ${SKAFFOLD_GO_GCFLAGS}"
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -mod=readonly -v -o /manager


# Now create separate deployment image
FROM gcr.io/distroless/base

ENV OPERATOR=/manager \
    USER_UID=1001 \
    USER_NAME=nonroot

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

WORKDIR /
COPY --from=build /manager .
ENTRYPOINT ["./manager"]