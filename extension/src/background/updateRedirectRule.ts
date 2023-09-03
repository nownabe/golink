import { getGolinkUrl } from "./storage";

export async function updateRedirectRule() {
  const ruleId = 1;

  const url = await getGolinkUrl();

  if (!url) {
    console.error("[updateRedirectRule] golink url is not configured");
    return;
  }

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
