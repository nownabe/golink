import { getGolinkUrl } from "./storage";

export async function updateRedirectRule() {
  console.debug("[updateRedirectRule] started");
  const ruleId = 1;

  const url = await getGolinkUrl();

  if (!url) {
    console.error("[updateRedirectRule] golink url is not configured");
    return;
  }

  console.debug(`[updateRedirectRule] updating redirect rule to ${url}`);
  let host;
  try {
    host = new URL(url).host;
  } catch (e) {
    console.error("Invalid URL:", url);
    console.error(e);
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

  await chrome.declarativeNetRequest.updateDynamicRules(updateRuleOptions);
  console.log(`[updateRedirectRule] updated redirect rule to ${url}`);
  console.debug("[updateRedirectRule] finished");
}
