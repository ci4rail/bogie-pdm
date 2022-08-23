import os
import asyncio
import logging
import signal
from nats.aio.client import Client as Nats
import argparse

tool_description = """
Tool to export data from the bogie application to csv files.
"""

# Configurations for network
EXPORT_SUBJECT = "bogie"
CONSUMER = "bogiecsv"
STREAM_NAME = "test"


async def error_cb(e):
    logging.error(e)


async def main(args):
    nc = Nats()
    loop = asyncio.get_event_loop()

    # Connect to global nats
    options = {
        "servers": [args.server],
        "ping_interval": 1,
        "max_outstanding_pings": 5,
        # "user_credentials": credsfile_path,
        "error_cb": error_cb,
        "max_reconnect_attempts": 10,
    }

    logging.info("Connecting to NATS server: %s", args.server)
    try:
        await nc.connect(**options)
    except Exception as e:
        logging.error("Error connecting to NATS server: %s", e)
        return

    # Create JetStream context.
    js = nc.jetstream()

    async def cb(msg):
        global up_to_date
        await msg.ack()
        logging.info("got data on subject %s: %s", msg.subject, msg.data)

    # Create single push based subscriber that is durable across restarts.
    await js.subscribe(EXPORT_SUBJECT, durable=CONSUMER, cb=cb, stream=STREAM_NAME)

    # The following shuts down gracefully when SIGINT or SIGTERM is received
    stop = {"stop": False}

    def signal_handler():
        stop["stop"] = True

    for sig in ("SIGINT", "SIGTERM"):
        loop.add_signal_handler(getattr(signal, sig), signal_handler)

    # Fetch and ack messagess from consumer.
    while not stop["stop"]:
        await asyncio.sleep(1)

    logging.info("Shutting down...")


def command_line_args_parsing():
    parser = argparse.ArgumentParser(description=tool_description)
    parser.add_argument("server", help="NATS server address (e.g. localhost:4222)")
    parser.add_argument("file", help="output files prefix, will be extended with sensor.csv and meta.csv")
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
