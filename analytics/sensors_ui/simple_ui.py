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
        self.num_samples = df.shape[0]
        self.slider = widgets.IntSlider(
            min=0,
            max=self.num_samples - 1,
            step=1,
            readout=False,
        )
        self.idx = 0
        self.render_sensors(df.iloc[self.idx])

        self.slider.observe(self.handle_slider_change, names="value")
        self.children = [self.map, self.slider, self.sensors_out]

    def handle_slider_change(self, change):
        self.idx = change.new
        self.render_sensors(self.dframe.iloc[self.idx])

    def render_map(self, df):
        center = (49.44, 11.06)
        map = ipyleaflet.Map(center=center, zoom=10)

        locations = df[["lat", "lon"]].dropna().values.tolist()
        ant_path = ipyleaflet.AntPath(
            locations=locations, delay=1000, color="#7590ba", pulse_color="#3f6fba"
        )
        map.add_layer(ant_path)

        return map

    def render_sensors(self, df):
        out = self.sensors_out
        with out:
            out.clear_output(wait=True)
            fig, (ax1, ax2, ax3, ax4) = plt.subplots(
                4, sharex=True, figsize=(FIG_SIZE_X, 6)
            )
            # fig, ax = plt.subplots(figsize=FIG_SIZE)
            x = list(range(0, df["sensor_data"].shape[0]))
            ax1.plot(x, df["sensor_data"]["sensor0"], label="Z rechts")
            ax1.legend(loc="upper right")
            ax2.plot(x, df["sensor_data"]["sensor3"], label="Z links")
            ax2.legend(loc="upper right")
            ax3.plot(x, df["sensor_data"]["sensor1"], label="Y rechts")
            ax3.legend(loc="upper right")
            ax4.plot(x, df["sensor_data"]["sensor2"], label="X mitte")
            ax4.legend(loc="upper right")
            ax4.set_ylabel("Acceleration")
            ax4.set_xlabel("Time")
            fig.canvas.header_visible = False
            fig.suptitle(
                f"Sample {self.idx+1} of {self.num_samples} {df['trigger_time']}",
                fontsize=16,
            )
            plt.show(fig)
