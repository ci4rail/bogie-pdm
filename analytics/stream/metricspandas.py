import proto.metrics_pb2 as metrics_pb2
import pandas as pd
import stream.timeconv as timeconv
import json
import base64


def metrics_nats_to_pandas(m):
    """Convert metrics message from nats to pandas dataframe"""

    df = pd.DataFrame()

    df["nats_rx_time"] = [m.metadata.timestamp]
    df["seq"] = [m.metadata.sequence.stream]

    # print("message data: %s" % m.data)

    js = json.loads(m.data)
    payload = base64.b64decode(js["data_base64"])  # strip dapr header
    data = metrics_pb2.Metrics()
    data.ParseFromString(payload)

    df["ts"] = [timeconv.pb_timestamp_to_local_datetime(data.ts)]
    df["accel_x_rms"] = [data.steady_drive.rms[0]]
    df["accel_y_rms"] = [data.steady_drive.rms[1]]
    df["accel_z_rms"] = [data.steady_drive.rms[2]]
    df["accel_x_max"] = [data.steady_drive.max[0]]
    df["accel_y_max"] = [data.steady_drive.max[1]]
    df["accel_z_max"] = [data.steady_drive.max[2]]
    df["gnss_mode"] = [data.gnss_raw.mode]
    if df["gnss_mode"].iloc[0] > 1:
        df["gnss_lat"] = [data.gnss_raw.lat]
        df["gnss_lon"] = [data.gnss_raw.lon]
        df["gnss_alt"] = [data.gnss_raw.alt]
        df["gnss_speed"] = [data.gnss_raw.speed]
        df["gnss_eph"] = [data.gnss_raw.eph]
    df["gnss_numsats"] = [data.gnss_raw.numsats]
    df["cellular_strength"] = [data.cellular.strength]
    df["cellular_operator"] = [data.cellular.operator]
    df["temperature_inbox"] = [data.temperature.inBox]
    df["internet_connected"] = [data.internet.connected]
    return df
