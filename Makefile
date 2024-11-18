-include .env

# Enable app debug logs.
DEBUG ?= 1
WEB_SRC ?= ./web_src
BLOG_WEB_ASSET_DIR=./web_src/dist
BLOG_DEVELOPMENT=1
BLOG_WEB_TEMPLATE_DIR=./templates
BLOG_CONTENT_DIR=./content

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
