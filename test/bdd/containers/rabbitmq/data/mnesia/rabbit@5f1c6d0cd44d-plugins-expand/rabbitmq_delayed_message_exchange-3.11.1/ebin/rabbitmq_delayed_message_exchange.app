{application, 'rabbitmq_delayed_message_exchange', [
	{description, "RabbitMQ Delayed Message Exchange"},
	{vsn, "3.11.1"},
	{modules, ['rabbit_delayed_message','rabbit_delayed_message_app','rabbit_delayed_message_sup','rabbit_delayed_message_utils','rabbit_exchange_type_delayed_message']},
	{registered, [rabbitmq_delayed_message_exchange_sup]},
	{applications, [kernel,stdlib,rabbit_common,rabbit]},
	{mod, {rabbit_delayed_message_app, []}},
	{env, []},
		{broker_version_requirements, ["3.11.0"]}
]}.