module.exports = {
    mode: 'development',
    entry: './src/main.js',
    output: {
        path: __dirname + "/src/dist/",
        filename: 'main_build.js'
    },

    watch: true,
    devtool: 'source-map'
}