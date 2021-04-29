module.exports = ({ env }) => {
  const plugins = {
    'postcss-import': {},
    'postcss-url': {},
    'postcss-preset-env': {
      features: {
        'nesting-rules': true,
      },
    },
  };

  if (env === 'production') {
    Object.assign(plugins, {
      '@fullhuman/postcss-purgecss': {
        content: ['./src/**/*.js'],
      },
      cssnano: {},
    });
  }

  return {
    plugins,
  };
};
