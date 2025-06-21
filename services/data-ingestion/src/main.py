import logging
import os
import time
from threading import Thread

from api.coingecko import CoinGeckoClient
from flask import Flask, jsonify
from kafka_client.producer import KafkaEventProducer
from prometheus_client import Counter, Histogram, generate_latest

logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)

app = Flask(__name__)

coingecko_client = CoinGeckoClient()
kafka_producer = None

api_calls_total = Counter(
    "coingecko_api_calls_total", "Total number of CoinGecko API calls", ["status"]
)
price_events_published = Counter(
    "price_events_published_total", "Total number of price events published", ["symbol"]
)
api_response_time = Histogram(
    "coingecko_api_response_seconds", "Time spent waiting for CoinGecko API response"
)


@app.route("/health")
def health():
    return jsonify({"status": "healthy"})


@app.route("/ready")
def ready():
    global kafka_producer
    ready_status = kafka_producer is not None and kafka_producer.producer is not None
    return jsonify({"status": "ready" if ready_status else "not ready"})


@app.route("/metrics")
def metrics():
    return generate_latest()


def main_loop():
    global kafka_producer

    kafka_servers = os.environ.get("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092")
    kafka_producer = KafkaEventProducer(bootstrap_servers=kafka_servers)

    polling_interval = int(os.environ.get("POLLING_INTERVAL", 60))

    while True:
        try:
            logger.info("Fetching price data from CoinGecko")

            with api_response_time.time():
                price_events = coingecko_client.fetch_price_data()

            if price_events:
                api_calls_total.labels(status="success").inc()
                success = kafka_producer.publish_price_events(price_events)
                if success:
                    for event in price_events:
                        price_events_published.labels(symbol=event["symbol"]).inc()
                    logger.info(
                        f"Successfully processed {len(price_events)} price events"
                    )
                else:
                    logger.error("Failed to publish price events to Kafka")
            else:
                api_calls_total.labels(status="failure").inc()
                logger.warning("No price data retrieved from CoinGecko")

        except Exception as e:
            api_calls_total.labels(status="error").inc()
            logger.error(f"Error in main loop: {e}")

        time.sleep(polling_interval)


if __name__ == "__main__":
    thread = Thread(target=main_loop, daemon=True)
    thread.start()

    port = int(os.environ.get("PORT", 8080))
    app.run(host="0.0.0.0", port=port)
