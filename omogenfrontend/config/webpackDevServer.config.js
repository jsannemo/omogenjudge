'use strict';

const ignoredFiles = require('react-dev-utils/ignoredFiles');
const paths = require('./paths');

module.exports = {
    // Enable gzip compression of generated files.
    compress: true,
    // Reduce logging noise
    clientLogLevel: 'none',
    quiet: true,
    // Enable hot reloading server.
    hot: true,
    sockPort: 3000,
    // It is important to tell WebpackDevServer to use the same "publicPath" path as
    // we specified in the webpack config. When homepage is '.', default to serving
    // from the root.
    // remove last slash so user can land on `/test` instead of `/test/`
    publicPath: paths.publicUrlOrPath.slice(0, -1),
    // Reportedly, this avoids CPU overload on some systems.
    // https://github.com/facebook/create-react-app/issues/293
    // src/node_modules is not ignored to support absolute imports
    // https://github.com/facebook/create-react-app/issues/1065
    watchOptions: {
        ignored: ignoredFiles(paths.appSrc),
    },
    host: '127.0.0.1',
    // The HTML pages are served on a different port, so CORS must be allowed
    headers: {
        "Access-Control-Allow-Origin": "*"
    },
    // Enables client overlay on compile errors
    overlay: true,
};