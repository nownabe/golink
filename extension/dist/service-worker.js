(()=>{async function e(){await chrome.declarativeNetRequest.updateDynamicRules({removeRuleIds:[1],addRules:[{id:1,priority:1,action:{type:"redirect",redirect:{transform:{scheme:"https",host:"your-host.com"}}},condition:{urlFilter:"*://go/*",resourceTypes:["main_frame"]}}]})}chrome.runtime.onInstalled.addListener(e)})();
//# sourceMappingURL=service-worker.js.map
