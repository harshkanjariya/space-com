import struct
import zlib
import argparse


class AOSFrame:
    def __init__(self, spacecraft_id, virtual_channel_id, data):
        self.spacecraft_id = spacecraft_id
        self.virtual_channel_id = virtual_channel_id
        self.data = data
        self.frame = self.create_frame()

    def create_frame(self):
        """Create an AOS frame with the given data."""
        # Ensure the spacecraft_id is in the range of 0-65535 (2 bytes)
        if not (0 <= self.spacecraft_id <= 0xFFFF):
            raise ValueError("Spacecraft ID must be between 0 and 65535")

        # Ensure the virtual_channel_id is in the range of 0-255 (1 byte)
        if not (0 <= self.virtual_channel_id <= 0xFF):
            raise ValueError("Virtual Channel ID must be between 0 and 255")

        # Header: spacecraft ID (2 bytes) and virtual channel ID (1 byte)
        header = struct.pack('>H', self.spacecraft_id) + struct.pack('B', self.virtual_channel_id)

        # Frame: header + data
        frame = header + self.data

        # Calculate CRC32 and append it (4 bytes)
        crc = zlib.crc32(frame) & 0xffffffff
        frame += struct.pack('>I', crc)

        return frame

    def get_frame(self):
        """Return the generated frame."""
        return self.frame


class AOSGenerator:
    def __init__(self, num_frames, spacecraft_id, virtual_channel_id, message):
        self.num_frames = num_frames
        self.spacecraft_id = spacecraft_id
        self.virtual_channel_id = virtual_channel_id
        self.message = message

    def split_message(self, max_data_length):
        """Split the message into chunks."""
        return [self.message[i:i + max_data_length] for i in range(0, len(self.message), max_data_length)]

    def generate_signals(self):
        signals = []
        message_chunks = self.split_message(100)

        for i in range(min(self.num_frames, len(message_chunks))):
            data = message_chunks[i]
            print("Encoding: " + str(data))
            data = data.ljust(100, b'\x00')
            frame = AOSFrame(self.spacecraft_id, self.virtual_channel_id, data)
            signals.append(frame.get_frame())
        return signals


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Generate AOS signals from an input message.")
    parser.add_argument("message", type=str, help="The input message to be converted into AOS signals.")
    parser.add_argument("--num-frames", type=int, default=10, help="Number of frames to generate (default: 10).")
    parser.add_argument("--spacecraft-id", type=int, required=True, help="Spacecraft ID for the AOS frames.")
    parser.add_argument("--virtual-channel-id", type=int, required=True, help="Virtual channel ID for the AOS frames.")

    args = parser.parse_args()

    # Validate input values
    if not (0 <= args.spacecraft_id <= 0xFFFF):
        raise ValueError("Spacecraft ID must be between 0 and 65535")
    if not (0 <= args.virtual_channel_id <= 0xFF):
        raise ValueError("Virtual Channel ID must be between 0 and 255")

    generator = AOSGenerator(args.num_frames, args.spacecraft_id, args.virtual_channel_id, args.message.encode())
    signals = generator.generate_signals()

    # Output the generated frames in hexadecimal format
    for i, signal in enumerate(signals):
        print(f"Frame {i + 1}: {signal.hex()}")
