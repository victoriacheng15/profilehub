# ProfileHub Default Template

Simple, single‑page template for ProfileHub projects. This is what `init` generates by default.

## Quick start

1) Edit your info in `config/config.yml`.

2) Run a local server (from the project root):

```bash
profilehub dev      # if you installed the CLI
# or
go run main.go dev  # if you're running from the repo
```

3) Build static files to `dist/`:

```bash
profilehub build
# or
go run main.go build
```

## What’s inside

```
.
├── config/
│   └── config.yml
└── src/
		├── index.html
		└── static/
				└── css/
						└── styles.css
```

## Config reference (config/config.yml)

```yaml
Params:
	Avatar: "https://example.com/avatar.png"   # or static path
	Name: "Your Name"
	Headline: "Short headline or bio"

	Theme:
		Background: "#ffffff"
		Text: "#111111"
		Button: "#2563eb"
		ButtonText: "#ffffff"
		ButtonHover: "#1d4ed8"

	Socials:
		- Icon: "fab fa-twitter"  # or emoji, e.g. "🐦"
			URL: "https://twitter.com/yourhandle"
		- Icon: "fab fa-linkedin"
			URL: "https://linkedin.com/in/yourhandle"

	Links:
		- Name: "My Website"
			URL: "https://yourwebsite.com"
		- Name: "Portfolio"
			URL: "https://yourportfolio.com"
```

## Customization

- Edit styles in `src/static/css/styles.css`.
- Replace the avatar with your own image (URL or place a file under `src/static/`).
- Font Awesome is already included via CDN in `src/index.html`; use icon classes like `<i class="fab fa-github"></i>` in your config’s `Socials`.

## Deploy

The `build` command outputs a fully rendered site in `dist/`. Deploy that folder to any static host (GitHub Pages, Vercel, Netlify, etc.).

### Deploy to GitHub Pages (recommended)

1) In your repository, go to Settings → Pages and set Source to “GitHub Actions”.

2) Add a workflow at `.github/workflows/deploy.yml`:

```yaml
name: Deploy to GitHub Pages

on:
	push:
		branches: [main]

permissions:
	contents: read
	pages: write
	id-token: write

concurrency:
	group: "pages"
	cancel-in-progress: true

jobs:
	build:
		runs-on: ubuntu-latest
		steps:
			- uses: actions/checkout@v4
			- uses: actions/setup-go@v5
				with:
					go-version: '1.22'

			- name: Build static site
				run: go run main.go build

			- name: Add .nojekyll
				run: mkdir -p dist && touch dist/.nojekyll

			- name: Upload artifact
				uses: actions/upload-pages-artifact@v3
				with:
					path: dist

	deploy:
		environment:
			name: github-pages
			url: ${{ steps.deployment.outputs.page_url }}
		runs-on: ubuntu-latest
		needs: build
		steps:
			- name: Deploy to GitHub Pages
				id: deployment
				uses: actions/deploy-pages@v4
```

3) Push to `main`. The site will be available at:
- User site: `https://<username>.github.io/`
- Project site: `https://<username>.github.io/<repo>/`

Tips
- Use relative asset paths like `static/...` instead of `/static/...` so assets work on project pages under a sub-path.
- Custom domain: set it in repo Settings → Pages. If you commit a `CNAME` file, ensure it’s preserved in your deploy output.