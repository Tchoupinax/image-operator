import gql from 'graphql-tag';
import { GraphQLClient, RequestOptions } from 'graphql-request';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
type GraphQLClientRequestHeaders = RequestOptions['requestHeaders'];
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type Image = {
  __typename?: 'Image';
  createdAt?: Maybe<Scalars['String']['output']>;
  destination?: Maybe<Source>;
  frequency?: Maybe<Scalars['String']['output']>;
  mode?: Maybe<Scalars['String']['output']>;
  name?: Maybe<Scalars['String']['output']>;
  source?: Maybe<Source>;
  status?: Maybe<Scalars['String']['output']>;
};

export type ImageBuilder = {
  __typename?: 'ImageBuilder';
  architecture?: Maybe<Scalars['String']['output']>;
  createdAt?: Maybe<Scalars['String']['output']>;
  name?: Maybe<Scalars['String']['output']>;
  source?: Maybe<Source>;
};

export type Mode =
  | 'OnceByTag'
  | 'OneShot'
  | 'Recurrent';

export type RootMutation = {
  __typename?: 'RootMutation';
  createImage?: Maybe<Image>;
};


export type RootMutationCreateImageArgs = {
  destinationRepositoryName: Scalars['String']['input'];
  destinationVersion: Scalars['String']['input'];
  mode: Mode;
  name: Scalars['String']['input'];
  sourceRepositoryName: Scalars['String']['input'];
  sourceVersion: Scalars['String']['input'];
};

export type RootQuery = {
  __typename?: 'RootQuery';
  imageBuilders?: Maybe<Array<Maybe<ImageBuilder>>>;
  images?: Maybe<Array<Maybe<Image>>>;
};

export type Source = {
  __typename?: 'Source';
  name?: Maybe<Scalars['String']['output']>;
  useAwsIRSA?: Maybe<Scalars['Boolean']['output']>;
  version?: Maybe<Scalars['String']['output']>;
};


export const FullData = gql`
    query fullData {
  images {
    name
    createdAt
    status
    source {
      name
      version
    }
    destination {
      name
      version
    }
  }
  imageBuilders {
    name
    createdAt
    architecture
  }
}
    `;
export const CreateImage = gql`
    mutation createImage($destinationRepositoryName: String!, $destinationVersion: String!, $sourceRepositoryName: String!, $sourceVersion: String!, $mode: Mode!, $name: String!) {
  createImage(
    destinationRepositoryName: $destinationRepositoryName
    destinationVersion: $destinationVersion
    mode: $mode
    name: $name
    sourceRepositoryName: $sourceRepositoryName
    sourceVersion: $sourceVersion
  ) {
    name
    createdAt
  }
}
    `;

export const FullDataDocument = gql`
    query fullData {
  images {
    name
    createdAt
    status
    source {
      name
      version
    }
    destination {
      name
      version
    }
  }
  imageBuilders {
    name
    createdAt
    architecture
  }
}
    `;
export const CreateImageDocument = gql`
    mutation createImage($destinationRepositoryName: String!, $destinationVersion: String!, $sourceRepositoryName: String!, $sourceVersion: String!, $mode: Mode!, $name: String!) {
  createImage(
    destinationRepositoryName: $destinationRepositoryName
    destinationVersion: $destinationVersion
    mode: $mode
    name: $name
    sourceRepositoryName: $sourceRepositoryName
    sourceVersion: $sourceVersion
  ) {
    name
    createdAt
  }
}
    `;

export type SdkFunctionWrapper = <T>(action: (requestHeaders?:Record<string, string>) => Promise<T>, operationName: string, operationType?: string, variables?: any) => Promise<T>;


const defaultWrapper: SdkFunctionWrapper = (action, _operationName, _operationType, _variables) => action();

export function getSdk(client: GraphQLClient, withWrapper: SdkFunctionWrapper = defaultWrapper) {
  return {
    fullData(variables?: FullDataQueryVariables, requestHeaders?: GraphQLClientRequestHeaders): Promise<FullDataQuery> {
      return withWrapper((wrappedRequestHeaders) => client.request<FullDataQuery>(FullDataDocument, variables, {...requestHeaders, ...wrappedRequestHeaders}), 'fullData', 'query', variables);
    },
    createImage(variables: CreateImageMutationVariables, requestHeaders?: GraphQLClientRequestHeaders): Promise<CreateImageMutation> {
      return withWrapper((wrappedRequestHeaders) => client.request<CreateImageMutation>(CreateImageDocument, variables, {...requestHeaders, ...wrappedRequestHeaders}), 'createImage', 'mutation', variables);
    }
  };
}
export type Sdk = ReturnType<typeof getSdk>;
export type FullDataQueryVariables = Exact<{ [key: string]: never; }>;


export type FullDataQuery = { __typename?: 'RootQuery', images?: Array<{ __typename?: 'Image', name?: string | null, createdAt?: string | null, status?: string | null, source?: { __typename?: 'Source', name?: string | null, version?: string | null } | null, destination?: { __typename?: 'Source', name?: string | null, version?: string | null } | null } | null> | null, imageBuilders?: Array<{ __typename?: 'ImageBuilder', name?: string | null, createdAt?: string | null, architecture?: string | null } | null> | null };

export type CreateImageMutationVariables = Exact<{
  destinationRepositoryName: Scalars['String']['input'];
  destinationVersion: Scalars['String']['input'];
  sourceRepositoryName: Scalars['String']['input'];
  sourceVersion: Scalars['String']['input'];
  mode: Mode;
  name: Scalars['String']['input'];
}>;


export type CreateImageMutation = { __typename?: 'RootMutation', createImage?: { __typename?: 'Image', name?: string | null, createdAt?: string | null } | null };
