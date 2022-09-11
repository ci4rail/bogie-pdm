import proto.bogie_pb2 as bogie_pb2
import pandas as pd
from dateutil import tz
import stream.timeconv as timeconv
import json
import base64


def bogie_nats_to_pandas(m):
    """Convert bogie message from nats to pandas dataframe"""

    df = pd.DataFrame()

    df["nats_rx_time"] = [m.metadata.timestamp]
    df["seq"] = [m.metadata.sequence.stream]

    js = json.loads(m.data)
    payload = base64.b64decode(js["data_base64"])  # strip dapr header
    data = bogie_pb2.Bogie()
    data.ParseFromString(payload)

    df["trigger_time"] = [timeconv.pb_timestamp_to_local_datetime(data.trigger_ts)]

    sensor_df = pd.DataFrame()

    min_len = 1e9
    for sensor in data.sensor_samples:
        if len(sensor.samples) < min_len:
            min_len = len(sensor.samples)

    for sensor in data.sensor_samples:
        sensor_df["sensor%d" % sensor.sensor_id] = list(sensor.samples[0:min_len])

    df["sensor_data"] = [sensor_df]

    df["bogie_id"] = [data.id]
    try:
        has_invalid = data.HasField("invalid")
    except ValueError:
        has_invalid = False
    if not has_invalid:
        df["invalid"] = [False]
    if not df["invalid"].iloc[0]:
        df["lat"] = [data.position.lat]
        df["lon"] = [data.position.lon]
        df["alt"] = [data.position.alt]
        df["speed"] = [data.position.speed]

    return df
