const { join } = require('path')
const { HotModuleReplacementPlugin, DefinePlugin } = require('webpack')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')

module.exports = {
    mode: 'development',
    entry: join(__dirname, 'src', 'index.jsx'),
    output: {
        path: join(__dirname, 'build'),
        publicPath: '/',
        filename: '[hash].js',
    },
    resolve: {
        extensions: ['.js', '.jsx'],
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env', '@babel/preset-react'],
                        plugins: ['@babel/plugin-proposal-class-properties'],
                    },
                },
            },
            {
                test: /\.scss$/,
                use: ['style-loader', 'css-loader', 'sass-loader'],
            },
            {
                test: /\.css$/,
                use: ['style-loader', 'css-loader'],
            },
            {
                test: /\.(png|jpg|gif)$/,
                use: 'file-loader',
            },
        ],
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: '[hash].css',
            allChunks: true,
        }),
        new DefinePlugin({
            'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'development'),
        }),
        new HtmlWebpackPlugin({
            template: join(__dirname, 'public', 'index.html'),
        }),
        new HotModuleReplacementPlugin(),
    ],
    optimization: {
        splitChunks: {
            name: 'vendor',
            filename: '[hash].js',
        },
        minimizer: [
            new UglifyJsPlugin({
                sourceMap: true,
            }),
        ],
    },
    devServer: {
        port: process.env.PORT || 8080,
        disableHostCheck: true,
        historyApiFallback: {
            filename: 'index.html',
        },
    },
}
