import configparser
import socket
import threading

def handle_client(client_socket):
    # Handle client requests here
    pass

def start_server(ip, port):
    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.bind((ip, port))
    server.listen(5)
    print(f"Server listening on {ip}:{port}")
    while True:
        client_socket, addr = server.accept()
        print(f"Accepted connection from {addr[0]}:{addr[1]}")
        client_handler = threading.Thread(target=handle_client, args=(client_socket,))
        client_handler.start()

if __name__ == "__main__":
    config = configparser.ConfigParser()
    config.read('config.ini')

    n = int(config['servers']['n'])
    threads = []

    for i in range(1, n+1):
        ip = config['servers'][f'server{i}_ip']
        port = int(config['servers'][f'server{i}_port'])
        thread = threading.Thread(target=start_server, args=(ip, port))
        thread.start()
        threads.append(thread)

    for thread in threads:
        thread.join()
