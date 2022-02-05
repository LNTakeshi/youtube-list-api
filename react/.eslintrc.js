module.exports = {
  env: {
    // 'browser': true,
    // 'es6': true,
    es2021: true
  },
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaFeatures: {
      jsx: true
    },
    ecmaVersion: 'latest',
    sourceType: 'module'
  },
  plugins: ['react', '@typescript-eslint'],
  rules: {
    quotes: ['error', 'single']
  },
  extends: ['plugin:prettier/recommended', 'plugin:react-hooks/recommended']
};
