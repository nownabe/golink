import { createPromiseClient } from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";

import { GolinkService } from "./gen/golink/v1/golink_connect";
import { apiEndpoint } from "./config";

const transport = createConnectTransport({
  baseUrl: apiEndpoint,
  credentials: "include",
});
const client = createPromiseClient(GolinkService, transport);
export default client;
