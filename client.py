# FILE          : client.py
# PROJECT       : Assignment3 - Services and Logging
# AUTHOR        : Youngwon Seo(8818834), Jiu Kim(8819115)
# DATE          : 2024.02.24
# DESCRIPTION   : This is a client-side file for implementing a multithreaded
#                 client that sends different types of login information to the server.

import socket
import argparse
import threading

# FUNCTION      : send_info()
# DESCRIPTION   : It sends login information to the server.
# PARAMETERS    : sock : socket.socket - The connection socket to the server
#                 loginInfo : str - Login information to send to the server
# RETURNS       : none


def send_info(sock, loginInfo):
    try:
        sock.sendall((loginInfo + "\n").encode('utf-8'))
        print(f"Information sent: {loginInfo}")
    except socket.error as err:
        print(f"Failed to send the information: {err}")


# FUNCTION      : create_connection_thread()
# DESCRIPTION   : It creates a thread that attempts to connect to the server and makes various login attempts.
# PARAMETERS    : host : str - The host address of the server
#                 port: int - The port number of the server
# RETURNS       : none
def create_connection_thread(host, port):
    try:
        # Attempt to connect to the server
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
            sock.connect((host, port))
            print(f"Connected to server: {host}:{port}")
            # Send the correct login information
            correctLoginInfo = "jkim9115:9115"
            send_info(sock, correctLoginInfo)

            # Sends invalid ID information
            incorrectIdInfo = "jkim7777:9115"
            send_info(sock, incorrectIdInfo)

            # Sends invalid password information
            incorrectPwInfo = "jkim9115:7777"
            send_info(sock, incorrectPwInfo)

            # Test making excessive requests to the server (rate limit test).
            for i in range(15):
                abuseInfo = f"abuseTest{i}"
                send_info(sock, abuseInfo)

    except socket.error as err:
        print(f"Failed to connect to server: {err}")


if __name__ == "__main__":
    # Create an ArgumentParser object to parse the command line arguments
    parser = argparse.ArgumentParser(
        description="Connect to the server as a client")
    parser.add_argument('--host', type=str,
                        help="The host address of the server", required=True)
    parser.add_argument('--port', type=int,
                        help="The port number of the server", required=True)

    args = parser.parse_args()

    # Create a thread list and start threads for multithreading
    threads = []
    for i in range(5):
        t = threading.Thread(target=create_connection_thread,
                             args=(args.host, args.port))
        threads.append(t)
        t.start()

    # Wait for all threads to terminate
    for t in threads:
        t.join()
