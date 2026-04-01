# Stockyard Branding Iron

**OG image generation API.** Generate social sharing images via URL params or API. Templates with customizable colors. Like Vercel OG but self-hosted. Single binary, no external dependencies.

## Usage

```html
<!-- In your HTML head -->
<meta property="og:image" content="https://your-server:9040/api/og?title=My+Post&subtitle=Read+more" />
```

```bash
# Generate via URL
curl "http://localhost:9040/api/og?title=Hello+World&subtitle=My+Blog" > og.svg

# Generate via POST
curl -X POST http://localhost:9040/api/og \
  -H "Content-Type: application/json" \
  -d '{"title":"Hello World","subtitle":"My Blog"}'

# Custom colors
curl "http://localhost:9040/api/og?title=Dark+Mode&bg=000000&fg=ffffff&accent=ff6600"
```

## Free vs Pro

| Feature | Free | Pro ($4.99/mo) |
|---------|------|----------------|
| Templates | 3 | Unlimited |
| Custom colors | ✓ | ✓ |
| Remove watermark | — | ✓ |
| Custom fonts | — | ✓ |

## License

Apache 2.0
