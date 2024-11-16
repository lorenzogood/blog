-include .env

# Enable app debug logs.
DEBUG ?= 1
WEB_SRC ?= ./web_src
NEWSLETTER_WEB_ASSET_DIR=./web_src/dist
NEWSLETTER_DEVELOPMENT=1
NEWSLETTER_WEB_TEMPLATE_DIR=./templates

export

.PHONY: dev
dev: web-deps
	@./tools/dev

.PHONY: web-deps
web-deps: $(WEB_SRC)/node_modules

$(WEB_SRC)/node_modules: $(WEB_SRC)/package-lock.json
	cd $(WEB_SRC) && npm ci

$(WEB_SRC)/dist: $(WEB_SRC)/node_modules $(WEB_SRC)/package.json $(WEB_SRC)/src $(WEB_SRC)/webpack.config.js 
	cd $(WEB_SRC) && node node_modules/webpack-cli/bin/cli.js --mode production 
