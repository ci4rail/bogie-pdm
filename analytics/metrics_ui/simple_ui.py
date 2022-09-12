import ipywidgets as widgets
import ipyleaflet
import matplotlib.pyplot as plt
import datetime
from colour import Color

FIG_SIZE_X = 12
FIG_SIZE = (FIG_SIZE_X, 2)


def make_color(value):
    scale = [x.hex for x in list(Color("red").range_to(Color("green"), 101))]
    if value > 100:
        value = 100
    return scale[int(value)]


def make_geojson(locations):
    """
    locations: list of [ts, lat, lon, value]
    """
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
                            [prev_loc[2], prev_loc[1]],
                            [this_loc[2], this_loc[1]],
                        ],
                    },
                    "properties": {"value": this_loc[3], "ts": this_loc[0]},
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

        # add hover tooltips
        self.tooltip_html = widgets.HTML(value="")
        self.tooltip_html.layout.margin = "0px 20px 20px 20px"
        self.tooltip_html.layout.visibility = "visible"
        # create widget for hover tooltips
        self.hover_control = ipyleaflet.WidgetControl(
            widget=self.tooltip_html, position="bottomright"
        )
        self.map.add_control(self.hover_control)

        self.tab_change(0)
        self.children = [map, tab]

    def handle_tab_change(self, change):
        self.tab_change(change.new)

    def tab_change(self, tabnum):
        if self.map_overlay is not None:
            self.map.remove(self.map_overlay)
            self.map_overlay = None

        l = None
        if tabnum == 0:
            l = self.lte_strength_layer_for_map()
        # elif tabnum == 1:
        #     self.render_acc_on_map()
        elif tabnum == 2:
            l = self.gnss_accuracy_layer_for_map()
        # elif tabnum == 3:
        #     self.render_temp_on_map()
        if l is not None:
            self.map.add_layer(l)
            self.map_overlay = l

    def render_map(self, df):
        center = (49.44, 11.06)
        map = ipyleaflet.Map(center=center, zoom=12)

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

    def lte_strength_layer_for_map(self):
        locations = (
            self.dframe[["ts", "gnss_lat", "gnss_lon", "cellular_strength"]]
            .dropna()
            .values.tolist()
        )

        g = ipyleaflet.GeoJSON(
            data=make_geojson(locations),
            style={"opacity": 0.8, "weight": 5},
            style_callback=self.style_callback_lte,
        )
        g.on_hover(self.overlay_on_hover)
        return g

    def style_callback_lte(self, feature):
        return {"color": make_color(feature["properties"]["value"])}

    def gnss_accuracy_layer_for_map(self):
        locations = (
            self.dframe[["ts", "gnss_lat", "gnss_lon", "gnss_eph"]]
            .dropna()
            .values.tolist()
        )
        g = ipyleaflet.GeoJSON(
            data=make_geojson(locations),
            style={"opacity": 0.8, "weight": 5},
            style_callback=self.style_callback_gnss,
        )
        g.on_hover(self.overlay_on_hover)
        return g

    def style_callback_gnss(self, feature):
        accuracy = feature["properties"]["value"]
        percent = 100 - (accuracy * 10)
        if percent < 0:
            percent = 0
        return {"color": make_color(percent)}

    def overlay_on_hover(self, feature, **kwargs):
        ts = datetime.datetime.fromisoformat(feature["properties"]["ts"])
        timestr = ts.strftime("%H:%M:%S")

        self.tooltip_html.value = f"{timestr} - {feature['properties']['value']}"
        self.tooltip_html.layout.visibility = "visible"
