import proto.bogie_pb2 as bogie_pb2
import pandas as pd
from dateutil import tz


def pb_timestamp_to_local_datetime(pb_timestamp):
    to_zone = tz.tzlocal()
    from_zone = tz.tzutc()
    utc = pb_timestamp.ToDatetime()
    utc = utc.replace(tzinfo=from_zone)
    return utc.astimezone(to_zone)


def bogie_nats_to_pandas(m):
    """Convert bogie message from nats to pandas dataframe"""

    df = pd.DataFrame()

    df["nats_rx_time"] = [m.metadata.timestamp]
    df["seq"] = [m.metadata.sequence.stream]
    data = bogie_pb2.Bogie()
    data.ParseFromString(m.data)

    df["trigger_time"] = [pb_timestamp_to_local_datetime(data.trigger_ts)]

    sensor_df = pd.DataFrame()
    for sensor in data.sensor_samples:
        sensor_df["sensor%d" % sensor.sensor_id] = list(sensor.samples)

    df["sensor_data"] = [sensor_df]
    return df
