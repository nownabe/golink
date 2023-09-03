import {
  fetchGolinkUrl,
  fetchIsManaged,
  saveGolinkUrl,
} from "../backgroundServices";

async function onSave() {
  const url = (<HTMLInputElement>document.getElementById("option-url")).value;
  try {
    await saveGolinkUrl(url);
    alert("Saved!");
  } catch (e) {
    console.error("[options.onSave] saveGolinkUrl failed:", e);
    alert(`Failed to save.`);
  }
}

async function initialize() {
  const url = (await fetchGolinkUrl()) ?? "";
  const input = <HTMLInputElement>document.getElementById("option-url");
  input.value = url;

  const isManaged = await fetchIsManaged();
  if (isManaged) {
    input.disabled = true;
    (<HTMLButtonElement>document.getElementById("save")).disabled = true;
    document.getElementById("managed")!.hidden = false;
  }
}

function onDOMContentLoaded() {
  console.debug("[options.onDOMContentLoaded] started");
  (async () => {
    try {
      await initialize();
    } catch (e) {
      console.error("[options.onDOMContentLoaded] initialize failed:", e);
    }
  })();
  console.debug("[options.onDOMContentLoaded] finished");

  return true;
}

document.addEventListener("DOMContentLoaded", onDOMContentLoaded);
document.getElementById("save")?.addEventListener("click", onSave);
