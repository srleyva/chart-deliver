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
RUN make && cp chart /chart

# Inject Binary into container
FROM sleyva97/base-layer:0.0.1
RUN gcloud components install kubectl
COPY --from=helm /linux-amd64/helm /bin/helm
COPY --from=build /chart /bin/chart
ADD cue-execute /usr/local/bin
ENTRYPOINT ["cue-execute"]
