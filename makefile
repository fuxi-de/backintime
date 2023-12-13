.PHONY: css
css:
	npx	tailwindcss -i ./src/templates/css/style.css -o dist/output.css --minify

## css-watch: watch build tailwindcss
.PHONY: css-watch
css-watch:
	npx tailwindcss -i src/templates/css/style.css -o dist/output.css --watch
