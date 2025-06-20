#!/bin/bash
set -e

if ! helm list | grep -q crypto-trackers; then
    echo "Error: crypto-trackers deployment not found"
    exit 1
fi

echo "Deployment status:"
kubectl get all -n crypto-trackers

KAFKA_READY=$(kubectl get pod kafka-0 -n crypto-trackers -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}')
ZK_READY=$(kubectl get pod zookeeper-0 -n crypto-trackers -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}')

echo "Pod status:"
echo "  Kafka: $KAFKA_READY"
echo "  ZooKeeper: $ZK_READY"

kubectl run verify-connectivity --image=busybox:1.35 --restart=Never -n crypto-trackers -- sh -c "nc -zv zookeeper-service 2181 && nc -zv kafka-service 9092" >/dev/null 2>&1
sleep 5
CONN_OUTPUT=$(kubectl logs verify-connectivity -n crypto-trackers 2>/dev/null || echo "failed")
kubectl delete pod verify-connectivity -n crypto-trackers >/dev/null 2>&1 || true

kubectl run kafka-verify --image=confluentinc/cp-kafka:7.4.0 --restart=Never -n crypto-trackers -- kafka-topics --bootstrap-server kafka-service:9092 --list >/dev/null 2>&1
sleep 10
TOPICS_OUTPUT=$(kubectl logs kafka-verify -n crypto-trackers 2>/dev/null || echo "failed")
kubectl delete pod kafka-verify -n crypto-trackers >/dev/null 2>&1 || true

echo "Connectivity: $(if echo "$CONN_OUTPUT" | grep -q "open"; then echo "OK"; else echo "FAILED"; fi)"
echo "Topics: $(if echo "$TOPICS_OUTPUT" | grep -q "crypto-prices" && echo "$TOPICS_OUTPUT" | grep -q "trading-signals"; then echo "crypto-prices, trading-signals"; else echo "MISSING"; fi)"

if [ "$KAFKA_READY" = "True" ] && [ "$ZK_READY" = "True" ] && echo "$CONN_OUTPUT" | grep -q "open" && echo "$TOPICS_OUTPUT" | grep -q "crypto-prices"; then
    echo "Verification successful"
else
    echo "Verification failed"
    exit 1
fi