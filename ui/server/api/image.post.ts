import { BackendWrapper } from "../../sdk/wrapper";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const body = await readBody(event)

  const backend = new BackendWrapper(config.public.graphqlApiUrl);
  await backend.sdk.createImage({
    destinationRepositoryName: body.destinationRepository,
    destinationVersion: body.destinationVersion,
    mode: body.mode,
    sourceRepositoryName: body.sourceRepository,
    sourceVersion: body.sourceVersion,
    name: body.name,
  });

  return "OK";
});
