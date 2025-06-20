import os
import time
import logging
from flask import Flask, jsonify
from threading import Thread
from api.coingecko import CoinGeckoClient
from kafka_client.producer import KafkaEventProducer

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

app = Flask(__name__)

coingecko_client = CoinGeckoClient()
kafka_producer = None

@app.route('/health')
def health():
    return jsonify({"status": "healthy"})

@app.route('/ready')
def ready():
    global kafka_producer
    ready_status = kafka_producer is not None and kafka_producer.producer is not None
    return jsonify({"status": "ready" if ready_status else "not ready"})

def main_loop():
    global kafka_producer
    
    kafka_servers = os.environ.get('KAFKA_BOOTSTRAP_SERVERS', 'kafka:9092')
    kafka_producer = KafkaEventProducer(bootstrap_servers=kafka_servers)
    
    polling_interval = int(os.environ.get('POLLING_INTERVAL', 60))
    
    while True:
        try:
            logger.info("Fetching price data from CoinGecko")
            price_events = coingecko_client.fetch_price_data()
            
            if price_events:
                success = kafka_producer.publish_price_events(price_events)
                if success:
                    logger.info(f"Successfully processed {len(price_events)} price events")
                else:
                    logger.error("Failed to publish price events to Kafka")
            else:
                logger.warning("No price data retrieved from CoinGecko")
                
        except Exception as e:
            logger.error(f"Error in main loop: {e}")
        
        time.sleep(polling_interval)

if __name__ == '__main__':
    thread = Thread(target=main_loop, daemon=True)
    thread.start()
    
    port = int(os.environ.get('PORT', 8080))
    app.run(host='0.0.0.0', port=port)