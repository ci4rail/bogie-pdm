import ipywidgets as widgets
import ipyleaflet
import matplotlib.pyplot as plt
from colour import Color

FIG_SIZE_X = 12
FIG_SIZE = (FIG_SIZE_X, 2)


def make_color(value):
    scale = [x.hex for x in list(Color("red").range_to(Color("green"), 101))]
    if value > 100:
        value = 100
    return scale[int(value)]


def make_geojson(locations):
    features = []
    for i in range(len(locations)):
        if i != 0:
            this_loc = locations[i]
            prev_loc = locations[i - 1]

            features.append(
                {
                    "type": "Feature",
                    "geometry": {
                        "type": "LineString",
                        "coordinates": [
                            [prev_loc[1], prev_loc[0]],
                            [this_loc[1], this_loc[0]],
                        ],
                    },
                    "properties": {"value": this_loc[2]},
                }
            )
    return {"type": "FeatureCollection", "features": features}


class MetricsUi(widgets.VBox):
    def __init__(self, df):
        super().__init__()
        self.dframe = df

        if df.shape[0] == 0:
            print("No DATA!")
            return

        print(
            f"Dataset starting {df.iloc[0]['ts']}  nats rx {df.iloc[0]['nats_rx_time']}"
        )

        lte_out = widgets.Output()
        self.render_lte(lte_out, df)

        acc_out = widgets.Output()
        self.render_acc(acc_out, df)

        gnss_out = widgets.Output()
        self.render_gnss(gnss_out, df)

        temp_out = widgets.Output()
        self.render_temp(temp_out, df)

        tab = widgets.Tab(children=[lte_out, acc_out, gnss_out, temp_out])
        tab.set_title(0, "LTE")
        tab.set_title(1, "Accel")
        tab.set_title(2, "GNSS")
        tab.set_title(3, "Temperature")

        tab.observe(self.handle_tab_change, names="selected_index")

        map = self.render_map(df)
        self.map = map
        self.map_overlay = None

        self.render_lte_on_map()
        self.children = [map, tab]

    def handle_tab_change(self, change):
        print("tab change", change.new)

    def render_map(self, df):
        center = (49.44, 11.06)
        map = ipyleaflet.Map(center=center, zoom=10)

        locations = df[["gnss_lat", "gnss_lon"]].dropna().values.tolist()
        ant_path = ipyleaflet.AntPath(
            locations=locations, delay=1000, color="#7590ba", pulse_color="#3f6fba"
        )
        map.add_layer(ant_path)

        return map

    def render_lte(self, lte_out, df):
        with lte_out:
            fig, ax = plt.subplots(figsize=FIG_SIZE)
            line = ax.plot(df["ts"], df["cellular_strength"])
            ax.set_ylabel("Strength [%]")
            ax.set_xlabel("Time")
            fig.canvas.header_visible = False
            plt.show(fig)

    def render_acc(self, acc_out, df):
        with acc_out:
            fig, (ax1, ax2) = plt.subplots(2, sharex=True, figsize=(FIG_SIZE_X, 6))
            l = ax1.plot(df["ts"], df["accel_x_rms"], label="vertical")
            l = ax1.plot(df["ts"], df["accel_y_rms"], label="side")
            l = ax1.plot(df["ts"], df["accel_z_rms"], label="forward")
            ax1.set_ylabel("RMS Acceleration [g]")
            ax1.set_xlabel("Time")
            ax1.legend()
            l = ax2.plot(df["ts"], df["accel_x_max"], label="vertical")
            l = ax2.plot(df["ts"], df["accel_y_max"], label="side")
            l = ax2.plot(df["ts"], df["accel_z_max"], label="forward")
            ax2.set_ylabel("MAX Acceleration [g]")
            ax2.set_xlabel("Time")
            ax2.legend()
            fig.canvas.header_visible = False
            plt.show(fig)

    def render_gnss(self, out, df):
        with out:
            fig, (ax1, ax2, ax3) = plt.subplots(3, sharex=True, figsize=(FIG_SIZE_X, 6))
            l = ax1.plot(df["ts"], df["gnss_speed"])
            ax1.set_ylabel("Speed [m/s]")
            l = ax2.plot(df["ts"], df["gnss_eph"])
            ax2.set_ylabel("Horizontal Error (m)")
            l = ax3.plot(df["ts"], df["gnss_alt"])
            ax3.set_ylabel("Alt (m)")
            ax3.set_xlabel("Time")
            fig.canvas.header_visible = False
            plt.show(fig)

    def render_temp(self, out, df):
        with out:
            fig, ax = plt.subplots(figsize=FIG_SIZE)
            line = ax.plot(df["ts"], df["temperature_inbox"])
            ax.set_ylabel("Temperature [°C]")
            ax.set_xlabel("Time")
            fig.canvas.header_visible = False
            plt.show(fig)

    def render_lte_on_map(self):
        locations = self.dframe[
            ["gnss_lat", "gnss_lon", "cellular_strength"]
        ].values.tolist()
        g = ipyleaflet.GeoJSON(
            data=make_geojson(locations),
            style={"opacity": 1, "weight": 10},
            style_callback=self.style_callback_lte,
        )
        self.map_overlay = g
        self.map.add_layer(g)

    def style_callback_lte(self, feature):
        return {"color": make_color(feature["properties"]["value"])}
