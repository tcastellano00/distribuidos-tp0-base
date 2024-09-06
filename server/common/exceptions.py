class ClientClosedException(Exception):
    def __init__(self, message="Closed connection by the client"):
        self.message = message
        super.__init__(message)