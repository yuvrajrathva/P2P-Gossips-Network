class Message:
    def __init__(self, type, data):
        self.type = type
        self.data = data

    def serialize(self):
        # Serialize message to send over network
        pass

    @classmethod
    def deserialize(cls, serialized_data):
        # Deserialize received data into Message object
        pass
