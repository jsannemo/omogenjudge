const webpack = require('webpack');
const {merge} = require('webpack-merge');

const common = require('./webpack.config.js');

module.exports = merge(common, {
  mode: 'development',
  devtool: 'inline-source-map',
  devServer: {
    hot: true,
    overlay: true,
    port: 3000,
    sockPort: 3000,
    publicPath: '/static/js/',
    headers: {"Access-Control-Allow-Origin": "*"},
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
  ]
});
