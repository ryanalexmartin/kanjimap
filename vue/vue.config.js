const { defineConfig } = require('@vue/cli-service')
const webpack = require('webpack');
const dotenv = require('dotenv');
const path = require('path');

// Determine which .env file to use
const envPath = path.resolve(__dirname, '../.env');

// Load environment variables
const env = dotenv.config({ path: envPath }).parsed || {};

module.exports = defineConfig({
  transpileDependencies: true,
  publicPath: '/',
})

