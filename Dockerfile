# Start from the latest Go 1.22 image
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o bin/big-john ./cmd
# need to run migrations
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz
# The 'tar xvz' command:
# x: extract files from the archive
# v: verbosely list files processed
# z: filter the archive through gzip (decompress)

# ***RUN STAGE***
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/bin/big-john ./big-john
COPY --from=build /app/home.html ./home.html
COPY --from=build /app/migrate ./migrate
COPY internal/db/postgresql/migration ./db/migration
COPY app.env .
COPY start.sh .
RUN chmod +x /app/start.sh  
# RUN ls -R

ENV LOG_LEVEL=1
ENV APP_ENV=production
ENV PORT=5001

EXPOSE 5001

# Run the binary
ENTRYPOINT [ "/app/start.sh" ]
CMD ["/app/big-john"]
