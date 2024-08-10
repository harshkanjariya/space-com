import argparse
from spacepackets.ecss.tc import PusTelecommand
from spacepackets.ecss.tm import PusTelemetry

def bytes_to_hex_string(data: bytes) -> str:
    """Convert bytes to a hex string without spaces."""
    return data.hex().upper()

def main(custom_message):
    # Convert custom message to bytes
    custom_message_bytes = custom_message.encode('utf-8')

    # Create a telecommand with custom payload
    ping_cmd = PusTelecommand(service=17, subservice=1, apid=0x01, app_data=custom_message_bytes)
    ping_cmd.data = custom_message_bytes  # Set custom message as data
    cmd_as_bytes = ping_cmd.pack()
    cmd_hex_string = bytes_to_hex_string(cmd_as_bytes)
    print("Telecommand: ", cmd_hex_string)  # Print only the hex string

    # Create a telemetry reply with custom payload
    ping_reply = PusTelemetry(service=17, subservice=2, source_data=custom_message_bytes, timestamp=bytes(), apid=0x01)
    tm_as_bytes = ping_reply.pack()
    tm_hex_string = bytes_to_hex_string(tm_as_bytes)
    print("Telemetry: ", tm_hex_string)  # Print only the hex string

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Generate PUS packets with custom message.")
    parser.add_argument("message", type=str, help="Custom message to include in the PUS packets.")
    args = parser.parse_args()

    main(args.message)
