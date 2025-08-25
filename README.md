# purr

![purr](https://img.shields.io/badge/privacy-focused-brightgreen)

A **privacy-focused web analytics tool**. Collects page view data without cookies or any personally-identifiable information, so no GDPR notice is required.

---

purr uses a **1x1 transparent pixel** to capture analytics.

1. **JavaScript-enabled tracking**: dynamically creates an `<img>` element with page info encoded in the URL.
2. **No-JS fallback**: uses a static `<img>` and the HTTP Referer to infer page details.

Both methods ensure **no cookies, no user identifiers**, and minimal data storage.

---

## Installation

Add this snippet to any page you want to track. Replace `https://example.com` with your server URL:

```html
<script>
(function () {
  try {
    var u = new URL("https://example.com/pixel.gif");
    u.searchParams.set("js", "1");
    u.searchParams.set("p", location.pathname);
    u.searchParams.set("t", document.title.slice(0, 120));
    u.searchParams.set("w", String(window.innerWidth || 0));
    u.searchParams.set("h", String(window.innerHeight || 0));
    u.searchParams.set("dpr", String(window.devicePixelRatio || 1));
    u.searchParams.set("lang", (navigator.language || "").slice(0, 16));
    u.searchParams.set("ref", document.referrer || "");

    var img = new Image(1, 1);
    img.decoding = "async";
    img.referrerPolicy = "strict-origin-when-cross-origin";
    img.src = u.toString();
  } catch (e) { /* ignore */ }
})();
</script>

<noscript>
  <img src="https://example.com/pixel.gif" width="1" height="1" alt="" />
</noscript>
```

## Building

### Prerequisites

- [Go 1.24](https://go.dev/) for the backend.
- Optional: Docker for containerized builds.

```bash
# Clone the repo
git clone https://github.com/yourusername/purr.git
cd purr

# Run the server
go run .
```

The server listens on port `8080` by default. Access it via:

- `/pixel.gif` - tracking pixel
- `/stats` - aggregated stats

To deploy in Docker:

```bash
# Build the Docker image
docker build -t purr:latest .

# Run the container
docker run -p 8080:8080 purr:latest
```
