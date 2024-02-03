import socket
import threading

class Peer:
    def __init__(self, ip, port):
        self.ip = ip
        self.port = port

    def connect(self, peer_ip, peer_port):
        try:
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
                sock.connect((peer_ip, peer_port))
                print(f"Connected to peer {peer_ip}:{peer_port}")
        except Exception as e:
            print(f"Error connecting to peer {peer_ip}:{peer_port}: {e}")

    def start(self):
        # Example: Start a thread for peer operations
        threading.Thread(target=self.run).start()

    def run(self):
        # Example: Logic for peer operations
        pass
