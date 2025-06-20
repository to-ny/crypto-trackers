#!/bin/bash
set -e

if helm list | grep -q crypto-trackers; then
    helm uninstall crypto-trackers
fi

sleep 10

if kubectl get namespace crypto-trackers >/dev/null 2>&1; then
    kubectl get all -n crypto-trackers || true
    
    if kubectl get pvc -n crypto-trackers 2>/dev/null | grep -q .; then
        kubectl delete pvc --all -n crypto-trackers --timeout=60s || true
    fi
    
    kubectl delete namespace crypto-trackers --timeout=60s || true
fi

echo "Cleanup complete"