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
  client: PromiseClient<typeof GolinkService>;

  constructor(url: string) {
    this.setUrl(url);
  }

  setUrl(url: string) {
    this.url = url;
    this.api = url + "api";

    const transport = createConnectTransport({
      baseUrl: this.api,
      credentials: "include",
    });
    this.client = createPromiseClient(GolinkService, transport);
  }

  async checkAuth(): boolean {
    try {
      const res = await fetch(this.api + "/healthz", {
        credentials: "include",
      });
      return res.status === 200;
    } catch (e) {
      console.error(e);
      return false;
    }
  }

  onSave = async () => {
    const res = await this.client.createGolink({
      name: "mylink",
      url: "https://example.com",
    });
    console.log(res);
  };

  openConsole = async () => {
    chrome.tabs.create({ url: this.url + "c" });
  };
}

async function getGolinkUrl(): Promise<string | null> {
  let url = (await chrome.storage.sync.get(golinkUrlKey))[golinkUrlKey];
  if (!url) {
    return null;
  }

  if (!url.endsWith("/")) {
    url += "/";
  }

  return url;
}

function updateClientFunc(popup: GolinkPopup) {
  const notConfiguredElem = document.getElementById("not-configured");
  const unauthenticatedElem = document.getElementById("unauthenticated");
  const golinkUiElem = document.getElementById("golink-ui");

  const showNotConfigured = () => {
    notConfiguredElem.hidden = false;
    unauthenticatedElem.hidden = true;
    golinkUiElem.hidden = true;
  };

  const showUnauthenticated = () => {
    notConfiguredElem.hidden = true;
    unauthenticatedElem.hidden = false;
    golinkUiElem.hidden = true;
  };

  const showGolinkUi = () => {
    notConfiguredElem.hidden = true;
    unauthenticatedElem.hidden = true;
    golinkUiElem.hidden = false;
  };

  return async () => {
    const url = await getGolinkUrl();

    if (!url) {
      showNotConfigured();
      return;
    }

    popup.setUrl(url);

    if (await popup.checkAuth()) {
      showGolinkUi();
    } else {
      showUnauthenticated();
    }
  };
}

async function openOptionsPage() {
  await chrome.runtime.openOptionsPage();
}

async function initialize() {
  const popup = new GolinkPopup();

  const updateClient = updateClientFunc(popup);
  await updateClient();

  document
    .getElementById("open-options")
    .addEventListener("click", openOptionsPage);
  document
    .getElementById("open-console")
    .addEventListener("click", popup.openConsole);
  document.getElementById("save").addEventListener("click", popup.onSave);
  document.getElementById("cancel").addEventListener("click", popup.onCancel);
  console.log("initialized");
}

initialize();
