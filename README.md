# Url shortener

An url shortener made with Go.

# Deployment

The server is deployed at the following address: https://url-shortener-p3em.onrender.com/

# Sleep preventer

The server is hosted on [Render](https://render.com/).

"Render spins down a Free web service that goes 15 minutes without receiving inbound traffic." ([source](https://docs.render.com/free#free-web-services))

To prevent the server from sleeping, a GitHub action was used to make an HTTP request to the server every 10 minutes (`.github/workflows/sleep-preventer.yml`).

A cron job is now set up to replace the GitHub action.
