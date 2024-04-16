FROM debian:bookworm-slim
RUN apt update && apt install -y ca-certificates
COPY rocketreport-amd64 /app/rocketreport
ENTRYPOINT ["/app/rocketreport"]
