# Covid data API (in GO)

---

ENDPOINTS:

GET `/data` to fetch data as JSON
GET `/update` to update data using firestore

---

start local server with VS Code debugger - press F5

start local server via command (pass in env var as cred)
```bash
GOOGLE_APPLICATION_CREDENTIALS=config/sd-covid-2-3c873e023505.json go run covid-api.go
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
