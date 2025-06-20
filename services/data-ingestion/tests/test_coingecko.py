import unittest
import requests
from datetime import datetime
from unittest.mock import patch, Mock
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))

class TestCoinGeckoClient(unittest.TestCase):
    def setUp(self):
        from api.coingecko import CoinGeckoClient
        self.client = CoinGeckoClient()
    
    def test_fetch_price_data_success(self):
        mock_response_data = {
            "bitcoin": {
                "usd": 67450.23,
                "usd_24h_vol": 28450000000,
                "usd_market_cap": 1330000000000,
                "usd_24h_change": 2.35
            },
            "ethereum": {
                "usd": 3450.89,
                "usd_24h_vol": 15200000000,
                "usd_market_cap": 414000000000,
                "usd_24h_change": -1.25
            }
        }
        
        with patch('requests.Session.get') as mock_get:
            mock_response = Mock()
            mock_response.status_code = 200
            mock_response.json.return_value = mock_response_data
            mock_response.raise_for_status.return_value = None
            mock_get.return_value = mock_response
            
            with patch('api.coingecko.datetime') as mock_datetime:
                mock_datetime.utcnow.return_value.isoformat.return_value = "2024-06-16T14:30:00"
                
                result = self.client.fetch_price_data()
        
        self.assertIsNotNone(result)
        self.assertEqual(len(result), 2)
        
        btc_event = next(event for event in result if event['symbol'] == 'BTC')
        eth_event = next(event for event in result if event['symbol'] == 'ETH')
        
        self.assertEqual(btc_event['timestamp'], "2024-06-16T14:30:00Z")
        self.assertEqual(btc_event['symbol'], "BTC")
        self.assertEqual(btc_event['price_usd'], 67450.23)
        self.assertEqual(btc_event['volume_24h'], 28450000000)
        self.assertEqual(btc_event['market_cap'], 1330000000000)
        self.assertEqual(btc_event['price_change_24h'], 2.35)
        self.assertEqual(btc_event['source'], "coingecko")
        
        self.assertEqual(eth_event['symbol'], "ETH")
        self.assertEqual(eth_event['price_usd'], 3450.89)
        self.assertEqual(eth_event['price_change_24h'], -1.25)
    
    def test_fetch_price_data_api_error(self):
        with patch('requests.Session.get') as mock_get:
            mock_response = Mock()
            mock_response.raise_for_status.side_effect = requests.exceptions.HTTPError()
            mock_get.return_value = mock_response
            
            result = self.client.fetch_price_data()
            
            self.assertIsNone(result)
    
    def test_fetch_price_data_timeout(self):
        with patch('requests.Session.get') as mock_get:
            mock_get.side_effect = requests.exceptions.Timeout()
            
            result = self.client.fetch_price_data()
            
            self.assertIsNone(result)
    
    def test_client_initialization(self):
        self.assertEqual(self.client.base_url, "https://api.coingecko.com/api/v3")
        self.assertIsNotNone(self.client.session)
        self.assertEqual(self.client.session.headers['User-Agent'], 'crypto-trackers/1.0')

if __name__ == '__main__':
    unittest.main()