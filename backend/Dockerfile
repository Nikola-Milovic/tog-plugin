FROM heroiclabs/nakama-pluginbuilder:2.14.0 AS builder

env GO111MODULE=on 
ENV CGO_ENABLED 1

WORKDIR /backend
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build --trimpath --buildmode=plugin -o ./tog-plugin.so

FROM heroiclabs/nakama:2.14.0

COPY --from=builder /backend/tog-plugin.so /nakama/data/modules
