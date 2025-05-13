# reserv

All important commands can be found in the Makefile. Just type `make help`

# Integration tests

TODO(@perebaj): Integrate the integration tests with the CI pipeline.

To run the integration tests, just type `make integration-test`, and all important containers will be started, and your tests will be executed.

Warning: `make docker-stop` is required to stop the containers after the tests are finished.

## Environment variables

- `POSTGRES_URL`: The URL for the Postgres database.
- `PORT`: The port of the server.
- `LOG_LEVEL`: The level of the logs.
- `LOG_FORMAT`: The format of the logs.
- `CLOUDFLARE_API_KEY`: The API key for the Cloudflare API.
- `CLOUDFLARE_ACCOUNT_ID`: The account ID for the Cloudflare API.
