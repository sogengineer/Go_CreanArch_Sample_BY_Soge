# ベースイメージとして Golang の Alpine バージョンを使用
FROM golang:1.22-alpine

# 作業ディレクトリを設定
WORKDIR /app
COPY ./src .

# 必要なパッケージをインストール
RUN apk upgrade --update && \
    apk --no-cache add git bash && \
    apk upgrade git && \
    apk add --no-cache gcc musl-dev && \
    apk add --no-cache -X http://dl-cdn.alpinelinux.org/alpine/edge/community watchexec

# Go 依存関係をインストール
RUN go mod download && \
    go install honnef.co/go/tools/cmd/staticcheck@v0.3.0 && \
    go install github.com/swaggo/swag/cmd/swag@latest

ENV PATH="/root/go/bin:${PATH}"
# ホットリロード用のコマンドを設定
CMD ["watchexec", "-w", ".", "--force-poll", "100", "-r", "go", "run", "main.go"]
