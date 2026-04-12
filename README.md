# gorate

A lightweight CLI tool for applying and inspecting rate-limit policies on HTTP endpoints during local development.

---

## Installation

```bash
go install github.com/yourusername/gorate@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/gorate.git && cd gorate && go build -o gorate .
```

---

## Usage

Apply a rate-limit policy to a local endpoint and proxy requests through it:

```bash
# Limit to 10 requests per second on a local API
gorate run --target http://localhost:8080 --rate 10 --per second

# Inspect current policy and live request stats
gorate inspect --target http://localhost:8080

# Apply a burst-aware policy
gorate run --target http://localhost:3000 --rate 100 --per minute --burst 20
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--target` | Target HTTP endpoint URL | *(required)* |
| `--rate` | Number of allowed requests | `60` |
| `--per` | Time window: `second`, `minute`, `hour` | `minute` |
| `--burst` | Max burst size above the rate limit | `0` |
| `--port` | Local proxy port to listen on | `9000` |

Once running, send your requests to `http://localhost:9000` and gorate will enforce the policy before forwarding them to the target.

---

## Why gorate?

Testing how your application behaves under rate-limiting shouldn't require deploying to a staging environment or mocking an external gateway. gorate lets you simulate real rate-limit conditions locally in seconds.

---

## Contributing

Pull requests and issues are welcome. Please open an issue before submitting large changes.

---

## License

[MIT](LICENSE)