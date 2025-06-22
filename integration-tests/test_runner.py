#!/usr/bin/env python3

import json
import logging
import os
import sys
import time
from typing import Dict, List, Optional

from kafka import KafkaConsumer, KafkaProducer

from data_generator import DataGenerator

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class IntegrationTestRunner:
    def __init__(self) -> None:
        self.kafka_host = os.getenv("KAFKA_HOST", "kafka-service:9092")
        self.namespace = os.getenv("NAMESPACE", "crypto-trackers")
        self.producer: Optional[KafkaProducer] = None
        self.consumer: Optional[KafkaConsumer] = None
        self.data_generator = DataGenerator()

    def setup_kafka(self) -> bool:
        try:
            self.producer = KafkaProducer(
                bootstrap_servers=[self.kafka_host],
                value_serializer=lambda x: json.dumps(x).encode("utf-8"),
                retries=5,
                retry_backoff_ms=1000,
            )

            self.consumer = KafkaConsumer(
                "trading-signals",
                bootstrap_servers=[self.kafka_host],
                value_deserializer=lambda m: json.loads(m.decode("utf-8")),
                consumer_timeout_ms=60000,
                auto_offset_reset="latest",
            )
            logger.info("Kafka connections established")
            return True
        except Exception as e:
            logger.error(f"Kafka setup failed: {e}")
            return False

    def scale_service(self, service_name: str, replicas: int) -> bool:
        cmd = f"kubectl scale deployment {service_name} --replicas={replicas} -n {self.namespace}"
        exit_code = os.system(cmd)
        if exit_code == 0:
            logger.info(f"Scaled {service_name} to {replicas} replicas")
            time.sleep(10)
            return True
        logger.error(f"Failed to scale {service_name}")
        return False

    def restart_service(self, service_name: str) -> bool:
        cmd = f"kubectl rollout restart deployment {service_name} -n {self.namespace}"
        exit_code = os.system(cmd)
        if exit_code == 0:
            logger.info(f"Restarted {service_name}")
            time.sleep(15)
            return True
        logger.error(f"Failed to restart {service_name}")
        return False

    def reset_consumer_offsets(self) -> bool:
        reset_cmd = f"kubectl exec -n {self.namespace} kafka-0 -- kafka-consumer-groups.sh --bootstrap-server localhost:9092 --group ma-signal-detector --reset-offsets --to-latest --topic crypto-prices --execute"
        os.system(reset_cmd)

        reset_cmd = f"kubectl exec -n {self.namespace} kafka-0 -- kafka-consumer-groups.sh --bootstrap-server localhost:9092 --group volume-spike-detector --reset-offsets --to-latest --topic crypto-prices --execute"
        os.system(reset_cmd)

        logger.info("Reset consumer offsets")
        return True

    def setup_test_environment(self) -> bool:
        logger.info("Setting up test environment...")

        if not self.scale_service("crypto-trackers-data-ingestion", 0):
            return False

        if not self.restart_service("crypto-trackers-ma-signal-detector"):
            return False

        if not self.restart_service("crypto-trackers-volume-spike-detector"):
            return False

        if not self.reset_consumer_offsets():
            return False

        logger.info("Test environment setup complete")
        return True

    def cleanup_test_environment(self) -> bool:
        logger.info("Cleaning up test environment...")

        if not self.scale_service("crypto-trackers-data-ingestion", 1):
            return False

        logger.info("Test environment cleanup complete")
        return True

    def inject_test_data(self, events: List[Dict]) -> bool:
        try:
            for event in events:
                self.producer.send("crypto-prices", event)
                time.sleep(1)

            self.producer.flush()
            logger.info(f"Injected {len(events)} test events")
            return True
        except Exception as e:
            logger.error(f"Failed to inject test data: {e}")
            return False

    def collect_signals(self, timeout_seconds: int = 60) -> List[Dict]:
        signals = []
        start_time = time.time()

        try:
            for message in self.consumer:
                signals.append(message.value)
                logger.info(f"Received signal: {message.value}")

                if time.time() - start_time > timeout_seconds:
                    break

        except Exception as e:
            logger.info(f"Consumer timeout or error: {e}")

        return signals

    def verify_signals(
        self, signals: List[Dict], expected_type: str, expected_count: int
    ) -> bool:
        matching_signals = [s for s in signals if s.get("signal_type") == expected_type]

        if len(matching_signals) == expected_count:
            logger.info(
                f"✓ {expected_type} test passed: {expected_count} signals received"
            )
            return True
        else:
            logger.error(
                f"✗ {expected_type} test failed: expected {expected_count}, got {len(matching_signals)}"
            )
            return False

    def run_scenario_test(
        self,
        scenario_name: str,
        events: List[Dict],
        expected_signal_type: str,
        expected_count: int,
    ) -> bool:
        logger.info(f"Running {scenario_name} scenario...")

        if not self.inject_test_data(events):
            return False

        logger.info("Waiting for signal processing...")
        signals = self.collect_signals(60)

        return self.verify_signals(signals, expected_signal_type, expected_count)

    def run_all_tests(self) -> bool:
        if not self.setup_kafka():
            return False

        if not self.setup_test_environment():
            return False

        test_results = []

        try:
            death_cross_events = self.data_generator.generate_death_cross_scenario()
            result = self.run_scenario_test(
                "Death Cross", death_cross_events, "moving_average_crossover", 1
            )
            test_results.append(("Death Cross", result))

            time.sleep(5)

            volume_spike_events = self.data_generator.generate_volume_spike_scenario()
            result = self.run_scenario_test(
                "Volume Spike", volume_spike_events, "volume_spike", 1
            )
            test_results.append(("Volume Spike", result))

        finally:
            self.cleanup_test_environment()

        logger.info("\nTest Results:")
        all_passed = True
        for test_name, result in test_results:
            status = "PASS" if result else "FAIL"
            logger.info(f"{test_name}: {status}")
            if not result:
                all_passed = False

        return all_passed


if __name__ == "__main__":
    runner = IntegrationTestRunner()
    success = runner.run_all_tests()
    sys.exit(0 if success else 1)
