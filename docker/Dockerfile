FROM debian:bookworm-slim
RUN apt update && apt install -y ca-certificates git
RUN mkdir /root/.rocketreport && git clone https://github.com/rocket-pool/rewards-trees /root/.rocketreport/reward-trees
COPY rocketreport-amd64 /app/rocketreport
ENTRYPOINT ["/app/rocketreport"]
