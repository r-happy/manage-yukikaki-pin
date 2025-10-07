type RuntimeEnv = typeof globalThis & {
  BACKEND_URL?: string;
  __env?: Record<string, string | undefined>;
  process?: {
    env?: Record<string, string | undefined>;
  };
};

const runtimeEnv = globalThis as RuntimeEnv;

const backendUrl =
  runtimeEnv.process?.env?.['BACKEND_URL'] ??
  runtimeEnv.__env?.['BACKEND_URL'] ??
  runtimeEnv.BACKEND_URL ??
  'http://localhost:1323';

export const environment = {
  production: false,
  backendUrl,
};
