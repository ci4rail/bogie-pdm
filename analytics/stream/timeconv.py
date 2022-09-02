from dateutil import tz


def pb_timestamp_to_local_datetime(pb_timestamp):
    to_zone = tz.tzlocal()
    from_zone = tz.tzutc()
    utc = pb_timestamp.ToDatetime()
    utc = utc.replace(tzinfo=from_zone)
    return utc.astimezone(to_zone)


def localtime_to_utc(localtime):
    to_zone = tz.tzutc()
    from_zone = tz.tzlocal()
    utc = localtime.astimezone(from_zone)
    return utc.replace(tzinfo=to_zone)
