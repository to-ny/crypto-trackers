import pytest
import json
from unittest.mock import Mock, patch, MagicMock
from src.main import app

class TestMainService:
    def setup_method(self):
        self.client = app.test_client()
        app.config['TESTING'] = True
    
    def test_health_endpoint(self):
        response = self.client.get('/health')
        
        assert response.status_code == 200
        data = json.loads(response.data)
        assert data['status'] == 'healthy'
    
    def test_ready_endpoint_with_kafka_producer(self):
        with patch('src.main.kafka_producer') as mock_producer:
            mock_producer.producer = Mock()
            
            response = self.client.get('/ready')
            
            assert response.status_code == 200
            data = json.loads(response.data)
            assert data['status'] == 'ready'
    
    def test_ready_endpoint_without_kafka_producer(self):
        with patch('src.main.kafka_producer', None):
            response = self.client.get('/ready')
            
            assert response.status_code == 200
            data = json.loads(response.data)
            assert data['status'] == 'not ready'
    
    def test_ready_endpoint_kafka_producer_no_connection(self):
        with patch('src.main.kafka_producer') as mock_producer:
            mock_producer.producer = None
            
            response = self.client.get('/ready')
            
            assert response.status_code == 200
            data = json.loads(response.data)
            assert data['status'] == 'not ready'