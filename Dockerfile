FROM golang:1.25 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api social/cmd/api/*.go

#The run stage
FROM scratch
WORKDIR /app
#Copy CA Certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]