### Building and running your application

Build the image from the `Dockerfile`:

```bash
docker build -t url-shortener .
```

Run the container, passing environment variables from your local `.env` file and mapping port `8080`:

```bash
docker run --rm --env-file .env -p 8080:8080 url-shortener
```

Your application will be available at http://localhost:8080.

### References
* [Docker's Go guide](https://docs.docker.com/language/golang/)
