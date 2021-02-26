FROM harbor.xhhk/common/golangbuilder:1.15-alpine AS builder

ARG APP_NAME="kube-sweeper"
ARG BUILD_VERSION=
ARG BUILD_TIME=
ARG BUILD_NAME_ENDFIX=
ARG COMMIT_SHA1=
ARG COMMIT_SHORT=

RUN mkdir -p /app/${APP_NAME}
WORKDIR /app/${APP_NAME}

COPY . ./

ENV APP_NAME="${APP_NAME}" BUILD_VERSION="${BUILD_VERSION}" BUILD_TIME="${BUILD_TIME}" BUILD_NAME_ENDFIX="${BUILD_NAME_ENDFIX}" COMMIT_SHORT="${COMMIT_SHORT}"
RUN go build -ldflags="-w -s -X 'main.BuildVersion=$(echo ${BUILD_TIME})' -X 'main.BuiltTime=$(echo ${BUILD_TIME})' -X 'main.BuildName=$(echo ${APP_NAME}_${BUILD_NAME_ENDFIX})' -X 'main.CommitID=$(echo ${COMMIT_SHORT})'" -o ${APP_NAME} -i init.go main.go

RUN mkdir -p target && \
    mv ${APP_NAME} target/${APP_NAME} \
    && chmod +x target/${APP_NAME}

FROM harbor.xhhk/common/workerimage:alpine

LABEL author=starviews.com

ARG APP_NAME="kube-sweeper"

RUN mkdir -p /app/${APP_NAME}
WORKDIR /app/${APP_NAME}
COPY --from=builder /app/${APP_NAME}/target/* ./
RUN mkdir -p ./log
RUN echo -e "#!/bin/sh\n$@\n./${APP_NAME}\n" > /usr/local/bin/entrypoint.sh \
    && chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]
