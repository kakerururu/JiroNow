module.exports = {
  jironow: {
    input: '../../backend/docs/openapi.yaml',
    output: {
      mode: 'split',
      target: './jironow/src/api',
      client: 'axios',
      mock: true,
    },
  },
};
