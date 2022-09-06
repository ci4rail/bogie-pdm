import ipywidgets as widgets
import ipyleaflet
import matplotlib.pyplot as plt

FIG_SIZE_X = 12
FIG_SIZE = (FIG_SIZE_X, 2)
NUM_SENSORS = 4


class SensorsUi(widgets.VBox):
    def __init__(self, df):
        super().__init__()
        self.dframe = df

        if df.shape[0] == 0:
            print("No DATA!")
            return

        self.sensors_out = widgets.Output()
        self.map = self.render_map(df)
        self.slider = widgets.IntSlider(
            min=0,
            max=df.shape[0] - 1,
            step=1,
            readout=False,
        )
        self.render_sensors(self.sensors_out, df.iloc[0])

        self.slider.observe(self.handle_slider_change, names="value")
        self.children = [self.map, self.slider, self.sensors_out]

    def handle_slider_change(self, change):
        idx = change.new
        self.render_sensors(self.sensors_out, self.dframe.iloc[idx])

    def render_map(self, df):
        center = (49.44, 11.06)
        map = ipyleaflet.Map(center=center, zoom=10)

        locations = df[["lat", "lon"]].dropna().values.tolist()
        ant_path = ipyleaflet.AntPath(
            locations=locations, delay=1000, color="#7590ba", pulse_color="#3f6fba"
        )
        map.add_layer(ant_path)

        return map

    def render_sensors(self, out, df):
        out.outputs = []
        print("render_sensors %s" % (df["trigger_time"]))
        with out:
            fig, ax = plt.subplots(figsize=FIG_SIZE)
            x = list(range(0, df["sensor_data"].shape[0]))
            for i in range(NUM_SENSORS):
                l = ax.plot(x, df["sensor_data"]["sensor%d" % i], label="sensor%d" % i)
            ax.set_ylabel("Acceleration")
            ax.set_xlabel("Time")
            ax.legend()
            fig.canvas.header_visible = False
            # fig.canvas.draw()
            plt.show(fig)
