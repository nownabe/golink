import { createPromiseClient, ConnectError } from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";
import type { PromiseClient } from "@bufbuild/connect";

import { GolinkService } from "../../gen/golink/v1/golink_connect";

const golinkUrlKey = "golinkUrl";

class GolinkPopup {
  client: PromiseClient<typeof GolinkService> | undefined;
  url: string | null;
  api: string | null = null;

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
    const notConfiguredElem = document.getElementById("not-configured")!;
    const unauthenticatedElem = document.getElementById("unauthenticated")!;
    const golinkUiElem = document.getElementById("golink-ui")!;

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
        document.getElementById("tab-title")!.innerText = title;
        (document.getElementById("url") as HTMLInputElement).value =
          tabs[0].url || "";
        const match = title.match(/go\/(?<name>[\w-]+)/);
        (document.getElementById("name") as HTMLInputElement).value =
          match?.groups?.name || "";
      }
      showGolinkUi();
      await this.listGolinks();
    } else {
      showUnauthenticated();
    }
  }

  public onSave = async () => {
    if (!this.client) {
      return;
    }

    const name = (document.getElementById("name") as HTMLInputElement).value;
    const url = (document.getElementById("url") as HTMLInputElement).value;

    if (!name || !url) {
      alert("golink name and url are required");
      return;
    }

    if (!this.validateName(name)) {
      alert("invalid golink name");
      return;
    }

    if (!this.validateUrl(url)) {
      alert("invalid URL");
      return;
    }

    try {
      await this.client.createGolink({
        name: name,
        url: url,
      });
      document.getElementById("save")!.hidden = true;
      document.getElementById("saved")!.hidden = false;
    } catch (rawErr) {
      const err = ConnectError.from(rawErr);
      alert(err.message);
    }
  };

  public openConsole = async () => {
    chrome.tabs.create({ url: this.url + "-/" });
  };

  private async checkAuth(): Promise<boolean> {
    try {
      console.log(`[checkAuth] trying to fetch ${this.url}health`);
      const res = await fetch(this.url + "health", {
        credentials: "include",
      });
      return res.status === 200;
    } catch (e) {
      console.error(e);
      return false;
    }
  }

  private buildClient() {
    if (!this.url || !this.api) return;

    const transport = createConnectTransport({
      baseUrl: this.api,
      credentials: "include",
    });
    this.client = createPromiseClient(GolinkService, transport);
  }

  private async listGolinks() {
    if (!this.client) {
      return;
    }

    const url = (document.getElementById("url") as HTMLInputElement).value;
    if (!url) return;

    try {
      const res = await this.client.listGolinksByUrl({ url });
      const golinks = res.golinks;

      if (golinks.length === 0) return;

      const ul = document.getElementById("golinks");
      if (!ul) return;

      golinks.forEach((golink) => {
        const consoleUrl = `${this.url}-/${golink.name}`;
        const li = document.createElement("li");
        const a = document.createElement("a");
        a.href = consoleUrl;
        a.innerText = "go/" + golink.name;
        a.onclick = () => chrome.tabs.create({ url: consoleUrl });
        li.appendChild(a);
        ul.appendChild(li);
      });

      document.getElementById("golinks-container")!.hidden = false;
    } catch (rawErr) {
      const err = ConnectError.from(rawErr);
      alert(err.message);
    }
  }

  private validateName(name: string): boolean {
    return /^[a-z0-9_-]+$/.test(name);
  }

  private validateUrl(url: string): boolean {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  }
}

async function initialize() {
  const popup = await GolinkPopup.create();
  await popup.initialize();

  document
    .getElementById("open-options")
    ?.addEventListener("click", async () => {
      await chrome.runtime.openOptionsPage();
    });
  document
    .getElementById("open-console")
    ?.addEventListener("click", popup.openConsole);
  document.getElementById("save")?.addEventListener("click", popup.onSave);
  console.log("initialized");
}

initialize();
