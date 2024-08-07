ARG PLUGIN_VERSION=3.12.0
ARG BASE_VERSION=3.12

FROM ubuntu:20.04 AS builder

ARG PLUGIN_VERSION

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y curl
RUN mkdir -p /plugins && \
    curl -fsSL \
    -o "/plugins/rabbitmq_delayed_message_exchange-${PLUGIN_VERSION}.ez" \
    https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/v${PLUGIN_VERSION}/rabbitmq_delayed_message_exchange-${PLUGIN_VERSION}.ez

FROM rabbitmq:${BASE_VERSION}-management-alpine

ARG PLUGIN_VERSION

# Copy plugin from builder stage
COPY --from=builder /plugins/rabbitmq_delayed_message_exchange-${PLUGIN_VERSION}.ez $RABBITMQ_HOME/plugins/

# Set the correct permissions for the .erlang.cookie
RUN chmod 600 /var/lib/rabbitmq/.erlang.cookie

# Enable RabbitMQ delayed message exchange plugin
RUN rabbitmq-plugins enable --offline rabbitmq_delayed_message_exchange

# Stage 2: Setup RabbitMQ with Go
# Ensure the .erlang.cookie has correct permissions again, as this may be reset
RUN chmod 600 /var/lib/rabbitmq/.erlang.cookie

# Command to start RabbitMQ and the Go application
CMD ["sh", "-c", "rabbitmq-server"]