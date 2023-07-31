import { MethodKind } from "@bufbuild/protobuf";
import { createPromiseClient } from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";

import type { PromiseClient } from "@bufbuild/connect";

import {
  CreateGolinkRequest,
  CreateGolinkResponse,
} from "../../gen/golink/v1/golink_pb";
import { GolinkService } from "../../gen/golink/v1/golink_connect";

const golinkUrlKey = "golinkUrl";

class GolinkPopup {
  client: PromiseClient<GolinkService>;

  constructor(client: PromiseClient<GolinkService>) {
    this.client = client;
  }

  setClient(client: PromiseClient<GolinkService>) {
    this.client = client;
  }

  onSave = async () => {
    console.log(this.client);
    const res = await this.client.createGolink({
      name: "mylink",
      url: "https://example.com",
    });
    console.log(res);
  };
}

async function buildClient() {
  const url = (await chrome.storage.sync.get(golinkUrlKey))[golinkUrlKey];
  const transport = createConnectTransport({
    baseUrl: url + "/api",
    credentials: "include",
  });
  return createPromiseClient(GolinkService, transport);
}

async function initialize() {
  const popup = new GolinkPopup(await buildClient());

  chrome.storage.onChanged.addListener(
    async (
      changes: { [key: string]: chrome.storage.StorageChange },
      namespace: string,
    ) => {
      if (namespace === "sync" && golinkUrlKey in changes) {
        popup.setClient(await this.buildClient());
      }
    },
  );

  document.getElementById("save").addEventListener("click", popup.onSave);
  document.getElementById("cancel").addEventListener("click", popup.onCancel);
  console.log("initialized");
}

initialize();
