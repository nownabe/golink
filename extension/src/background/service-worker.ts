import {
  fetchGolinkUrl,
  fetchGolinkUrlName,
  isManaged,
  isManagedName,
  saveGolinkUrl,
  saveGolinkUrlName,
} from "./messageListeners";
import { Router } from "./router";

async function initialize() {
  console.log("[initialize] initializing");

  const router = new Router();
  router.on(saveGolinkUrlName, saveGolinkUrl);
  router.on(fetchGolinkUrlName, fetchGolinkUrl);
  router.on(isManagedName, isManaged);
  chrome.runtime.onMessage.addListener(router.listener());

  // await updateRedirectRule();
  console.log("[initialize] initialized");
}

function onInstalled() {
  console.log("[onInstalled] started");
  (async () => {
    try {
      await initialize();
    } catch (e) {
      console.error("[onInstalled]", e);
    }
  })();

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
