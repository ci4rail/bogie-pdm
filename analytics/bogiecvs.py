import os
import asyncio
import logging
import argparse
import proto.bogie_pb2 as bogie_pb2
import csvwrite.csvwrite as csvwrite
from stream.stream import NatsStream

tool_description = """
Tool to export data from the bogie application to csv files.
"""

# Configurations for network
EXPORT_SUBJECT = "bogie"
# CONSUMER = "bogiecsv2"
STREAM_NAME = "test"


async def main(args):
    csvwrite.set_csv_dialect(decimalsep=",")
    meta_csv = csvwrite.MetaCsv(args.file + "-meta.csv")
    sensor_csv = csvwrite.SensorCsv(args.file + "-sensor.csv")

    logging.info("Connecting to NATS server: %s", args.server)
    ns = await NatsStream.from_start_time(
        args.server, STREAM_NAME, EXPORT_SUBJECT, "1990-01-01T00:00:00.000000Z"
    )

    while True:
        try:
            msg = await ns.next_msg()
        except asyncio.TimeoutError:
            logging.info("Timeout waiting for next message. Stop")
            break

        m = decode_message(msg.data)
        print("got data on subject %s: " % (m))
        meta_csv.write(m)
        sensor_csv.write(m.sensor_samples)

        await ns.ack(msg)

    meta_csv.close()
    sensor_csv.close()


def decode_message(msg):
    data = bogie_pb2.Bogie()
    data.ParseFromString(msg)
    return data


def command_line_args_parsing():
    parser = argparse.ArgumentParser(description=tool_description)
    parser.add_argument("server", help="NATS server address (e.g. localhost:4222)")
    parser.add_argument(
        "file",
        help="output files prefix, will be extended with sensor.csv and meta.csv",
    )
    parser.add_argument(
        "-r",
        "--runtime",
        help="runtime in seconds. (default: run forever)",
        type=int,
        default=None,
    )
    return parser.parse_args()


if __name__ == "__main__":
    logging.basicConfig(
        level=os.environ.get("LOGLEVEL", "INFO").upper(),
        format="%(asctime)s %(name)-12s %(levelname)-8s %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )
    args = command_line_args_parsing()
    asyncio.run(main(args))
