import unittest
from src.seed import SeedNode

class TestSeedNode(unittest.TestCase):
    def test_seed_node_creation(self):
        seed_node = SeedNode()
        self.assertIsNotNone(seed_node.active_peers)

    def test_seed_node_add_remove_peer(self):
        seed_node = SeedNode()
        seed_node.add_peer("192.168.0.1", 5000)
        seed_node.add_peer("192.168.0.2", 5000)
        self.assertIn(("192.168.0.1", 5000), seed_node.active_peers)
        seed_node.remove_peer("192.168.0.1", 5000)
        self.assertNotIn(("192.168.0.1", 5000), seed_node.active_peers)

if __name__ == '__main__':
    unittest.main()
