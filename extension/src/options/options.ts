const storageKey = "golinkUrl";

async function onSave() {
  const url = (<HTMLInputElement>document.getElementById("option-url")).value;
  await chrome.storage.sync.set({ [storageKey]: url });
  const response = await chrome.runtime.sendMessage({
    type: "saveGolinkUrl",
    url: url,
  });
  console.log(response);
  alert("Saved!");
}

async function restoreOptions() {
  const url = (await chrome.storage.sync.get(storageKey))[storageKey];
  const input = <HTMLInputElement>document.getElementById("option-url");
  if (input) {
    input.value = url ?? "";
  }
}

document.addEventListener("DOMContentLoaded", restoreOptions);
document.getElementById("save")?.addEventListener("click", onSave);
