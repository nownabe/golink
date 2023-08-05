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
  client: PromiseClient<typeof GolinkService> | undefined;
  url: string | null;

  constructor(url: string | null) {
    this.url = url;
    if (this.url) {
      if (!this.url.endsWith("/")) {
        this.url += "/";
      }
      this.api = this.url + "api";
      this.buildClient();
    }
  }

  public static async create(): Promise<GolinkPopup> {
    const url = (await chrome.storage.sync.get(golinkUrlKey))[golinkUrlKey];
    return new GolinkPopup(url);
  }

  public async initialize() {
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

    if (!this.url) {
      showNotConfigured();
      return;
    }

    if (await this.checkAuth()) {
      const tabs = await chrome.tabs.query({
        active: true,
        currentWindow: true,
      });
      if (tabs.length > 0) {
        const title = tabs[0].title || "";
        document.getElementById("tab-title").innerText = title;
        document.getElementById("url").value = tabs[0].url;
        const match = title.match(/go\/(?<name>[\w-]+)/);
        document.getElementById("name").value = match?.groups?.name || "";
      }
      showGolinkUi();
    } else {
      showUnauthenticated();
    }
  }

  private async checkAuth(): boolean {
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

  public onSave = async () => {
    const res = await this.client.createGolink({
      name: "mylink",
      url: "https://example.com",
    });
    console.log(res);
  };

  public openConsole = async () => {
    chrome.tabs.create({ url: this.url + "c" });
  };

  private buildClient() {
    if (!this.url) return;

    const transport = createConnectTransport({
      baseUrl: this.api,
      credentials: "include",
    });
    this.client = createPromiseClient(GolinkService, transport);
  }
}

async function initialize() {
  const popup = await GolinkPopup.create();
  await popup.initialize();

  document
    .getElementById("open-options")
    .addEventListener("click", async () => {
      await chrome.runtime.openOptionsPage();
    });
  document
    .getElementById("open-console")
    .addEventListener("click", popup.openConsole);
  document.getElementById("save").addEventListener("click", popup.onSave);
  console.log("initialized");
}

initialize();
