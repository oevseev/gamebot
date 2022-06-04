const path = require('path');
const { env } = require('process');

module.exports = (env) => {
  return {
    mode: 'development',
    entry: './src/index.tsx',
    devtool: 'inline-source-map',
    module: {
      rules: [
        {
          test: /\.tsx?$/,
          use: 'ts-loader',
          exclude: /node_modules/,
        },
        {
          test: /\.css$/i,
          include: path.resolve(__dirname, 'src'),
          use: ['style-loader', 'css-loader', 'postcss-loader'],
        },
      ],
    },
    resolve: {
      extensions: ['.tsx', '.ts', '.js'],
    },
    output: {
      filename: 'static/bundle.js',
      path: path.resolve(__dirname, 'dist'),
    },
    devServer: {
      compress: true,
      devMiddleware: {
        index: false,
      },
      static: {
        directory: path.join(__dirname, 'public'),
        publicPath: '/static',
      },
      proxy: [
        {
          context: ['**', '!/static/**'],
          target: env.webserver_endpoint,
          secure: false,
        },
      ],
    },
  };
};
