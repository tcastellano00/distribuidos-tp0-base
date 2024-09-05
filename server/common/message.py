from common.utils import Bet

class ClientMessageParser:
    def __init__(self, client_message):
        self.client_message = client_message

    def get_bets(self) -> list[Bet]:
        bets = []

        for message in self.client_message.split('\n'):
            message_info = message.split('|')
            bets.append(
                Bet(
                    message_info[0], 
                    message_info[1], 
                    message_info[2], 
                    message_info[3], 
                    message_info[4], 
                    message_info[5]
                )
            )

        return bets
    