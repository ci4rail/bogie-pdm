import ipywidgets as widgets
import ipyleaflet
import matplotlib.pyplot as plt
import pandas as pd

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

    def handle_click(self, **kwargs):
        print(kwargs)

    def render_map(self, df):
        center = (49.44, 11.06)
        map = ipyleaflet.Map(center=center, zoom=10)

        locations = df[["lat", "lon"]].dropna().values.tolist()
        ant_path = ipyleaflet.AntPath(
            locations=locations, delay=1000, color="#7590ba", pulse_color="#3f6fba"
        )
        map.add_layer(ant_path)

        i = 1
        for _, row in df.iterrows():
            if not pd.isna(row["lat"]):
                marker = ipyleaflet.Marker(
                    location=(row["lat"], row["lon"]), title=f"Sample {i}"
                )
                marker.on_click(self.handle_click)
                map.add_layer(marker)
            i += 1
        return map

    def scale_sensor_data(self, s):
        return (s - 0.6) * 125

    def render_sensors(self, df):
        out = self.sensors_out
        with out:
            out.clear_output(wait=True)
            fig, (ax1, ax3, ax4) = plt.subplots(
                3,
                sharex=True,
                sharey=True,
                figsize=(FIG_SIZE_X, 6),
            )
            for ax in [ax1, ax3, ax4]:
                ax1.set_ylim(-2, 2)
                ax1.autoscale(enable=False, axis="y")

            x = [(4-i) / 1000 for i in range(0, df["sensor_data"].shape[0])]
            l = ax1.plot(
                x,
                self.scale_sensor_data(df["sensor_data"]["sensor0"]),
                label="Z",
                color="blue",
            )
            ax1.legend(loc="upper right")
            # l = ax2.plot(x, df["sensor_data"]["sensor3"], label="Z links")
            # ax2.legend(loc="upper right")
            l = ax3.plot(
                x,
                self.scale_sensor_data(df["sensor_data"]["sensor1"]),
                label="Y",
                color="orange",
            )
            ax3.legend(loc="upper right")
            l = ax4.plot(
                x,
                self.scale_sensor_data(df["sensor_data"]["sensor2"]),
                label="X",
                color="green",
            )
            ax4.legend(loc="upper right")
            ax4.set_ylabel("Acceleration (g)")
            ax4.set_xlabel("Time (s)")
            fig.canvas.header_visible = False
            timestr = df["trigger_time"].strftime("%Y-%m-%d %H:%M:%S")
            fig.suptitle(
                f"Sample {self.idx+1} of {self.num_samples}, recorded {timestr}",
                fontsize=16,
            )
            plt.show(fig)
