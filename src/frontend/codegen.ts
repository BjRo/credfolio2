import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  schema: "../backend/internal/graphql/schema/schema.graphqls",
  documents: ["src/**/*.tsx", "src/**/*.ts", "src/graphql/**/*.graphql"],
  generates: {
    "./src/graphql/generated/": {
      preset: "client",
      presetConfig: {
        gqlTagName: "graphql",
      },
      config: {
        scalars: {
          DateTime: "string",
          JSON: "Record<string, unknown>",
        },
      },
    },
  },
  ignoreNoDocuments: true,
};

export default config;
