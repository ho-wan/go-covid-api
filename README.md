# Covid data API (in GO)

`go-covid-api` fetches data from GCP firestore as json and serves it through the `/data` endpoint. This is used for the `covid-chart-app` which fetches the data and plots cases against date by country.

---

ENDPOINTS:

GET `/data` to fetch data as JSON
GET `/update` to update data using firestore

---

start local server with VS Code debugger - press F5

start local server via command line (pass in env var as cred)
```bash
GOOGLE_APPLICATION_CREDENTIALS=config/{filename.json} go run covid-api.go
```

---

automatic deploy

- `cloudbuild.yaml` configured to deploy automatically when pushing to master branch on github

---

manual build & deploy

```bash
gcloud builds submit --tag gcr.io/sd-covid-2/covid-api
```
```bash
gcloud run deploy --image gcr.io/sd-covid-2/covid-api --platform managed covid-api
```

---

## Overview

Steps:

`py-covid-data-fetch`
- Python service running on GCP Cloud Run
- Fetch CSV data from [JHU github](https://github.com/CSSEGISandData/COVID-19)
- Parses data: last 28 days, sum by country, omit lat/long coords
- Loads data into json store (GCP firestore)
- Runs once a day through GCP Cloud Scheduler calling endpoint, and calls `go-covid-api` to update

`go-covid-api`
- Golang service running on GCP Cloud Run
- Fetch json data from GCP firestore, and loaded into state through `/update`
- Data served through `/data`


`covid-chart-app`
- React app running on AWS 
- Fetches data from `go-covid-api`
- Plots data using Nivo-charts
