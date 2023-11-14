const createProxyMiddleware  = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: process.env.GORYA_API_ADDR || 'http://localhost:8080',
      changeOrigin: true,
    })
  );
};
