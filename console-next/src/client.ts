import { createPromiseClient } from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";

import { GolinkService } from "./gen/golink/v1/golink_connect";

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "/api";

const transport = createConnectTransport({
  baseUrl: apiUrl,
  credentials: "include",
});
const client = createPromiseClient(GolinkService, transport);

export default client;
