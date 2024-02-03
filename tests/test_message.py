import unittest
from src.message import Message

class TestMessage(unittest.TestCase):
    def test_message_creation(self):
        message = Message("type", "data")
        self.assertEqual(message.type, "type")
        self.assertEqual(message.data, "data")

    # Add more test cases for message serialization/deserialization if needed

if __name__ == '__main__':
    unittest.main()
