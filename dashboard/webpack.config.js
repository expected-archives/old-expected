const { join } = require('path')
const { HotModuleReplacementPlugin, DefinePlugin } = require('webpack')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')

module.exports = {
    mode: 'development',
    entry: join(__dirname, 'src', 'index.tsx'),
    output: {
        path: join(__dirname, 'build'),
        publicPath: '/',
        filename: '[hash].js',
    },
    resolve: {
        extensions: ['.js', '.ts', '.tsx'],
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                loader: 'awesome-typescript-loader',
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
