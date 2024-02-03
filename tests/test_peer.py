import unittest
from src.peer import Peer

class TestPeer(unittest.TestCase):
    def test_peer_creation(self):
        peer = Peer("192.168.0.1", 5000)
        self.assertEqual(peer.ip, "192.168.0.1")
        self.assertEqual(peer.port, 5000)

    def test_peer_connection(self):
        peer = Peer("192.168.0.1", 5000)
        peer.connect("192.168.0.2", 5000)
        # Assert connection success

if __name__ == '__main__':
    unittest.main()
