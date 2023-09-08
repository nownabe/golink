import { addListenerOnGolinkUrlChanged, getGolinkUrl } from "../storage";
import { updateRedirectRule } from "./updateRedirectRule";

async function initialize() {
  console.debug("[initialize] started");

  const url = await getGolinkUrl();
  if (url) {
    await updateRedirectRule(url);
  }

  addListenerOnGolinkUrlChanged((newUrl, oldUrl) => {
    console.debug(`[golinkChanged] newUrl = '${newUrl}', oldUrl = '${oldUrl}'`);
    (async () => {
      if (newUrl) {
        updateRedirectRule(newUrl);
      }
    })();
  });

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

// function onStartup() {
//   console.log("[onStartup]");
//   (async () => {
//     await initialize();
//   })();

//   return true;
// }

chrome.runtime.onInstalled.addListener(onInstalled);
// chrome.runtime.onStartup.addListener(onStartup);
