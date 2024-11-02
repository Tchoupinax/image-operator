import { GraphQLClient } from "graphql-request";

import { getSdk, type Sdk } from "./backend.generated";

export class BackendWrapper {
  private client: GraphQLClient;
  public sdk: Sdk;

  constructor(backendUrl: string) {
    this.client = new GraphQLClient(backendUrl);
    this.sdk = getSdk(this.client);
  }
}
