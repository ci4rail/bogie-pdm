import proto.bogie_pb2 as bogie_pb2
import pandas as pd


def bogie_nats_to_pandas(m):
    """Convert bogie message from nats to pandas dataframe"""

    df = pd.DataFrame()

    df["nats_rx_time"] = [m.metadata.timestamp]
    data = bogie_pb2.Bogie()
    data.ParseFromString(m.data)

    df["trigger_time"] = [data.trigger_ts.ToDatetime()]

    sensor_df = pd.DataFrame()
    for sensor in data.sensor_samples:
        sensor_df["sensor%d" % sensor.sensor_id] = list(sensor.samples)

    df["sensor_data"] = [sensor_df]
    return df
