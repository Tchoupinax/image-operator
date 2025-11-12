import { type ImageBuilder, type Image } from "../../sdk/backend.generated";
import { BackendWrapper } from "../../sdk/wrapper";

type Data = {
  images: Array<Image>,
  imageBuilders: Array<ImageBuilder>,
}

export default defineEventHandler(async (): Promise<Data> => {
  const config = useRuntimeConfig();

  const backend = new BackendWrapper(config.public.graphqlApiUrl);
  try {

    const { images, imageBuilders } = await backend.sdk.fullData();
    return { images, imageBuilders } as Data;
  } catch (_) {
    return { images: [], imageBuilders: []}
  }
});
