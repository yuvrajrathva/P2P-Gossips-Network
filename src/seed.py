class SeedNode:
    def __init__(self):
        self.active_peers = set()

    def add_peer(self, peer_ip, peer_port):
        self.active_peers.add((peer_ip, peer_port))

    def remove_peer(self, peer_ip, peer_port):
        if (peer_ip, peer_port) in self.active_peers:
            self.active_peers.remove((peer_ip, peer_port))

    def broadcast_dead_peer(self, peer_ip, peer_port):
        # Example: Broadcast dead peer details to other peers
        pass
