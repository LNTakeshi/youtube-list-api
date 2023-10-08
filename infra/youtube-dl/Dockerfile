# bump: golang /GOLANG_VERSION=([\d.]+)/ docker:golang|^1
# bump: golang link "Release notes" https://golang.org/doc/devel/release.html
ARG GOLANG_VERSION=1.20.4
# bump: yt-dlp /YT_DLP=([\d.-]+)/ https://github.com/yt-dlp/yt-dlp.git|/^\d/|sort
# bump: yt-dlp link "Release notes" https://github.com/yt-dlp/yt-dlp/releases/tag/$LATEST
ARG YT_DLP=2023.06.22

FROM golang:$GOLANG_VERSION AS base
ARG YT_DLP

WORKDIR /src

RUN \
  curl -L https://github.com/LNTakeshi/yt-dlp/releases/download/2023.08.25/yt-dlp_linux -o /usr/local/bin/yt-dlp && \
  chmod a+x /usr/local/bin/yt-dlp

COPY . ./
RUN \
  go mod download && \
  go build registry/api/run/main.go


FROM gcr.io/distroless/base-debian11 AS release
WORKDIR /src

COPY --from=base /src/main /src/main
COPY --from=base /usr/local/bin/yt-dlp /usr/local/bin/yt-dlp

ENV PORT 8080

EXPOSE $PORT

CMD ["./main"]