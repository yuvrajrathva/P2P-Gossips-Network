import configparser
import socket

def connect_to_server(ip, port):
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect((ip, port))
    print(f"Connected to server at {ip}:{port}")

    # Add client logic here
    # For example, send and receive data

    # client_socket.close()

if __name__ == "__main__":
    config = configparser.ConfigParser()
    config.read('config.ini')

    n = int(config['servers']['n'])
    
    print("Available servers:")
    for i in range(1, n+1):
        ip = config['servers'][f'server{i}_ip']
        port = int(config['servers'][f'server{i}_port'])
        print(f"Server {i}: {ip}:{port}")

    choice = int(input("Enter the server number you want to connect to: "))

    if choice < 1 or choice > n:
        print("Invalid choice.")
    else:
        ip = config['servers'][f'server{choice}_ip']
        port = int(config['servers'][f'server{choice}_port'])
        connect_to_server(ip, port)
