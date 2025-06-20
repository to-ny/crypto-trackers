import os
import sys
import unittest
from unittest.mock import Mock, patch

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "src"))


class TestKafkaEventProducer(unittest.TestCase):
    def setUp(self):
        with patch("kafka_client.producer.KafkaProducer"):
            from kafka_client.producer import KafkaEventProducer

            self.producer = KafkaEventProducer()

    def test_init_default_values(self):
        with patch("kafka_client.producer.KafkaProducer"):
            from kafka_client.producer import KafkaEventProducer

            producer = KafkaEventProducer()
            self.assertEqual(producer.topic, "crypto-prices")
            self.assertEqual(producer.bootstrap_servers, "kafka:9092")

    def test_init_custom_values(self):
        with patch("kafka_client.producer.KafkaProducer"):
            from kafka_client.producer import KafkaEventProducer

            producer = KafkaEventProducer(
                bootstrap_servers="localhost:9092", topic="test-topic"
            )
            self.assertEqual(producer.topic, "test-topic")
            self.assertEqual(producer.bootstrap_servers, "localhost:9092")

    def test_publish_price_events_success(self):
        mock_producer = Mock()
        mock_future = Mock()
        mock_producer.send.return_value = mock_future

        self.producer.producer = mock_producer

        price_events = [
            {
                "timestamp": "2024-06-16T14:30:00Z",
                "symbol": "BTC",
                "price_usd": 67450.23,
                "volume_24h": 28450000000,
                "market_cap": 1330000000000,
                "price_change_24h": 2.35,
                "source": "coingecko",
            },
            {
                "timestamp": "2024-06-16T14:30:00Z",
                "symbol": "ETH",
                "price_usd": 3450.89,
                "volume_24h": 15200000000,
                "market_cap": 414000000000,
                "price_change_24h": -1.25,
                "source": "coingecko",
            },
        ]

        result = self.producer.publish_price_events(price_events)

        self.assertTrue(result)
        self.assertEqual(mock_producer.send.call_count, 2)
        mock_producer.flush.assert_called_once()

        calls = mock_producer.send.call_args_list
        self.assertEqual(calls[0][1]["key"], "BTC")
        self.assertEqual(calls[1][1]["key"], "ETH")
        self.assertEqual(calls[0][1]["value"], price_events[0])
        self.assertEqual(calls[1][1]["value"], price_events[1])

    def test_publish_price_events_no_producer(self):
        self.producer.producer = None

        with patch.object(self.producer, "_connect") as mock_connect:
            mock_connect.return_value = None

            result = self.producer.publish_price_events([])

            self.assertFalse(result)
            mock_connect.assert_called_once()

    def test_close(self):
        mock_producer = Mock()
        self.producer.producer = mock_producer

        self.producer.close()

        mock_producer.close.assert_called_once()

    def test_close_no_producer(self):
        self.producer.producer = None

        self.producer.close()


if __name__ == "__main__":
    unittest.main()
