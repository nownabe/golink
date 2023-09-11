import { updateRedirectRule } from "../background/updateRedirectRule";
import { getGolinkUrl, getIsManaged, setGolinkUrl } from "../storage";

async function onSave() {
  const url = (<HTMLInputElement>document.getElementById("option-url")).value;
  try {
    await setGolinkUrl(url);
    await updateRedirectRule(url);
    alert("Saved!");
  } catch (e) {
    console.error("[options.onSave] saveGolinkUrl failed:", e);
    alert(`Failed to save.`);
  }
}

async function initialize() {
  console.debug("[initialize] starged");

  const url = await getGolinkUrl();
  const isManaged = await getIsManaged();

  const input = <HTMLInputElement>document.getElementById("option-url");
  const button = <HTMLButtonElement>document.getElementById("save");

  input.value = url || "";

  if (isManaged) {
    document.getElementById("managed")!.hidden = false;
  } else {
    input.disabled = false;
    button.disabled = false;
  }

  console.debug("[initialize] finished");
}

function onDOMContentLoaded() {
  console.debug("[onDOMContentLoaded] started");
  (async () => {
    try {
      await initialize();
    } catch (e) {
      console.error("[onDOMContentLoaded] initialize failed:", e);
    }
  })();
  console.debug("[onDOMContentLoaded] finished");

  return true;
}

document.addEventListener("DOMContentLoaded", onDOMContentLoaded);
document.getElementById("save")?.addEventListener("click", onSave);
