all:
	npx tailwindcss -i ./input.css -o ./public/style.css
	templ generate

.PHONY: all
