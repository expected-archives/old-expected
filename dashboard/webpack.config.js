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
        filename: '[hash].js',
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
                    },
                },
            },
            {
                test: /\.less$/,
                use: ['style-loader', 'css-loader', {
                    loader: 'less-loader',
                    options: {
                        javascriptEnabled: true,
                    },
                }],
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
    },
}
