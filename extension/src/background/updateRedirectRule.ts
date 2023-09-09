import { getIsFinishedFirstOpen, setIsFinishedFirstOpen } from "../storage";

const ruleId = 1;

export async function updateRedirectRule(url: string) {
  console.debug("[updateRedirectRule] started");

  if (!url) {
    console.error(`[updateRedirectRule] golink url is empty: '${url}`);
    return;
  }

  let host;
  try {
    host = new URL(url).host;
  } catch (e) {
    console.error(`[updateRedirectRUle] invalid url: '${url}:`, e);
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
  console.debug("[updateRedirectRule] updateRuleOptions", updateRuleOptions);

  console.debug(`[updateRedirectRule] updating redirect rule to ${url}`);
  await chrome.declarativeNetRequest.updateDynamicRules(updateRuleOptions);

  // To tell Chrome Browser that http://go/ is a valid url.
  if (!(await getIsFinishedFirstOpen())) {
    await openGoTab(url);
  }

  console.debug("[updateRedirectRule] finished");
}

async function openGoTab(url: string) {
  console.debug(`[openGoTab] started`);

  if (!url.endsWith("/")) {
    url += "/";
  }
  const consoleUrl = url + "-/";

  const goTab = await chrome.tabs.create({ url: "http://go/" });
  console.debug(`[openGoTab] opened http://go/:`, goTab);

  const onUpdated = (
    tabId: number,
    changeInfo: chrome.tabs.TabChangeInfo,
    tab: chrome.tabs.Tab
  ) => {
    if (tabId === goTab.id && changeInfo.status === "complete") {
      console.debug("[openGoTab] loading goTab is completed", tab);

      chrome.tabs.remove(tabId);
      console.debug("[openGoTab] removed goTab");

      chrome.tabs.onUpdated.removeListener(onUpdated);
      console.debug(`[openGoTab] removed listener`);

      if (tab.url === consoleUrl) {
        console.debug(
          `[openGoTab] succeed to open and redirect http://go/ to ${consoleUrl}`
        );
        (async () => {
          await setIsFinishedFirstOpen(true);
        })();
      }
    }
  };

  chrome.tabs.onUpdated.addListener(onUpdated);

  console.debug(`[openGoTab] finished`);
}
