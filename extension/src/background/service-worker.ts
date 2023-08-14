const golinkUrlKey = "golinkUrl";

async function updateRedirectRule(url: string) {
  const ruleId = 1;

  console.log(`Updating redirect rule to ${url}`);
  const host = new URL(url).host;

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
}

async function initialize() {
  const url = (await chrome.storage.sync.get(golinkUrlKey))[golinkUrlKey];

  await updateRedirectRule(url);

  chrome.storage.onChanged.addListener(
    async (
      changes: { [key: string]: chrome.storage.StorageChange },
      namespace: string,
    ) => {
      if (namespace === "sync" && golinkUrlKey in changes) {
        console.log("Golink URL changed", changes[golinkUrlKey].newValue);
        await updateRedirectRule(changes[golinkUrlKey].newValue);
      }
    },
  );
}

chrome.runtime.onInstalled.addListener(initialize);
