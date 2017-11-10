FROM alpine:latest
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk update && apk add --upgrade ca-certificates
RUN mkdir -p /opt/vice/
WORKDIR /opt/vice/
ADD vice-worker /opt/vice/
RUN chmod +x /opt/vice/vice-worker
ENV WORKERTYPE=import \
    RETHINKDB_LOCATION=localhost \
    RETHINKDB_DATABASE=vice \
    RABBITMQ_LOCATION=localhost \
    RABBITMQ_USER=admin \
    RABBITMQ_PASS=admin \
    STORAGE_BASEPATH=/tmp/
CMD /opt/vice/vice-worker \
    $( if [[ $WORKERTYPE == "import" ]]; then echo "--import"; fi) \
    $( if [[ $WORKERTYPE == "export" ]]; then echo "--export"; fi) \
    --rethinkdb-location $RETHINKDB_LOCATION \
    --rethinkdb-database $RETHINKDB_DATABASE \
    --rabbitmq-location $RABBITMQ_LOCATION \
    --rabbitmq-user $RABBITMQ_USER \
    --rabbitmq-pass $RABBITMQ_PASS \
    -- storage-basepath $STORAGE_BASEPATH