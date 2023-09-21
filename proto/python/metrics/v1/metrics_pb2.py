# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: metrics/v1/metrics.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='metrics/v1/metrics.proto',
  package='metrics',
  syntax='proto3',
  serialized_options=_b('Z\nmetrics/v1'),
  serialized_pb=_b('\n\x18metrics/v1/metrics.proto\x12\x07metrics\x1a\x1fgoogle/protobuf/timestamp.proto\"\xaf\x05\n\x07Metrics\x12&\n\x02ts\x18\x01 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x32\n\x0csteady_drive\x18\x02 \x01(\x0b\x32\x1c.metrics.Metrics.SteadyDrive\x12+\n\x08position\x18\x03 \x01(\x0b\x32\x19.metrics.Metrics.Position\x12\x31\n\x0btemperature\x18\x04 \x01(\x0b\x32\x1c.metrics.Metrics.Temperature\x12+\n\x08internet\x18\x05 \x01(\x0b\x32\x19.metrics.Metrics.Internet\x12+\n\x08\x63\x65llular\x18\x06 \x01(\x0b\x32\x19.metrics.Metrics.Cellular\x12*\n\x08gnss_raw\x18\x07 \x01(\x0b\x32\x18.metrics.Metrics.GnssRaw\x1a\'\n\x0bSteadyDrive\x12\x0b\n\x03max\x18\x01 \x03(\x01\x12\x0b\n\x03rms\x18\x02 \x03(\x01\x1aO\n\x08Position\x12\r\n\x05valid\x18\x01 \x01(\x08\x12\x0b\n\x03lat\x18\x02 \x01(\x02\x12\x0b\n\x03lon\x18\x03 \x01(\x02\x12\x0b\n\x03\x61lt\x18\x04 \x01(\x02\x12\r\n\x05speed\x18\x05 \x01(\x02\x1a\x1c\n\x0bTemperature\x12\r\n\x05inBox\x18\x01 \x01(\x02\x1a\x1d\n\x08Internet\x12\x11\n\tconnected\x18\x01 \x01(\x08\x1a>\n\x08\x43\x65llular\x12\x10\n\x08operator\x18\x01 \x01(\t\x12\x10\n\x08strength\x18\x02 \x01(\x02\x12\x0e\n\x06\x63\x65llid\x18\x03 \x01(\t\x1ak\n\x07GnssRaw\x12\x0b\n\x03lat\x18\x01 \x01(\x02\x12\x0b\n\x03lon\x18\x02 \x01(\x02\x12\x0b\n\x03\x61lt\x18\x03 \x01(\x02\x12\r\n\x05speed\x18\x04 \x01(\x02\x12\x0b\n\x03\x65ph\x18\x05 \x01(\x02\x12\x0c\n\x04mode\x18\x06 \x01(\x05\x12\x0f\n\x07numsats\x18\x07 \x01(\x05\x42\x0cZ\nmetrics/v1b\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_timestamp__pb2.DESCRIPTOR,])




_METRICS_STEADYDRIVE = _descriptor.Descriptor(
  name='SteadyDrive',
  full_name='metrics.Metrics.SteadyDrive',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='max', full_name='metrics.Metrics.SteadyDrive.max', index=0,
      number=1, type=1, cpp_type=5, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='rms', full_name='metrics.Metrics.SteadyDrive.rms', index=1,
      number=2, type=1, cpp_type=5, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=404,
  serialized_end=443,
)

_METRICS_POSITION = _descriptor.Descriptor(
  name='Position',
  full_name='metrics.Metrics.Position',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='valid', full_name='metrics.Metrics.Position.valid', index=0,
      number=1, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='lat', full_name='metrics.Metrics.Position.lat', index=1,
      number=2, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='lon', full_name='metrics.Metrics.Position.lon', index=2,
      number=3, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='alt', full_name='metrics.Metrics.Position.alt', index=3,
      number=4, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='speed', full_name='metrics.Metrics.Position.speed', index=4,
      number=5, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=445,
  serialized_end=524,
)

_METRICS_TEMPERATURE = _descriptor.Descriptor(
  name='Temperature',
  full_name='metrics.Metrics.Temperature',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='inBox', full_name='metrics.Metrics.Temperature.inBox', index=0,
      number=1, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=526,
  serialized_end=554,
)

_METRICS_INTERNET = _descriptor.Descriptor(
  name='Internet',
  full_name='metrics.Metrics.Internet',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='connected', full_name='metrics.Metrics.Internet.connected', index=0,
      number=1, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=556,
  serialized_end=585,
)

_METRICS_CELLULAR = _descriptor.Descriptor(
  name='Cellular',
  full_name='metrics.Metrics.Cellular',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='operator', full_name='metrics.Metrics.Cellular.operator', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='strength', full_name='metrics.Metrics.Cellular.strength', index=1,
      number=2, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='cellid', full_name='metrics.Metrics.Cellular.cellid', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=587,
  serialized_end=649,
)

_METRICS_GNSSRAW = _descriptor.Descriptor(
  name='GnssRaw',
  full_name='metrics.Metrics.GnssRaw',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='lat', full_name='metrics.Metrics.GnssRaw.lat', index=0,
      number=1, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='lon', full_name='metrics.Metrics.GnssRaw.lon', index=1,
      number=2, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='alt', full_name='metrics.Metrics.GnssRaw.alt', index=2,
      number=3, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='speed', full_name='metrics.Metrics.GnssRaw.speed', index=3,
      number=4, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='eph', full_name='metrics.Metrics.GnssRaw.eph', index=4,
      number=5, type=2, cpp_type=6, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='mode', full_name='metrics.Metrics.GnssRaw.mode', index=5,
      number=6, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='numsats', full_name='metrics.Metrics.GnssRaw.numsats', index=6,
      number=7, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=651,
  serialized_end=758,
)

_METRICS = _descriptor.Descriptor(
  name='Metrics',
  full_name='metrics.Metrics',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='ts', full_name='metrics.Metrics.ts', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='steady_drive', full_name='metrics.Metrics.steady_drive', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='position', full_name='metrics.Metrics.position', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='temperature', full_name='metrics.Metrics.temperature', index=3,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='internet', full_name='metrics.Metrics.internet', index=4,
      number=5, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='cellular', full_name='metrics.Metrics.cellular', index=5,
      number=6, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='gnss_raw', full_name='metrics.Metrics.gnss_raw', index=6,
      number=7, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_METRICS_STEADYDRIVE, _METRICS_POSITION, _METRICS_TEMPERATURE, _METRICS_INTERNET, _METRICS_CELLULAR, _METRICS_GNSSRAW, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=71,
  serialized_end=758,
)

_METRICS_STEADYDRIVE.containing_type = _METRICS
_METRICS_POSITION.containing_type = _METRICS
_METRICS_TEMPERATURE.containing_type = _METRICS
_METRICS_INTERNET.containing_type = _METRICS
_METRICS_CELLULAR.containing_type = _METRICS
_METRICS_GNSSRAW.containing_type = _METRICS
_METRICS.fields_by_name['ts'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_METRICS.fields_by_name['steady_drive'].message_type = _METRICS_STEADYDRIVE
_METRICS.fields_by_name['position'].message_type = _METRICS_POSITION
_METRICS.fields_by_name['temperature'].message_type = _METRICS_TEMPERATURE
_METRICS.fields_by_name['internet'].message_type = _METRICS_INTERNET
_METRICS.fields_by_name['cellular'].message_type = _METRICS_CELLULAR
_METRICS.fields_by_name['gnss_raw'].message_type = _METRICS_GNSSRAW
DESCRIPTOR.message_types_by_name['Metrics'] = _METRICS
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Metrics = _reflection.GeneratedProtocolMessageType('Metrics', (_message.Message,), dict(

  SteadyDrive = _reflection.GeneratedProtocolMessageType('SteadyDrive', (_message.Message,), dict(
    DESCRIPTOR = _METRICS_STEADYDRIVE,
    __module__ = 'metrics.v1.metrics_pb2'
    # @@protoc_insertion_point(class_scope:metrics.Metrics.SteadyDrive)
    ))
  ,

  Position = _reflection.GeneratedProtocolMessageType('Position', (_message.Message,), dict(
    DESCRIPTOR = _METRICS_POSITION,
    __module__ = 'metrics.v1.metrics_pb2'
    # @@protoc_insertion_point(class_scope:metrics.Metrics.Position)
    ))
  ,

  Temperature = _reflection.GeneratedProtocolMessageType('Temperature', (_message.Message,), dict(
    DESCRIPTOR = _METRICS_TEMPERATURE,
    __module__ = 'metrics.v1.metrics_pb2'
    # @@protoc_insertion_point(class_scope:metrics.Metrics.Temperature)
    ))
  ,

  Internet = _reflection.GeneratedProtocolMessageType('Internet', (_message.Message,), dict(
    DESCRIPTOR = _METRICS_INTERNET,
    __module__ = 'metrics.v1.metrics_pb2'
    # @@protoc_insertion_point(class_scope:metrics.Metrics.Internet)
    ))
  ,

  Cellular = _reflection.GeneratedProtocolMessageType('Cellular', (_message.Message,), dict(
    DESCRIPTOR = _METRICS_CELLULAR,
    __module__ = 'metrics.v1.metrics_pb2'
    # @@protoc_insertion_point(class_scope:metrics.Metrics.Cellular)
    ))
  ,

  GnssRaw = _reflection.GeneratedProtocolMessageType('GnssRaw', (_message.Message,), dict(
    DESCRIPTOR = _METRICS_GNSSRAW,
    __module__ = 'metrics.v1.metrics_pb2'
    # @@protoc_insertion_point(class_scope:metrics.Metrics.GnssRaw)
    ))
  ,
  DESCRIPTOR = _METRICS,
  __module__ = 'metrics.v1.metrics_pb2'
  # @@protoc_insertion_point(class_scope:metrics.Metrics)
  ))
_sym_db.RegisterMessage(Metrics)
_sym_db.RegisterMessage(Metrics.SteadyDrive)
_sym_db.RegisterMessage(Metrics.Position)
_sym_db.RegisterMessage(Metrics.Temperature)
_sym_db.RegisterMessage(Metrics.Internet)
_sym_db.RegisterMessage(Metrics.Cellular)
_sym_db.RegisterMessage(Metrics.GnssRaw)


DESCRIPTOR._options = None
# @@protoc_insertion_point(module_scope)
