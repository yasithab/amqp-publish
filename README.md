# RabbitMQ Message Publisher

This is a command-line tool for publishing messages to a RabbitMQ server. It connects to a RabbitMQ instance and publishes messages specified via command-line flags or input file.

## Usage

```shell
./rabbitmq-publisher [flags]
```

### Flags

- `-uri string`: AMQP URI (e.g., `amqp://<user>:<password>@<host>:port>/[vhost]`), required.
- `-exchange string`: Exchange name (Optional).
- `-routing-key string`: Routing key. Use a queue name with a blank exchange to publish directly to a queue.
- `-body string`: Message body.
- `-input-file string`: Input file path containing messages, one per line. If specified, it overrides the `-body` flag.
- `-persistent`: Use persistent delivery mode (default is non-persistent).

### Note

The application requires that the `uri` is provided. Either the `exchange` or the `routing-key` must be specified. Similarly, you must provide either a `body` or an `input-file`.

## Example

To publish a single message directly using the `body`:

```shell
./rabbitmq-publisher --uri="amqp://guest:guest@localhost:5672/" --routing-key="myKey" --body="Hello, World!"
```

To publish multiple messages from a file, where each line in the file is treated as a separate message:

```shell
./rabbitmq-publisher --uri="amqp://guest:guest@localhost:5672/" --routing-key="myKey" --input-file="data.txt"
```

## Error Handling

- The application will exit with an error message if the `uri` is not provided, or if both `exchange` and `routing-key`, or both `body` and `input-file` are missing.
- Other errors during message reading or publishing will be logged to the console. Messages will continue to be published unless there's a failure in setting up the connection or channel.

## Dependencies

This application requires the `github.com/rabbitmq/amqp091-go` package to handle AMQP protocol interactions.

### Installation

To use this application, ensure you have Go installed and then install the dependency using:

```shell
go get github.com/rabbitmq/amqp091-go
```

Compile the package using:

```shell
go build -o rabbitmq-publisher
```

This generates an executable you can use to publish messages to RabbitMQ.

## Acknowledgments

This application was inspired by and built upon the work from the [amqp-publish](https://github.com/selency/amqp-publish) GitHub repository by the Selency team. Their project provided the foundation and insights necessary to develop this tool.

## License

This project is licensed under [MIT License](LICENSE). 