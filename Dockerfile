# Download helm
FROM alpine as helm
RUN apk update && apk add curl ca-certificates
RUN curl -o helm.tar.gz https://storage.googleapis.com/kubernetes-helm/helm-v2.12.0-linux-amd64.tar.gz && \
	tar -xvf helm.tar.gz

# Build Binary
FROM golang:alpine as build
RUN apk update && apk add make git tree
RUN mkdir -p $GOPATH/src/github.com/srleyva/chart-deliver
WORKDIR $GOPATH/src/github.com/srleyva/chart-deliver
ADD ./ ./
RUN make && cp chart-generator /chart-generator

ARG chart_cmd=print

# Inject Binary into container
FROM scratch
COPY --from=helm /etc/ssl/certs /etc/ssl/certs
COPY --from=helm /linux-amd64/helm /bin/helm
COPY --from=build /chart-generator /bin/chart-generator
CMD ["chart-generator"]
