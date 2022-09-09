import stream.stream as stream
import stream.metricspandas as metricpandas
import pandas as pd
from dateutil import tz
import datetime


async def nats_stream_fetch(
    server, nats_creds, stream_name, export_subject, start_time, end_time
):
    ns = await stream.NatsStream.from_start_time(
        server, nats_creds, stream_name, export_subject, start_time
    )

    end_time = end_time.replace(tzinfo=tz.tzlocal())
    # print("nats_stream_fetch until %s\n" % end_time)
    timeout = 2.0
    df = pd.DataFrame()
    while True:
        msg = await ns.next_msg(timeout=timeout)
        if msg is None:
            # print("msg is None")
            break
        await ns.ack(msg)
        df2 = metricpandas.metrics_nats_to_pandas(msg)

        ts = df2["nats_rx_time"].iloc[0].replace(tzinfo=tz.tzlocal())
        if ts > end_time:
            break
        # print(
        #     "loaded time %s seq %d" % (df2.iloc[0]["nats_rx_time"], df2.iloc[0]["seq"])
        # )
        if df2["ts"].iloc[0] >= start_time:
            df = pd.concat([df, df2], axis=0)

        timeout = 0.5
    return df


async def load_from_time(
    ui, server, nats_creds, stream_name, export_subject, time, time_range
):
    print("load_from_time" + str(time))
    ui.busy()
    df = await nats_stream_fetch(
        server, nats_creds, stream_name, export_subject, time, time + time_range
    )
    ui.new_data(df)


async def loader(queue, ui, server, nats_creds, stream_name, export_subject):
    time_range = ui.time_range_seconds()
    start_time = datetime.datetime.now() - datetime.timedelta(hours=1)
    start_time = start_time.replace(tzinfo=tz.tzlocal())

    while True:
        print("loader loop loading from %s" % start_time)
        await load_from_time(
            ui, server, nats_creds, stream_name, export_subject, start_time, time_range
        )
        print("loader loop waiting")
        command = await queue.get()
        queue.task_done()
        if command["type"] == "start_date":
            start_time = command["value"]
        elif command["type"] == "time_range":
            time_range = ui.time_range_seconds()
        else:
            print("unknown command")
