# 程序编译镜像
FROM golang:1.22 as build

ENV GO111MODULE=on                               \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ./dist/nav-green-download   ./cmd/nav-green-download  \
    && GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ./dist/scripts           ./cmd/scripts 

# 构建镜像
FROM alpine:3.20 as prod

ENV SMOC_DIR=./nc_files/smoc        \
    MFWAM_DIR=./nc_files/mfwam      \
    SEA_ICE_DIR=./nc_files/sea_ice  \
    EC_DIR=./nc_files/ec_0p25       \
    LOG_FILE=""                     \
    LOG_LEVEL="info"         

# 在build阶段复制时区到
COPY --from=build /go/release/dist/nav-green-download   /nav-green-download
COPY --from=build /go/release/dist/scripts              /scripts

CMD [ "/nav-green-download"]