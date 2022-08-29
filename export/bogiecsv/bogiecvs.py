import os
import asyncio
import logging
import nats
from nats.js.api import ConsumerConfig, DeliverPolicy, AckPolicy
import argparse
import bogie_pb2
import csvwrite

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
    try:
        nc = await nats.connect(args.server)
    except Exception as e:
        logging.error("Error connecting to NATS server: %s", e)
        return

    # Create JetStream context.
    js = nc.jetstream()

    config = ConsumerConfig(
        deliver_policy=DeliverPolicy.BY_START_TIME,
        opt_start_time="1990-01-01T00:00:00.000000Z",
        ack_policy=AckPolicy.EXPLICIT,
        # durable_name=CONSUMER,
    )

    sub = await js.subscribe(EXPORT_SUBJECT, stream=STREAM_NAME, config=config)

    while True:
        try:
            msg = await sub.next_msg(timeout=2.0)
        except asyncio.TimeoutError:
            logging.info("Timeout waiting for next message. Stop")
            break

        await msg.ack()

        m = decode_message(msg.data)
        print("got data on subject %s: " % (m))
        meta_csv.write(m)
        sensor_csv.write(m.sensor_samples)

    meta_csv.close()


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
