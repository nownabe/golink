const host = "your-host.com";

async function initialize() {
  const redirectRule = {
    id: 1,
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
    removeRuleIds: [1],
    addRules: [redirectRule],
  };

  await chrome.declarativeNetRequest.updateDynamicRules(updateRuleOptions);
}

chrome.runtime.onInstalled.addListener(initialize);
