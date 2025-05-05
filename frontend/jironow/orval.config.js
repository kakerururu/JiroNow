module.exports = {
  jironow: {
    input: '../docs/openapi.yaml',  // frontend/web から見たパス
    output: {
      mode: 'split',
      target: './jironow/src/api',
      client: 'axios',
      mock: true,
    },
  },
};
