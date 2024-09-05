import socket
import logging

import signal

from common.utils import Bet
from common.utils import store_bets
from .protocol import Protocol
from .message import ClientMessageParser

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_is_running = True

        # Initialize signals
        self.initialize_signals()

    def initialize_signals(self):
        signal.signal(signal.SIGTERM, self.stop)

    def stop(self, signum, frame):
        logging.info("action: server_stop | result: in_progress")
        self._server_is_running = False
        self._server_socket.close()
        logging.info("action: server_stop | result: success")


    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        try:
            while self._server_is_running:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
        except OSError: 
            logging.error("action: server_run | result: stopped")


    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """

        protocol = Protocol(client_sock)

        try:
            client_msg = protocol.receive()
            addr = client_sock.getpeername()

            parser = ClientMessageParser(client_msg)
            bets = parser.get_bets()

            store_bets(bets)

            logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)} | peso_kb: {len(client_msg) / 1024}')

            protocol.send(True, "message received")
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
