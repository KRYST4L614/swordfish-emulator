# Базовый минимальный образ
FROM golang:alpine3.18 as builder

WORKDIR /build

RUN apk update && apk add --no-cache git 

RUN git clone --branch main --single-branch https://gitlab.com/IgorNikiforov/swordfish-emulator-go.git

WORKDIR /build/swordfish-emulator-go

RUN go mod download && go build -o /bin/swordfish-emulator-go ./cmd/emulator/main.go

FROM alpine:3.18 as main 

RUN apk update && apk add bash

# Create a group and user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Tell docker that all future commands should run as the appuser user

WORKDIR /bin/

COPY --from=builder /bin/swordfish-emulator-go /bin/swordfish-emulator-go

COPY --from=builder /build/swordfish-emulator-go/datasets /bin/datasets

COPY --from=builder /build/swordfish-emulator-go/scripts /bin/scripts

COPY --from=builder /build/swordfish-emulator-go/database /bin/database

COPY --from=builder /build/swordfish-emulator-go/configs /bin/configs

RUN chmod +x /bin/swordfish-emulator-go

RUN chown -R appuser /bin/swordfish-emulator-go

RUN chown -R appuser /bin/datasets

RUN chown -R appuser /bin/scripts

RUN chown -R appuser /bin/database

RUN chown -R appuser /bin/configs

RUN chmod 777 /bin


# Устанавливаем точку входа
ENTRYPOINT ["bash", "/bin/scripts/nfs_setup.sh", "/bin/swordfish-emulator-go"]

USER appuser

ENTRYPOINT ["/bin/swordfish-emulator-go"]