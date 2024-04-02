# Url shortener

An url shortener made with Go.

# Deployment

The server is deployed at the following address: https://url-shortener-p3em.onrender.com/

# Sleep preventer

The server is hosted on [Render](https://render.com/).

"Render spins down a Free web service that goes 15 minutes without receiving inbound traffic." ([source](https://docs.render.com/free#free-web-services))

To prevent the server from sleeping, a cron job is set up to make an HTTP request to the server every 14 minutes.
The cron job can be found in the `.github/workflows/sleep-preventer.yml` file.
