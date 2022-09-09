import stream.stream as stream
import stream.bogiepandas as bogiepandas
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
        df2 = bogiepandas.bogie_nats_to_pandas(msg)

        ts = df2["nats_rx_time"].iloc[0].replace(tzinfo=tz.tzlocal())
        if ts > end_time:
            break
        # print(
        #     "loaded time %s seq %d" % (df2.iloc[0]["nats_rx_time"], df2.iloc[0]["seq"])
        # )
        if df2["trigger_time"].iloc[0] >= start_time:
            df = pd.concat([df, df2], axis=0)

        timeout = 0.5
    await ns.close()
    return df
