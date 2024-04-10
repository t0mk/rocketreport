FROM ubuntu:latest
RUN apt update && apt install -y ca-certificates
WORKDIR /app
COPY rocketreport /app/rocketreport
ENTRYPOINT ["/app/rocketreport"]