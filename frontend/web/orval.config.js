module.exports = {
  jironow: {
    input: '../docs/openapi.yaml',  // frontend/web から見たパス
    output: {
      mode: 'split',
      target: './src/api/',         // web/src/api に生成
      client: 'axios',
      mock: true,
    },
  },
};
