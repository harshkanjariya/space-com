#!/bin/bash

# aos signal
python3 aos-signal-generator.py "This is first message" --num-frames 10 --spacecraft-id 123 --virtual-channel-id 124

# telemetry signal
python3 telemetry-signal-generator.py -pt 5 "This is message"

# ccsds signal

