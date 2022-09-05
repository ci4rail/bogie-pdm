
import ipywidgets as widgets
import ipyleaflet
import pandas as pd
import plotly.graph_objects as go
import datetime
from colour import Color

class MetricsUi(widgets.VBox):
    def __init__(self, queue):
        super().__init__()
        
        self.queue = queue
        
        self.start_date_picker = widgets.DatetimePicker()
        self.time_range_picker = widgets.Dropdown(
            options=[('1 min', 60), ('5 min', 300), ('15 min', 900), ('1 hour', 3600), ('1 day', 86400)],
            value=300,
            description='Time range',
            disabled=False,
        )
        self.metric_picker = widgets.Dropdown(
            options=['lte', 'gnss', 'accel', 'position'],
            value='position',
            description='Metric',
            disabled=False,
        )
        self.busy_indicator = widgets.Text(
            disabled=True
        )
        self.control_box = widgets.Box([self.start_date_picker, self.time_range_picker, self.metric_picker, self.busy_indicator])

        center = (49.44, 11.06)
        self.map = ipyleaflet.Map(center=center, zoom=10)

        self.accel_wf_fig = go.FigureWidget()
        for i in range(3):
            self.accel_wf_fig.add_trace(go.Scatter(x=[], y=[], name="sensor%d" % i))

        self.time_range_picker.value = 300

        self.start_date_picker.observe(self.handle_start_date_change, names='value')
        self.time_range_picker.observe(self.handle_time_range_change, names='value')
        self.metric_picker.observe(self.handle_metric_change, names='value')

        self.df = pd.DataFrame()

        self.draw()
        self.children = [self.control_box, self.map]

    def draw_map(self):
        def marker_color(value):
            scale = [
                x.hex
                for x in list(Color("red").range_to(Color("green"), 101))
            ]
            if value > 100:
                value = 100
            return scale[int(value)]

        map = self.map
        df = self.df

        for layer in map.layers:
            if not isinstance(layer, ipyleaflet.TileLayer):
                map.remove_layer(layer)
        
        if self.metric_picker.value == "lte":
            for ts, lat, lon, strength in zip(df["ts"], df["gnss_lat"], df["gnss_lon"], df["cellular_strength"]):
                #print("add circle", lat, lon, strength)
                circle = ipyleaflet.Circle(location=(lat, lon), draggable=False, radius=5, color=marker_color(strength))
                map.add_layer(circle)
                # popup = ipyleaflet.Popup(
                #     location=(lat, lon),
                #     child=widgets.HTML(value=f"ts: {ts} strength: {strength}%"),
                #     close_button=False,
                #     auto_close=False,
                #     close_on_escape_key=True,
                # )
                # map.add_layer(popup)                  
                


    def draw_accel(self):
        df = self.df
        fig = self.accel_wf_fig
        
        #fig.update_layout({"title": f"Dataset from {df.iloc[0]['ts']} seq {df.iloc[0]['seq']}"})
        fig.update_layout({"title": "Accelerometer"})

        for i in range(3):
            fig.data[i].x = df["ts"]
        
        fig.data[0].y = df["accel_x_rms"]
        fig.data[1].y = df["accel_y_rms"]
        fig.data[2].y = df["accel_z_rms"]


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
        #self.draw_accel()


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
        self.df = df
        self.draw()