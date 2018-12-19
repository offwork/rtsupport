const path = require('path');

module.exports = {
  mode: 'none',
  entry: './index.js',
  output: {
    filename: 'main.js',
    path: path.resolve(__dirname, '.')
  },
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loader: 'babel-loader'
      }
    ]
  },
  devServer: {
    contentBase: path.join(__dirname, '.'),
    index: 'index.html',
    filename: 'main.js',
    compress: true,
    port: 9200
  }
};
