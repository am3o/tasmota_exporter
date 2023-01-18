FROM scratch

COPY tasmota_exporter /tasmota_exporter

ENTRYPOINT ["/tasmota_exporter"]