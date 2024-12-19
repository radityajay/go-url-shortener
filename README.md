# Shortened URL API

A simple API for shortening URLs built with Golang.

## Features

* Shorten long URLs to unique short URLs.
* Automatic redirection from short URLs to the original URL.
* API endpoints to create and retrieve short URLs.
* Custom aliases for short URLs.
* Integration with storage services (PostgreSQL) for data persistence.

## How to Run

1. Clone this repository.
2. Install dependencies: `go mod download`
3. Configure database and environment variables.
4. Run the application: `go run main.go`

## API Endpoints

### Create Short URLs

**POST {{BASE_URL}}/shortener**

**Request Body:**

```json
{
 "url": "https://www.exampleurlpanjang.com/example-url-panjang",
 [Optional] "custom_alias": "short-url"
}
```

**Response**

```json
{
 "success": true,
 "url": "http://localhost:4000/s/short-url",
 "message": "Success"
}
```

### Short URL Redirect
**GET {{BASE_URL}}/s/:short_code**

**Example:**

`GET /s/short-url` will redirect to `https://www.exampleurllenght.com/example-url-lenght`
