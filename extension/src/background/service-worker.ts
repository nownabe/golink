import { saveGolinkUrl, saveGolinkUrlName } from "../messageListeners";
import { Router, SimpleResponse } from "../router";

export const golinkUrlKey = "golinkUrl";

export async function updateRedirectRule(url: string) {
  const ruleId = 1;

  console.log(`[updateRedirectRule] updating redirect rule to ${url}`);
  let host;
  try {
    host = new URL(url).host;
  } catch (e) {
    console.log("Invalid URL:", url);
    console.log(e);
    return;
  }

  const redirectRule = {
    id: ruleId,
    priority: 1,
    action: {
      type: "redirect",
      redirect: {
        transform: { scheme: "https", host: host },
      },
    },
    condition: {
      urlFilter: "*://go/*",
      resourceTypes: ["main_frame"],
    },
  };

  const updateRuleOptions = {
    removeRuleIds: [ruleId],
    addRules: [redirectRule],
  };

  await chrome.declarativeNetRequest.updateDynamicRules(updateRuleOptions);
  console.log("[updateRedirectRule] updated redirect rule");
}

async function initialize() {
  console.log("[initialize] initializing");

  const url = (await chrome.storage.sync.get(golinkUrlKey))[golinkUrlKey];

  await updateRedirectRule(url);

  chrome.storage.onChanged.addListener(
    (
      changes: { [key: string]: chrome.storage.StorageChange },
      namespace: string
    ) => {
      (async () => {
        if (namespace === "sync" && golinkUrlKey in changes) {
          console.log(
            "[storage.onChanged] Golink URL changed",
            changes[golinkUrlKey].newValue
          );
        }
      })();
    }
  );

  const router = new Router();
  router.on(saveGolinkUrlName, saveGolinkUrl);
  chrome.runtime.onMessage.addListener(router.listener());

  console.log("[initialize] initialized");
}

function onInstalled() {
  console.log("[onInstalled]");
  (async () => {
    await initialize();
  })();
}

function onStartup() {
  console.log("[onStartup]");
  (async () => {
    await initialize();
  })();
}

chrome.runtime.onInstalled.addListener(onInstalled);
chrome.runtime.onStartup.addListener(onStartup);
