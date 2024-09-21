# Url shortener

An url shortener made with Go.

# Local development

To run the server locally, create an .env file and set the environment variables.

```bash
cp .env.example .env
```

Then run the server with the following command:

```bash
go run url-shortener
```

Navigate to `http://localhost:8080` to access the server.

# Deployment

The server is deployed at the following address: https://url-shortener-p3em.onrender.com/

# Sleep preventer

The server is hosted on [Render](https://render.com/).

"Render spins down a Free web service that goes 15 minutes without receiving inbound traffic." ([source](https://docs.render.com/free#free-web-services))

~~To prevent the server from sleeping, a GitHub action is used to make an HTTP request to the server every 10 minutes (`.github/workflows/sleep-preventer.yml`).~~

To prevent the server from sleeping, a cron job is set up to make an HTTP request to the server every 10 minutes.
