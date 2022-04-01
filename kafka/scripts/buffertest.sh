#!/usr/bin/env bash
#
# Try running this script against you local kafka to see how
# kafkacat/kcat in cooperation with Kafka will buffer events
# before flushing them.
#
# Also try using the -X options to change this behavior.
#

KAFKA=${KAFKA:-localhost:9092}
TOPIC=buffertest
GROUP=$TOPIC.buffertest.sh
NUM_MESSAGES=${1:-3}  # use arg1 or use default=3
COMMON_ARGS="
-q -b $KAFKA
"

CONSUMER_ARGS="
-u -o stored -C -K: -G $GROUP $TOPIC
"
# -X heartbeat.interval.ms=100
# -X session.timeout.ms=100

PRODUCER_ARGS="
-P -K: -t $TOPIC
"
# tested params with no effect!
# -u
# -m 1
# -m 0.01
# -X batch.num.messages=1
# -X queue.buffering.max.messages=1
# -X queue.buffering.max.ms=10
# -X auto.commit.interval.ms=10
# -X linger.ms=0
# -X socket.timeout.ms=10
# -X max.in.flight.requests.per.connection=1
# -X request.timeout.ms=10
# -X message.timeout.ms=10
# -X offset.store.sync.interval.ms=1
# -X message.copy.max.bytes=10
# -X socket.send.buffer.bytes=10
# -X delivery.timeout.ms=10
# -X queue.buffering.max.kbytes=0.001  ERROR

KCAT="kcat $COMMON_ARGS"

kkcat() {
    log running $KCAT $*
    $KCAT "$@"
}

log(){ echo "[$(date):INFO] $*" 1>&2; }

# produce N messages and stop
producer() {
    log "producing $1 messages"
    for i in `seq $1`; do
        sleep 0.1
        echo "$i:message $(date)"
        if (($i%10 == 0 || $i == 1))
        then log "produced $i/$1 messages"
        fi
    done | kkcat -c$1 $PRODUCER_ARGS
    log "produced $1 messages"
}

# consume and N messages and stop
consumer() {
    log "starting consumer to read $1 messages"
    n=0
    kkcat -c$1 $CONSUMER_ARGS | {
        local n=0
        while read msg; do log "consumed: $msg"; ((n++)); done
        log "consumed $n messages"
    }
}

# produce/consume single message
produce() { echo "$@" | kkcat -T $PRODUCER_ARGS; }
consume() {             kkcat "$@" $CONSUMER_ARGS; }

cleanup() {
    # cleanup pending messages from failed tests
    # by reading messages for one second
    log "cleaning up topic $TOPIC"
    consume -e | wc -l | {
        read num
        log "cleaned $num messages from $TOPIC"
    }
}

check() {
    log "sending single message as health check"
    produce "1:test" | grep "1:test" >/dev/null &&
    log "[OK] produced message, waiting for consumer" &&
    consume -e | grep "1:test" >/dev/null
}

# synchronously consume and produce one message
# to create topic and check basic pipeline health
if cleanup && check
then log "[OK] can send and receive events"
else log "[ERR] failed to send and receive event"; exit 1
fi

# start actual test
consumer $NUM_MESSAGES&
sleep 0.5
producer $NUM_MESSAGES&
log "waiting for consumer and producer to stop after $NUM_MESSAGES"
wait
