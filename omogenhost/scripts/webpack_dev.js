const Webpack = require("webpack");
const WebpackDevServer = require('webpack-dev-server');
const webpackConfig = require("../webpack.dev");

process.on('unhandledRejection', err => {
  throw err;
});

webpackConfig.mode = 'development';
const compiler = Webpack(webpackConfig);
const devServerOptions = {...webpackConfig.devServer};
const devServer = new WebpackDevServer(compiler, devServerOptions);
devServer.listen(3000, "127.0.0.1", err => {
  if (err) {
    return console.log(err);
  }
});
