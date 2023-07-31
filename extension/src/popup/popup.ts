import { MethodKind } from "@bufbuild/protobuf";
import { createPromiseClient } from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";

import {
  CreateGolinkRequest,
  CreateGolinkResponse,
} from "../../gen/golink/v1/golink_pb";
import { GolinkService } from "../../gen/golink/v1/golink_connect";

const apiBaseUrl = "http://localhost:8081/api/";

class GolinkPopup {
  // TODO: Don't use any
  constructor(client: any) {
    self.client = client;
  }

  async onSave() {
    const res = await client.createGolink({
      name: "mylink",
      url: "https://example.com",
    });
    console.log(res);
  }
}

async function initialize() {
  const transport = createConnectTransport({
    baseUrl: apiBaseUrl,
    credentials: "include",
  });
  const client = createPromiseClient(GolinkService, transport);

  const popup = new GolinkPopup(client);

  document.getElementById("save").addEventListener("click", popup.onSave);
  document.getElementById("cancel").addEventListener("click", popup.onCancel);
  console.log("initialized");
}

initialize();
