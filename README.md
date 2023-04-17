# Crop Connect Backend

## Setup

1. Install dependencies

```bash
go mod download
```

2. Create a `.env` file in the root directory and add the following environment variables:

```bash
cp .env.example .env
```

Note:

- APP_DOMAIN delimiter is a comma

3. Run the server

```bash
go run main.go
```
