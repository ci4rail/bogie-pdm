import csv
import datetime

_decimal_sep = "."


def set_csv_dialect(decimalsep):
    csv.register_dialect(
        "bogie_csv",
        delimiter=";",
        quotechar='"',
        doublequote=True,
        skipinitialspace=True,
        lineterminator="\n",
        quoting=csv.QUOTE_MINIMAL,
    )
    global _decimal_sep
    _decimal_sep = decimalsep


def _float_to_csv_string(fmt, value) -> str:
    global _decimal_sep
    s = "{:.6f}".format(value)
    s = s.replace(".", _decimal_sep)
    return s


class MetaCsv:
    def __init__(self, file_name) -> None:

        file = open(file_name, "w")
        writer = csv.writer(file, dialect="bogie_csv")
        self.writer = writer
        # Write Comment row
        row = [
            "bogie id",
            "trigger type",
            "trigger time",
            "lat",
            "lon",
            "alt",
            "speed",
            "x_max",
            "x_rms",
            "y_max",
            "y_rms",
            "z_max",
            "z_rms",
        ]
        writer.writerow(row)

    def write(self, bogie_data):
        gps_fmt = "{:.6f}"
        sd_fmt = "{:.6f}"

        row = [
            bogie_data.id,
            bogie_data.trigger_type,
            bogie_data.trigger_ts.ToDatetime(),
            _float_to_csv_string(gps_fmt, bogie_data.position.lat),
            _float_to_csv_string(gps_fmt, bogie_data.position.lon),
            _float_to_csv_string(gps_fmt, bogie_data.position.alt),
            _float_to_csv_string(gps_fmt, bogie_data.position.speed),
            _float_to_csv_string(sd_fmt, bogie_data.steady_drive.max[0]),
            _float_to_csv_string(sd_fmt, bogie_data.steady_drive.rms[0]),
            _float_to_csv_string(sd_fmt, bogie_data.steady_drive.max[1]),
            _float_to_csv_string(sd_fmt, bogie_data.steady_drive.rms[1]),
            _float_to_csv_string(sd_fmt, bogie_data.steady_drive.max[2]),
            _float_to_csv_string(sd_fmt, bogie_data.steady_drive.rms[2]),
        ]
        self.writer.writerow(row)

    def close(self):
        self.writer.close()


class SensorCsv:
    def __init__(self, file_name) -> None:

        file = open(file_name, "w")
        writer = csv.writer(file, dialect="bogie_csv")
        self.writer = writer
        # Write Comment row
        row = [
            "sensor_id",
            "timestamp",
            "value",
        ]
        writer.writerow(row)

    def write(self, sensor_data):
        value_fmt = "{:.6f}"

        for sensor in sensor_data:
            id = sensor.sensor_id
            ts = sensor.first_sample_ts.ToDatetime()
            td = 1 / sensor.sample_rate

            for sample in sensor.samples:
                row = [id, ts, _float_to_csv_string(value_fmt, sample)]
                self.writer.writerow(row)
                ts = ts + datetime.timedelta(seconds=td)

        self.writer.writerow(row)

    def close(self):
        self.writer.close()
