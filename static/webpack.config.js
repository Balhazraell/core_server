module.exports = {
    mode: 'development',
    entry: './js/src/main.js',
    output: {
        path: __dirname + "/js/dist/",
        filename: 'main_build.js'
    },

    watch: true,
    devtool: 'source-map'
}