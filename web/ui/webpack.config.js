const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

module.exports = (_, { mode = 'development' }) => {
  return {
    mode,
    entry: {
      index: './index.js',
    },
    output: {
      path: path.join(__dirname, 'www'),
    },
    devtool: 'source-map',
    module: {
      rules: [
        {
          test: /\.p?css$/,
          use: [
            MiniCssExtractPlugin.loader,
            { loader: 'css-loader', options: { importLoaders: 1 } },
            'postcss-loader',
          ],
        },
        {
          test: /\.(svg|png|ico|jpe?g|gif)(\?.*)?$/i,
          use: {
            loader: 'url-loader',
            options: {
              limit: 1,
            },
          },
        },
        {
          test: /\.xlsx(\?.*)?$/i,
          use: {
            loader: 'url-loader',
            options: {
              limit: 1,
            },
          },
        },
        {
          test: /\.(woff2?|eot|ttf|otf)(\?.*)?$/i,
          use: {
            loader: 'url-loader',
            options: {
              limit: 1,
            },
          },
        },
      ],
    },
    plugins: [
      new HtmlWebpackPlugin({
        publicPath: '/',
        template: './index.html',
        favicon: './assets/favicon.ico',
      }),
      new MiniCssExtractPlugin(),
    ],
  };
};
