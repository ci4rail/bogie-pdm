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

    min_len = 1E9
    for sensor in data.sensor_samples:
        if len(sensor.samples) < min_len:
            min_len = len(sensor.samples)

    for sensor in data.sensor_samples:
        sensor_df["sensor%d" % sensor.sensor_id] = list(sensor.samples[0:min_len])

    df["sensor_data"] = [sensor_df]
    return df
