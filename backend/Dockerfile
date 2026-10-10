FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY bin/server ./server

EXPOSE 8080

ENTRYPOINT ["/server"]