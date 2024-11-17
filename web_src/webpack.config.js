const path = require('path');
const glob = require('glob');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');
const RemoveEmptyScriptsPlugin = require('webpack-remove-empty-scripts');

const getEntryPoints = () => {
	const entries = {};

	const scssFiles = glob.sync('./src/**/[^_]*.scss');
	scssFiles.forEach(file => {
		const name = path.basename(file, '.scss');
		entries[`css/${name}`] = "./" + file;
	});

	const jsFiles = glob.sync('./src/**/*[^_]*.js');
	jsFiles.forEach(file => {
		const name = path.basename(file, '.js');
		entries[`js/${name}`] = "./" + file;
	});

	return entries;
};

module.exports = {
	entry: getEntryPoints(),
	output: {
		path: path.resolve(__dirname, 'dist'),
		clean: true,
		filename: "[name]~[contenthash].js",
	},
	module: {
		rules: [
			{
				test: /\.scss$/i,
				use: [MiniCssExtractPlugin.loader, "css-loader", "sass-loader"],
			},
		],
	},
	plugins: [
		new RemoveEmptyScriptsPlugin(),
		new MiniCssExtractPlugin({
			filename: "[name]~[contenthash].css",
		})
	],
	optimization: {
		sideEffects: true,
		removeEmptyChunks: true,
		minimizer: [
			new CssMinimizerPlugin({
				minimizerOptions: {
					preset: [
						'default',
						{
							discardComments: { removeAll: true },
						},
					],
				},
			}),
			'...',
		],
	},
};
