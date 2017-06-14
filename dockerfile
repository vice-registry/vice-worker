FROM alpine:latest
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN mkdir -p /opt/vice/
WORKDIR /opt/vice/
ADD vice-worker /opt/vice/
RUN chmod +x /opt/vice/vice-worker
ENV WORKERTYPE=import \
    COUCHBASE_LOCATION=localhost \
    COUCHBASE_USER=admin \
    COUCHBASE_PASS=admin \
    RABBITMQ_LOCATION=localhost \
    RABBITMQ_USER=admin \
    RABBITMQ_PASS=admin \
    STORAGE_BASEPATH=/tmp/
CMD /opt/vice/vice-worker \
    $( if [[ $WORKERTYPE == "import" ]]; then echo "--import"; fi) \
    $( if [[ $WORKERTYPE == "export" ]]; then echo "--export"; fi) \
    --couchbase-location $COUCHBASE_LOCATION \
    --couchbase-user $COUCHBASE_USER \
    --couchbase-pass $COUCHBASE_PASS \
    --rabbitmq-location $RABBITMQ_LOCATION \
    --rabbitmq-user $RABBITMQ_USER \
    --rabbitmq-pass $RABBITMQ_PASS \
    -- storage-basepath $STORAGE_BASEPATH