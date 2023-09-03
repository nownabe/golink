import {
  fetchGolinkUrl,
  fetchGolinkUrlName,
  isManaged,
  isManagedName,
  saveGolinkUrl,
  saveGolinkUrlName,
} from "./messageListeners";
import { Router } from "./router";
import { updateRedirectRule } from "./updateRedirectRule";

async function initialize() {
  console.debug("[initialize] started");

  const router = new Router();
  router.on(saveGolinkUrlName, saveGolinkUrl);
  router.on(fetchGolinkUrlName, fetchGolinkUrl);
  router.on(isManagedName, isManaged);
  chrome.runtime.onMessage.addListener(router.listener());

  await updateRedirectRule();

  console.debug("[initialize] finished");
}

function onInstalled() {
  console.debug("[onInstalled] started");
  (async () => {
    try {
      await initialize();
    } catch (e) {
      console.error("[onInstalled]", e);
    }
  })();

  console.debug("[onInstalled] finished");
  return true;
}

function onStartup() {
  console.log("[onStartup]");
  (async () => {
    await initialize();
  })();

  return true;
}

chrome.runtime.onInstalled.addListener(onInstalled);
// chrome.runtime.onStartup.addListener(onStartup);
