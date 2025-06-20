import requests
import logging
from datetime import datetime
from typing import Optional, Dict, Any

logger = logging.getLogger(__name__)

class CoinGeckoClient:
    def __init__(self):
        self.base_url = "https://api.coingecko.com/api/v3"
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'crypto-trackers/1.0'
        })
    
    def fetch_price_data(self) -> Optional[list]:
        try:
            url = f"{self.base_url}/simple/price"
            params = {
                'ids': 'bitcoin,ethereum',
                'vs_currencies': 'usd',
                'include_market_cap': 'true',
                'include_24hr_vol': 'true',
                'include_24hr_change': 'true'
            }
            
            response = self.session.get(url, params=params, timeout=10)
            response.raise_for_status()
            
            data = response.json()
            
            price_events = []
            timestamp = datetime.utcnow().isoformat() + 'Z'
            
            if 'bitcoin' in data:
                btc_data = data['bitcoin']
                price_events.append({
                    "timestamp": timestamp,
                    "symbol": "BTC",
                    "price_usd": btc_data.get('usd', 0),
                    "volume_24h": btc_data.get('usd_24h_vol', 0),
                    "market_cap": btc_data.get('usd_market_cap', 0),
                    "price_change_24h": btc_data.get('usd_24h_change', 0),
                    "source": "coingecko"
                })
            
            if 'ethereum' in data:
                eth_data = data['ethereum']
                price_events.append({
                    "timestamp": timestamp,
                    "symbol": "ETH",
                    "price_usd": eth_data.get('usd', 0),
                    "volume_24h": eth_data.get('usd_24h_vol', 0),
                    "market_cap": eth_data.get('usd_market_cap', 0),
                    "price_change_24h": eth_data.get('usd_24h_change', 0),
                    "source": "coingecko"
                })
            
            logger.info(f"Successfully fetched price data for {len(price_events)} symbols")
            return price_events
            
        except requests.exceptions.RequestException as e:
            logger.error(f"API request failed: {e}")
            return None
        except Exception as e:
            logger.error(f"Unexpected error fetching price data: {e}")
            return None