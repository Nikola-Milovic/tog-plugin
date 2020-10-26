FROM heroiclabs/nakama-pluginbuilder:2.13.0 AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 1

WORKDIR /backend
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build --trimpath --buildmode=plugin -o ./backend.so

FROM heroiclabs/nakama:2.13.0

COPY --from=builder /backend/tog-plugin.so /nakama/data/modules
