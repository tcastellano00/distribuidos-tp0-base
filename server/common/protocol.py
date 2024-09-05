
RECEIVE_SIZE_BYTES = 1024

class Protocol:
    def __init__(self, socket):
        self.socket = socket

    def receive(self):
        buffer = b""

        while True:
            data = self.socket.recv(RECEIVE_SIZE_BYTES)
            
            if not data:
                return None

            buffer += data

            if b'\n\n' in buffer:
                return buffer[:buffer.index(b'\n\n')].strip().decode('utf-8')

    def send(self, success, text):
        message = ("OK" if success else "Error") + "," + text + "\n"
        data = message.encode('utf-8')

        sent_bytes = 0
        bytes_to_send = len(data)

        while sent_bytes < bytes_to_send:
            n = self.socket.send(data[sent_bytes:])
            if (n == 0):
                return None
            sent_bytes += n
            
