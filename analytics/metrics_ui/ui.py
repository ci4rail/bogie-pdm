import ipywidgets as widgets
import ipyleaflet
import pandas as pd
import datetime
from colour import Color
import matplotlib.pyplot as plt
from IPython.display import display


class MetricsUi(widgets.VBox):
    def __init__(self, queue):
        super().__init__()

        self.queue = queue

        self.start_date_picker = widgets.DatetimePicker()
        self.time_range_picker = widgets.Dropdown(
            options=[
                ("1 min", 60),
                ("5 min", 300),
                ("15 min", 900),
                ("1 hour", 3600),
                ("1 day", 86400),
            ],
            value=300,
            description="Time range",
            disabled=False,
        )
        self.metric_picker = widgets.Dropdown(
            options=["lte", "gnss", "accel", "position"],
            value="position",
            description="Metric",
            disabled=False,
        )
        self.busy_indicator = widgets.Text(disabled=True)
        self.control_box = widgets.Box(
            [
                self.start_date_picker,
                self.time_range_picker,
                self.metric_picker,
                self.busy_indicator,
            ]
        )

        self.plot_outputs = []
        for i in range(2):
            output = widgets.Output()
            self.plot_outputs.append(output)

        tab = widgets.Tab(children=self.plot_outputs)
        tab.set_title(0, "LTE")
        tab.set_title(1, "Accel")

        center = (49.44, 11.06)
        self.map = ipyleaflet.Map(center=center, zoom=10)

        self.time_range_picker.value = 3600

        self.start_date_picker.observe(self.handle_start_date_change, names="value")
        self.time_range_picker.observe(self.handle_time_range_change, names="value")
        self.metric_picker.observe(self.handle_metric_change, names="value")

        self.df = pd.DataFrame()

        self.draw()
        self.children = [self.control_box, self.map, tab]

    def draw_map(self):
        def marker_color(value):
            scale = [x.hex for x in list(Color("red").range_to(Color("green"), 101))]
            if value > 100:
                value = 100
            return scale[int(value)]

        map = self.map
        df = self.df

        for layer in map.layers:
            if not isinstance(layer, ipyleaflet.TileLayer):
                map.remove_layer(layer)

        if self.metric_picker.value == "lte":
            for ts, lat, lon, strength in zip(
                df["ts"], df["gnss_lat"], df["gnss_lon"], df["cellular_strength"]
            ):
                # print("add circle", lat, lon, strength)
                circle = ipyleaflet.Circle(
                    location=(lat, lon),
                    draggable=False,
                    radius=5,
                    color=marker_color(strength),
                )
                map.add_layer(circle)
                # popup = ipyleaflet.Popup(
                #     location=(lat, lon),
                #     child=widgets.HTML(value=f"ts: {ts} strength: {strength}%"),
                #     close_button=False,
                #     auto_close=False,
                #     close_on_escape_key=True,
                # )
                # map.add_layer(popup)

    def draw(self):
        """
        draw dataframe row at idx
        """
        df = self.df
        if df.shape[0] == 0:
            self.busy_indicator.value = "No DATA!"
            return
        self.busy_indicator.value = f"Dataset starting {df.iloc[0]['ts']}"
        print("draw")
        self.draw_map()
        self.draw_lte()
        display(self)

    def draw_lte(self):
        print("draw_lte")
        df = self.df
        output = self.plot_outputs[0]
        with output:
            fig, ax = plt.subplots()
            line = ax.plot(df["ts"], df["cellular_strength"])
            ax.set_title("LTE")
            ax.set_ylabel("Strength [%]")
            ax.set_xlabel("Time")
            plt.show(fig)

    def handle_start_date_change(self, change):
        self.queue.put_nowait({"type": "start_date", "value": change.new})

    def handle_time_range_change(self, change):
        self.queue.put_nowait({"type": "time_range", "value": change.new})

    def handle_metric_change(self, change):
        self.draw()

    def time_range_seconds(self):
        return datetime.timedelta(seconds=self.time_range_picker.value)

    def busy(self):
        self.busy_indicator.value = "Loading..."

    def new_data(self, df):
        print("new_data")
        self.df = df
        self.draw()
