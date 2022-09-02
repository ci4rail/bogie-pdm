from dateutil import tz


def pb_timestamp_to_local_datetime(pb_timestamp):
    to_zone = tz.tzlocal()
    from_zone = tz.tzutc()
    utc = pb_timestamp.ToDatetime()
    utc = utc.replace(tzinfo=from_zone)
    return utc.astimezone(to_zone)


def localtime_to_utc(localtime):
    return localtime.astimezone(tz.tzutc())
