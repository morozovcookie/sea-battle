#golang:1.15.3-alpine3.12
FROM golang@sha256:6bdf3cc6cb0509e860e626a339c19b7e842efb8444960bb72d654e33607c34f5

RUN addgroup -S seabattle && \
    adduser -S seabattle -G seabattle

WORKDIR /src

ADD go.mod go.sum ./

RUN go mod download

ADD battle_field.go battle_field_cell.go battle_ship.go game.go game_state.go ./
ADD errors.go ./
ADD http ./http/
ADD cmd ./cmd/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w" \
    -o ./out/sea-battle \
    ./cmd/sea-battle/main.go


FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc/passwd /etc/passwd

USER seabattle

COPY --from=0 /src/out/sea-battle /usr/bin/

ENTRYPOINT ["sea-battle"]
