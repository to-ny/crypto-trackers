import json
import logging
from typing import Any

from kafka import KafkaProducer
from kafka.errors import KafkaError

logger = logging.getLogger(__name__)


class KafkaEventProducer:
    def __init__(
        self, bootstrap_servers: str = "kafka:9092", topic: str = "crypto-prices"
    ) -> None:
        self.topic = topic
        self.producer: KafkaProducer | None = None
        self.bootstrap_servers = bootstrap_servers
        self._connect()

    def _connect(self) -> None:
        try:
            self.producer = KafkaProducer(
                bootstrap_servers=self.bootstrap_servers,
                value_serializer=lambda v: json.dumps(v).encode("utf-8"),
                key_serializer=lambda k: k.encode("utf-8") if k else None,
                retries=3,
                retry_backoff_ms=1000,
                request_timeout_ms=30000,
            )
            logger.info(f"Connected to Kafka at {self.bootstrap_servers}")
        except Exception as e:
            logger.error(f"Failed to connect to Kafka: {e}")
            self.producer = None

    def publish_price_events(self, price_events: list[dict[str, Any]]) -> bool:
        if not self.producer:
            logger.warning("Kafka producer not connected, attempting reconnection")
            self._connect()
            if not self.producer:
                return False

        try:
            for event in price_events:
                symbol = event.get("symbol")
                future = self.producer.send(self.topic, value=event, key=symbol)

                future.add_callback(self._on_send_success)
                future.add_errback(self._on_send_error)

            self.producer.flush()
            logger.info(f"Published {len(price_events)} price events to {self.topic}")
            return True

        except KafkaError as e:
            logger.error(f"Kafka error publishing events: {e}")
            return False
        except Exception as e:
            logger.error(f"Unexpected error publishing events: {e}")
            return False

    def _on_send_success(self, record_metadata: Any) -> None:
        logger.debug(
            f"Message sent to {record_metadata.topic} partition {record_metadata.partition} offset {record_metadata.offset}"
        )

    def _on_send_error(self, excp: Exception) -> None:
        logger.error(f"Failed to send message: {excp}")

    def close(self) -> None:
        if self.producer:
            self.producer.close()
            logger.info("Kafka producer closed")
