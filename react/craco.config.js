const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const path = require('path');
path.resolve(__dirname, 'src/index.tsx');

function webpackConfig(config) {
  config.plugins.push(new MiniCssExtractPlugin({ filename: 'styles.css' }));
  config.module = {
    rules: [
      {
        test: /\.(ts|tsx)?$/,
        use: [
          { loader: 'babel-loader' },
          {
            loader: '@linaria/webpack-loader',
            options: {
              cacheDirectory: 'src/.linaria_cache',
              sourceMap: process.env.NODE_ENV !== 'production'
            }
          },
          {
            loader: 'esbuild-loader',
            options: {
              loader: 'tsx',
              target: 'esnext'
            }
          }
        ]
      },
      {
        test: /\.css$/,
        use: [
          {
            loader: MiniCssExtractPlugin.loader
          },
          {
            loader: 'css-loader',
            options: {
              sourceMap: process.env.NODE_ENV !== 'production'
            }
          }
        ]
      }
    ]
  };
  return config;
}
module.exports = {
  webpack: { configure: webpackConfig }
};
