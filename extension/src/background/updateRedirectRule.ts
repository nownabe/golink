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

  console.debug("[updateRedirectRule] finished");
}
