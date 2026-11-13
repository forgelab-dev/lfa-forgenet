# syntax=docker/dockerfile:1
#
# Image du serveur ForgeNet : le binaire, le dashboard compilé (ADR-0009) et la
# CLI docker avec Compose, que le serveur appelle pour les stacks (ADR-0003).
# Construite depuis la racine du dépôt :
#   docker build -f deploy/Dockerfile -t forgenet-server .
#
# Chaque image de base est épinglée par version et par empreinte.

FROM node:26.4.0-alpine@sha256:725aeba2364a9b16beae49e180d83bd597dbd0b15c47f1f28875c290bfd255b9 AS web
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY web/ ./
RUN npm run build

FROM golang:1.26.11-alpine@sha256:ae52b5bc6ed045a9ed26ab03dc19327a802a9cc9aa8b27ed57f5c6b60579b0bb AS server
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
ARG VERSION=0.1.0-alpha.dev
ARG GIT_COMMIT=unknown
# SQLite en Go pur (modernc) : binaire statique, sans cgo.
RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags "-s -w -X github.com/forgelab-dev/lfa-forgenet/internal/version.Version=${VERSION} -X github.com/forgelab-dev/lfa-forgenet/internal/version.GitCommit=${GIT_COMMIT}" \
    -o /out/forgenet-server ./cmd/forgenet-server

FROM docker:29.9.3-cli@sha256:5a7af0dcb0e46c803efe9b45371dd21cb500e5b4453250394ac6f06d49b6e6f1

# Un utilisateur sans droits : l'accès à Docker passe par le groupe du socket,
# ajouté au lancement (group_add), jamais par root.
RUN addgroup -S -g 10001 forgenet \
    && adduser -S -D -H -u 10001 -G forgenet -h /var/lib/forgenet forgenet \
    && mkdir -p /var/lib/forgenet \
    && chown forgenet:forgenet /var/lib/forgenet \
    && chmod 0700 /var/lib/forgenet

COPY --from=server /out/forgenet-server /usr/local/bin/forgenet-server
COPY --from=web /src/web/dist /usr/share/forgenet/web

# HOME et DOCKER_CONFIG pointent vers /tmp : le système de fichiers du
# conteneur est en lecture seule (docker-compose.yml).
ENV FORGENET_ENVIRONMENT=production \
    FORGENET_SERVER_HOST=0.0.0.0 \
    FORGENET_DATABASE_PATH=/var/lib/forgenet/forgenet.db \
    FORGENET_WEB_ROOT=/usr/share/forgenet/web \
    FORGENET_LOG_FORMAT=json \
    HOME=/tmp \
    DOCKER_CONFIG=/tmp/.docker

USER 10001:10001
WORKDIR /var/lib/forgenet
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/forgenet-server"]
