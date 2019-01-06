FROM alpine:latest

ARG CONTAINER_BINARY

#ENV MQTT_USE_TLS "true"
#ENV MQTT_TLS_CA_PATH "/etc/mqtt-kube-operator/certs/DST_Root_CA_X3.pem"
ENV MQTT_USERNAME "xdk"
ENV MQTT_PASSWORD "xdk"
ENV MQTT_HOST "51.15.72.31"
ENV MQTT_PORT "1883"

RUN mkdir -p  /opt/mqtt-kube-operator
RUN chmod +x /opt/mqtt-kube-operator
COPY ./$CONTAINER_BINARY /opt/mqtt-kube-operator
#COPY ./certs/DST_Root_CA_X3.pem /etc/mqtt-kube-operator/certs/DST_Root_CA_X3.pem
WORKDIR /opt
ENTRYPOINT ["/opt/mqtt-kube-operator"]
