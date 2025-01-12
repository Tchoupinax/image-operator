export type RegistryImage = {
  downloadCount?: number;
  isOfficial: boolean;
  name: string;
  registry: "Amazon ECR" | "Quay.io" | "DockerHub";
}

export default defineEventHandler(async (event): Promise<Array<RegistryImage>> => {
  const { repo } = getQuery(event)
  const string = repo;

  const dockerHubUrl = `https://hub.docker.com/v2/search/repositories/?query=${string}`;
  const quayUrl = `https://quay.io/api/v1/find/repositories?query=${string}`
  const amazonUrl = `https://api.us-east-1.gallery.ecr.aws/searchRepositoryCatalogData`

  const [dockerHub, quay, amazon] = await Promise.all([
    fetch(dockerHubUrl).then(res => res.json()),
    fetch(quayUrl).then(res => res.json()),
    fetch(amazonUrl, {
      method: "POST",
      body: JSON.stringify({
        searchTerm: string,
        sortConfiguration: { sortKey: "POPULARITY" }
      }),
    }).then(res => res.json())
  ])

  const result = ([
    ...amazon.repositoryCatalogSearchResultList.map((r) => ({
      downloadCount: r.downloadCount,
      isOfficial: r.registryVerified,
      name: r.repositoryName,
      registry: "Amazon ECR",
    } satisfies RegistryImage)),
    ...quay.results.map((r: any) => ({
      name: `${r.namespace.name}/${r.name}`,
      registry: "Quay.io",
      isOfficial: false,
    } satisfies RegistryImage)),
    ...dockerHub.results.map((s: any) => ({
      downloadCount: s.pull_count,
      isOfficial: s.is_official,
      name: s.repo_name,
      registry: "DockerHub",
    } satisfies RegistryImage))
  ] as Array<RegistryImage>)

  return limitResultReturned(sortResult(result));
});

function sortResult(result: Array<RegistryImage>): Array<RegistryImage> {
  return result.sort((a, b) => {
    if (a.isOfficial && !b.isOfficial) {
      return -1;
    }
    if (!a.isOfficial && b.isOfficial) {
      return 1;
    }
    // Secondary criterion: Higher downloadCount comes first
    const downloadCountA = a.downloadCount || 0; // Treat missing as 0
    const downloadCountB = b.downloadCount || 0; // Treat missing as 0
    return downloadCountB - downloadCountA;
  });
}

function limitResultReturned(result: Array<RegistryImage>): Array<RegistryImage> {
  return result.slice(0, 1000);
}