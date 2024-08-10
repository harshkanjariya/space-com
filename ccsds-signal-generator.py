import argparse
from gnuradio import gr, blocks
import ccsds

class CCSDSHexEncoder(gr.top_block):
    def __init__(self, input_message):
        gr.top_block.__init__(self)

        # Convert the input message to a byte array
        input_data = list(input_message.encode())

        # Define parameters for CCSDS
        frame_size = 223         # Example frame size (based on CCSDS standards)
        rs_interleave = 1        # Reed-Solomon interleaving (can be adjusted)
        scrambling = True        # Apply scrambling (True/False based on needs)

        # CCSDS Encoder block
        self.ccsds_encoder = ccsds.encoder(frame_size, rs_interleave, scrambling)

        # Source block with the input data
        self.source = blocks.vector_source_b(input_data, repeat=False)

        # Sink to capture the encoded output
        self.hex_sink = blocks.vector_sink_b()

        # Connect the blocks
        self.connect(self.source, self.ccsds_encoder, self.hex_sink)

    def get_hex_output(self):
        self.run()
        encoded_bytes = self.hex_sink.data()
        return ''.join(format(byte, '02x') for byte in encoded_bytes)

if __name__ == "__main__":
    # Parse command-line arguments
    parser = argparse.ArgumentParser(description="CCSDS Hex Signal Generator")
    parser.add_argument("input_message", type=str, help="The message to encode in CCSDS format")
    args = parser.parse_args()

    # Create CCSDSHexEncoder object with the input message
    encoder = CCSDSHexEncoder(args.input_message)

    # Get and print the hex output
    hex_output = encoder.get_hex_output()
    print(f"Input Message: {args.input_message}")
    print(f"CCSDS Hex Signal: {hex_output}")
