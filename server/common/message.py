from common.utils import Bet

CLIENT_MESSAGE_TYPE_BET = "bet"
CLIENT_MESSAGE_TYPE_READY = "ready"
CLIENT_MESSAGE_TYPE_RESULTS = "results"

class ClientMessageParser:
    def __init__(self, client_message):
        self.client_message = client_message

    def get_type(self):
        return (self.client_message.split('\n')[0])

    def get_bets(self) -> list[Bet]:
        #Validacion
        if (self.get_type() != CLIENT_MESSAGE_TYPE_BET):
            return None

        bets = []

        for message in self.client_message.split('\n')[1:]:
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

    def get_client_id(self) -> str:
        #Validacion
        if (self.get_type() == CLIENT_MESSAGE_TYPE_BET):
            return None

        return self.client_message.split('\n')[1]
    
    