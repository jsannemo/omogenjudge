const path = require('path');
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
    contentBase: path.join(__dirname, './frontend/assets'),
    headers: {"Access-Control-Allow-Origin": "*"},
    proxy: [
      {
        context: ['/omogen.webapi', '/problems/img/'],
        target: 'http://localhost:56744',
      },
    ],
    historyApiFallback: {
      index: 'index.html'
    }
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
  ]
});
