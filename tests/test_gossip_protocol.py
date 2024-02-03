import unittest
from src.gossip_protocol import GossipProtocol

class TestGossipProtocol(unittest.TestCase):
    def test_gossip_protocol_creation(self):
        seed_nodes = [("192.168.0.1", 5000), ("192.168.0.2", 5000)]
        gossip_protocol = GossipProtocol(seed_nodes)
        self.assertEqual(gossip_protocol.seed_nodes, seed_nodes)

    # Add more test cases for Gossip protocol functionalities

if __name__ == '__main__':
    unittest.main()
