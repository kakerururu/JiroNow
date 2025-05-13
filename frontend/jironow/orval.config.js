export default {
  jironow: {
    input: '../../backend/docs/openapi.yaml',
    output: {
      mode: 'split',
      target: './src/api',
      client: 'axios',
      mock: true,
    },
  },
};
