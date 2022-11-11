#docker build -t my-prometheus .
#docker run -d -p 9090:9090 my-prometheus
FROM prom/prometheus
ADD prometheus.yml /etc/prometheus/
ADD test_server.json /etc/prometheus/