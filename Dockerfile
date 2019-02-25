# First Stage
FROM golang:1.11.5-alpine3.9 as builder

RUN mkdir /app
ADD . /app/
WORKDIR /app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


# Second Stage
FROM alpine
EXPOSE 80
EXPOSE 443
EXPOSE 18080
CMD ["/app/bin/m3d"]
ENTRYPOINT ["/app/bin/m3d", "--base", "/app/dhnt"]

# Copy from first stage
COPY --from=builder /app /app