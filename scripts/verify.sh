#!/bin/bash
set -e

check_pod_ready() {
    kubectl get pod $1 -n crypto-trackers -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}' 2>/dev/null || echo "False"
}

check_deployment_ready() {
    kubectl get deployment $1 -n crypto-trackers -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0"
}

check_health() {
    kubectl exec -n crypto-trackers deployment/$1 -- curl -s http://localhost:8080/health 2>/dev/null | grep -o '"status":"[^"]*"' | cut -d'"' -f4 || echo "FAILED"
}

run_temp_pod() {
    local name=$1
    local image=$2
    local cmd=$3
    
    kubectl run $name --image=$image --restart=Never -n crypto-trackers -- sh -c "$cmd" >/dev/null 2>&1
    sleep 5
    local output=$(kubectl logs $name -n crypto-trackers 2>/dev/null || echo "failed")
    kubectl delete pod $name -n crypto-trackers >/dev/null 2>&1 || true
    echo "$output"
}

if ! helm list -n crypto-trackers | grep -q crypto-trackers; then
    echo "Error: deployment not found"
    exit 1
fi

KAFKA_READY=$(check_pod_ready kafka-0)
ZK_READY=$(check_pod_ready zookeeper-0)

DATA_READY=$(check_deployment_ready crypto-trackers-data-ingestion)
MA_READY=$(check_deployment_ready crypto-trackers-ma-signal-detector)
VOLUME_READY=$(check_deployment_ready crypto-trackers-volume-spike-detector)
ALERT_READY=$(check_deployment_ready crypto-trackers-alert-service)
PROMETHEUS_READY=$(check_deployment_ready prometheus)

CONN_OUTPUT=$(run_temp_pod verify-connectivity busybox:1.35 "nc -zv zookeeper-service 2181 && nc -zv kafka-service 9092")
TOPICS_OUTPUT=$(run_temp_pod kafka-verify confluentinc/cp-kafka:7.4.0 "kafka-topics --bootstrap-server kafka-service:9092 --list")

DATA_HEALTH=$(check_health crypto-trackers-data-ingestion)
MA_HEALTH=$(check_health crypto-trackers-ma-signal-detector)
VOLUME_HEALTH=$(check_health crypto-trackers-volume-spike-detector)
ALERT_HEALTH=$(check_health crypto-trackers-alert-service)

PROMETHEUS_HEALTH=$(run_temp_pod verify-prometheus busybox:1.35 "wget -qO- prometheus-service:9090/-/healthy || echo 'FAILED'")

echo "Infrastructure: Kafka=$KAFKA_READY ZooKeeper=$ZK_READY"
echo "Services: Data=$([[ $DATA_READY == "1" ]] && echo "READY" || echo "NOT READY") MA=$([[ $MA_READY == "1" ]] && echo "READY" || echo "NOT READY") Volume=$([[ $VOLUME_READY == "1" ]] && echo "READY" || echo "NOT READY") Alert=$([[ $ALERT_READY == "1" ]] && echo "READY" || echo "NOT READY")"
echo "Monitoring: Prometheus=$([[ $PROMETHEUS_READY == "1" ]] && echo "READY" || echo "NOT READY")"
echo "Connectivity: $([[ $CONN_OUTPUT =~ "open" ]] && echo "OK" || echo "FAILED")"
echo "Topics: $([[ $TOPICS_OUTPUT =~ "crypto-prices" && $TOPICS_OUTPUT =~ "trading-signals" ]] && echo "OK" || echo "MISSING")"
echo "Health: Data=$DATA_HEALTH MA=$MA_HEALTH Volume=$VOLUME_HEALTH Alert=$ALERT_HEALTH"
echo "Monitoring Health: Prometheus=$([[ $PROMETHEUS_HEALTH =~ "ok" ]] && echo "OK" || echo "FAILED")"

if [[ $KAFKA_READY == "True" && $ZK_READY == "True" && $CONN_OUTPUT =~ "open" && $TOPICS_OUTPUT =~ "crypto-prices" && $DATA_READY == "1" && $MA_READY == "1" && $VOLUME_READY == "1" && $ALERT_READY == "1" && $PROMETHEUS_READY == "1" && $DATA_HEALTH == "healthy" && $MA_HEALTH == "healthy" && $VOLUME_HEALTH == "healthy" && $ALERT_HEALTH == "healthy" ]]; then
    echo "System verification SUCCESSFUL"
else
    echo "System verification FAILED"
    exit 1
fi