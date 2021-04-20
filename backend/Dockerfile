FROM heroiclabs/nakama-pluginbuilder:3.2.1 AS builder

env GO111MODULE=on 
ENV CGO_ENABLED 1

WORKDIR /backend

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build --trimpath --buildmode=plugin -o ./tog-plugin.so

FROM heroiclabs/nakama:3.2.1

COPY --from=builder /backend/tog-plugin.so /nakama/data/modules

RUN mkdir /nakama/resources

COPY --from=builder /backend/resources/ /nakama/resources
