import time
import socket
import threading

class GossipProtocol:
    def __init__(self, seed_nodes):
        self.seed_nodes = seed_nodes
        self.peers = set()

    def join_network(self):
        # Connect to seed nodes and get list of active peers
        for seed_node in self.seed_nodes:
            self.connect_to_seed(seed_node)

    def connect_to_seed(self, seed_node):
        try:
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
                sock.connect(seed_node)
                # Receive list of active peers from seed
                # Add received peers to self.peers
        except Exception as e:
            print(f"Error connecting to seed node {seed_node}: {e}")

    def gossip(self, message):
        # Broadcast message to all connected peers
        for peer in self.peers:
            # Send message to peer

    def check_liveness(self):
        # Check liveness of connected peers
        for peer in self.peers:
            # Send heartbeat message to peer and wait for response
            # If no response, mark peer as dead and handle dead peer

    def handle_dead_peer(self, peer_ip, peer_port):
        # Remove dead peer from self.peers
        # Notify seed nodes about dead peer
        for seed_node in self.seed_nodes:
            seed_node.broadcast_dead_peer(peer_ip, peer_port)

    def start(self):
        # Start Gossip protocol operations in a separate thread
        threading.Thread(target=self.run).start()

    def run(self):
        # Main loop for Gossip protocol operations
        while True:
            # Periodically check liveness of peers
            self.check_liveness()
            # Sleep for a fixed interval
            time.sleep(10)
