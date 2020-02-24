# Covid data API (in GO)

run locally with VS Code debugger - press F5

run locally via command (pass in env var as cred)
```bash
GOOGLE_APPLICATION_CREDENTIALS=config/sd-covid-2-3c873e023505.json go run covid-api.go
```

```bash
gcloud builds submit --tag gcr.io/sd-covid-2/covid-api
```
```bash
gcloud run deploy --image gcr.io/sd-covid-2/covid-api --platform managed
```
