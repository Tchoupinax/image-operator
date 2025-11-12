import { BackendWrapper } from "../../sdk/wrapper";
import { logger } from "../tools/logger";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const body = await readBody(event)

  const payload = {
    destinationRepositoryName: body.destinationRepository,
    destinationVersion: body.destinationVersion,
    mode: body.mode,
    sourceRepositoryName: body.sourceRepository,
    sourceVersion: body.sourceVersion,
    name: body.name,
  }

  logger.info(payload, 'Create image manually')

  const backend = new BackendWrapper(config.public.graphqlApiUrl);
  const image = await backend.sdk.createImage(payload);

  return "OK";
});
