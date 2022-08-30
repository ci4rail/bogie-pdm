import asyncio
import nats
from nats.js.api import ConsumerConfig, DeliverPolicy, AckPolicy, ReplayPolicy


class NatsStream:
    @classmethod
    async def create(
        cls,
        server,
        stream,
        subject,
        durable_name=None,
        delivery_policy=DeliverPolicy.NEW,
        opt_start_time=None,
        opt_start_seq=None,
        replayPolicy=ReplayPolicy.INSTANT,
    ):
        self = NatsStream()
        nc = await nats.connect(server)

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
    async def from_start_time(cls, server, stream, subject, start_time):
        """
        Create a new ephemeral NatsStream object that starts from a given time.
        """
        return await cls.create(
            server,
            stream,
            subject,
            delivery_policy=DeliverPolicy.BY_START_TIME,
            opt_start_time=start_time,
        )

    @classmethod
    async def from_durable_all(cls, server, stream, subject, durable_name):
        return await cls.create(
            server,
            stream,
            subject,
            durable_name=durable_name,
            delivery_policy=DeliverPolicy.ALL,
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
