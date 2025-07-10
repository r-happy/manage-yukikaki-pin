export const environment = {
  production: true,
  backendUrl:
    (typeof process !== 'undefined' ? process.env['BACKEND_URL'] : undefined) ||
    'http://localhost:1323',
};
