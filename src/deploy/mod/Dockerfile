# Main Dockerfile of the project. It creates an image that serves the application. This image should
# be built from the dist directory.

# Scratch can be used as the base image because the backend is compiled to include all
# its dependencies.
FROM scratch
LABEL maintainer "Piotr Bryk <bryk@google.com>"

# Add all files from current working directory to the root of the image, i.e., copy dist directory
# layout to the root directory.
ADD . /

# The port that the application listens on.
# TODO(bryk): Parametrize this argument so that other build tools are aware of the exposed port.
EXPOSE 9090 8443
ENTRYPOINT ["/dashboard"]
