import json
import time
from datetime import datetime, timedelta
from typing import List, Dict
import math

class DataGenerator:
    def __init__(self):
        self.base_price = 67000.0
        self.base_volume = 25000000000
        
    def generate_golden_cross_scenario(self) -> List[Dict]:
        events = []
        start_time = datetime.now()
        
        for i in range(60):
            timestamp = start_time + timedelta(seconds=i*60)
            
            if i < 30:
                price = self.base_price - (30-i) * 50
            else:
                price = self.base_price + (i-30) * 100
                
            event = {
                "timestamp": timestamp.isoformat() + "Z",
                "symbol": "BTC",
                "price_usd": round(price, 2),
                "volume_24h": self.base_volume + (i * 100000000),
                "market_cap": round(price * 19700000, 0),
                "price_change_24h": round((price - self.base_price) / self.base_price * 100, 2),
                "source": "coingecko"
            }
            events.append(event)
            
        return events
        
    def generate_death_cross_scenario(self) -> List[Dict]:
        events = []
        start_time = datetime.now()
        
        for i in range(60):
            timestamp = start_time + timedelta(seconds=i*60)
            
            if i < 30:
                price = self.base_price + (30-i) * 50
            else:
                price = self.base_price - (i-30) * 100
                
            event = {
                "timestamp": timestamp.isoformat() + "Z",
                "symbol": "BTC",
                "price_usd": round(price, 2),
                "volume_24h": self.base_volume + (i * 100000000),
                "market_cap": round(price * 19700000, 0),
                "price_change_24h": round((price - self.base_price) / self.base_price * 100, 2),
                "source": "coingecko"
            }
            events.append(event)
            
        return events
        
    def generate_volume_spike_scenario(self) -> List[Dict]:
        events = []
        start_time = datetime.now()
        
        for i in range(60):
            timestamp = start_time + timedelta(seconds=i*60)
            price = self.base_price + math.sin(i * 0.1) * 100
            
            if i == 45:
                volume = self.base_volume * 1.5
            else:
                volume = self.base_volume + (i * 50000000)
                
            event = {
                "timestamp": timestamp.isoformat() + "Z",
                "symbol": "BTC",
                "price_usd": round(price, 2),
                "volume_24h": volume,
                "market_cap": round(price * 19700000, 0),
                "price_change_24h": round((price - self.base_price) / self.base_price * 100, 2),
                "source": "coingecko"
            }
            events.append(event)
            
        return events