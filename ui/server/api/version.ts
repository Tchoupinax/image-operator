export default defineEventHandler(() => {
  const config = useRuntimeConfig();
  return config.public.operatorVersion;
});
