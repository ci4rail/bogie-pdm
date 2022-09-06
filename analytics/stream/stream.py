import asyncio
import nats
from nats.js.api import ConsumerConfig, DeliverPolicy, AckPolicy, ReplayPolicy
import stream.timeconv as timeconv


class NatsStream:
    @classmethod
    async def create(
        cls,
        server,
        credsfile_path,
        stream,
        subject,
        durable_name=None,
        delivery_policy=DeliverPolicy.NEW,
        opt_start_time=None,
        opt_start_seq=None,
        replayPolicy=ReplayPolicy.INSTANT,
    ):
        self = NatsStream()

        options = {
            "servers": [server],
        }

        if credsfile_path:
            options["user_credentials"] = credsfile_path

        nc = await nats.connect(**options)

        # Create JetStream context.
        js = nc.jetstream()

        config = ConsumerConfig(
            durable_name=durable_name,
            deliver_policy=delivery_policy,
            opt_start_time=opt_start_time,
            opt_start_seq=opt_start_seq,
            ack_policy=AckPolicy.EXPLICIT,
            replay_policy=replayPolicy,
        )

        sub = await js.subscribe(subject, stream=stream, config=config)

        self.sub = sub
        self.nc = nc
        return self

    @classmethod
    async def from_start_time(cls, server, credsfile_path, stream, subject, start_time):
        """
        Create a new ephemeral NatsStream object that starts from a given time.
        start_time must be a Datetime object in local time
        """
        start_time = timeconv.localtime_to_utc(start_time)
        stime = start_time.strftime("%Y-%m-%dT%H:%M:%S.%fZ")
        #print("start time %s" % stime)

        return await cls.create(
            server,
            credsfile_path,
            stream,
            subject,
            delivery_policy=DeliverPolicy.BY_START_TIME,
            opt_start_time=stime,
        )

    @classmethod
    async def from_seq(cls, server, credsfile_path, stream, subject, seq):
        """
        Create a new ephemeral NatsStream object that starts from a given sequence.
        """
        return await cls.create(
            server,
            credsfile_path,
            stream,
            subject,
            delivery_policy=DeliverPolicy.BY_START_SEQUENCE,
            opt_start_seq=int(seq),
        )

    async def next_msg(self, timeout=2.0):
        """
        Wait for next message.
        Returns: Message object or None if timeout.
        """
        try:
            msg = await self.sub.next_msg(timeout=timeout)
        except asyncio.TimeoutError:
            return None
        return msg

    async def ack(self, msg):
        """
        Acknowledge message.
        """
        await msg.ack()

    async def close(self):
        """
        Close connection.
        """
        await self.nc.close()
