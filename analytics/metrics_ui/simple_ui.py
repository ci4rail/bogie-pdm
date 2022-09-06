import ipywidgets as widgets
import ipyleaflet
import matplotlib.pyplot as plt

FIG_SIZE_X = 12
FIG_SIZE = (FIG_SIZE_X, 2)


def marker_color(value):
    scale = [x.hex for x in list(Color("red").range_to(Color("green"), 101))]
    if value > 100:
        value = 100
    return scale[int(value)]


def render_map(df):
    center = (49.44, 11.06)
    map = ipyleaflet.Map(center=center, zoom=10)

    locations = df[["gnss_lat", "gnss_lon"]].dropna().values.tolist()
    ant_path = ipyleaflet.AntPath(
        locations=locations, delay=1000, color="#7590ba", pulse_color="#3f6fba"
    )
    map.add_layer(ant_path)

    return map


def render_lte(lte_out, df):
    with lte_out:
        fig, ax = plt.subplots(figsize=FIG_SIZE)
        line = ax.plot(df["ts"], df["cellular_strength"])
        ax.set_ylabel("Strength [%]")
        ax.set_xlabel("Time")
        fig.canvas.header_visible = False
        plt.show(fig)


def render_acc(acc_out, df):
    with acc_out:
        fig, ax = plt.subplots(figsize=FIG_SIZE)
        l = ax.plot(df["ts"], df["accel_x_rms"], label="vertical")
        l = ax.plot(df["ts"], df["accel_y_rms"], label="side")
        l = ax.plot(df["ts"], df["accel_z_rms"], label="forward")
        ax.set_ylabel("Acceleration [m/s^2]")
        ax.set_xlabel("Time")
        ax.legend()
        fig.canvas.header_visible = False
        plt.show(fig)


def render_gnss(out, df):
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


def render_temp(out, df):
    with out:
        fig, ax = plt.subplots(figsize=FIG_SIZE)
        line = ax.plot(df["ts"], df["temperature_inbox"])
        ax.set_ylabel("Temperature [°C]")
        ax.set_xlabel("Time")
        fig.canvas.header_visible = False
        plt.show(fig)


def render_ui(df):
    if df.shape[0] == 0:
        print("No DATA!")
        return

    print(f"Dataset starting {df.iloc[0]['ts']}")

    lte_out = widgets.Output()
    render_lte(lte_out, df)

    acc_out = widgets.Output()
    render_acc(acc_out, df)

    gnss_out = widgets.Output()
    render_gnss(gnss_out, df)

    temp_out = widgets.Output()
    render_temp(temp_out, df)

    tab = widgets.Tab(children=[lte_out, acc_out, gnss_out, temp_out])
    tab.set_title(0, "LTE")
    tab.set_title(1, "Accel")
    tab.set_title(2, "GNSS")
    tab.set_title(3, "Temperature")

    map = render_map(df)
    vbox = widgets.VBox([map, tab])
    return vbox
