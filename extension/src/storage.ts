const golinkUrlKey = "golinkUrl" as const;
const isFinishedFirstOpenKey = "isFinishedFirstOpen" as const;
const managedInstanceUrlKey = "golinkInstanceUrl" as const;

export async function getGolinkUrl(): Promise<string | null> {
  console.debug(`[getGolinkUrl] started`);

  const managedResult = await chrome.storage.managed.get(managedInstanceUrlKey);
  if (managedResult && managedInstanceUrlKey in managedResult) {
    const url = managedResult[managedInstanceUrlKey] || null;
    console.debug(`[getGolinkUrl] got url from managed storage`, url);
    return url;
  }

  const syncResult = await chrome.storage.sync.get(golinkUrlKey);
  if (syncResult && golinkUrlKey in syncResult) {
    const url = syncResult[golinkUrlKey] || null;
    console.debug(`[getGolinkUrl] got url from sync storage`, url);
    return url;
  }

  console.debug(`[getGolinkUrl] url not found in storage`);
  return null;
}

export async function setGolinkUrl(url: string) {
  console.debug(`[setGolinkUrl] started setting URL ${url}`);
  await chrome.storage.sync.set({ [golinkUrlKey]: url });
  console.debug(`[setGolinkUrl] finished setting URL ${url}`);
}

export async function getIsManaged(): Promise<boolean> {
  console.debug(`[isManaged] started`);
  const result = await chrome.storage.managed.get(managedInstanceUrlKey);
  const isManaged = result && managedInstanceUrlKey in result ? true : false;
  console.debug(`[isManaged] isManaged = ${isManaged}`);
  console.debug(`[isManaged] finished`);
  return isManaged;
}

type listenerFn = (newVal: string | null, oldVal: string | null) => void;

export function addListenerOnGolinkUrlChanged(callback: listenerFn) {
  chrome.storage.onChanged.addListener((changes, areaName) => {
    console.log(`[onChanged] started`);
    if (areaName === "sync" && changes[golinkUrlKey]) {
      callback(
        changes[golinkUrlKey].newValue || null,
        changes[golinkUrlKey].oldValue || null
      );
    } else if (areaName === "managed" && changes[managedInstanceUrlKey]) {
      callback(
        changes[managedInstanceUrlKey].newValue || null,
        changes[managedInstanceUrlKey].oldValue || null
      );
    }
    console.debug(`[onChanged] finished`);
  });
}

export async function getIsFinishedFirstOpen(): Promise<boolean> {
  console.debug(`[getIsFinishedFirstOpen] started`);
  const result = await chrome.storage.local.get(isFinishedFirstOpenKey);
  const isFinished = result && result[isFinishedFirstOpenKey] ? true : false;
  console.debug(`[getIsFinishedFirstOpen] isFinished = ${isFinished}`);
  console.debug(`[getIsFinishedFirstOpen] finished`);
  return isFinished;
}

export async function setIsFinishedFirstOpen(
  isFinished: boolean
): Promise<void> {
  console.debug(`[setIsFinishedFirstOpen] started`);
  await chrome.storage.local.set({ [isFinishedFirstOpenKey]: isFinished });
  console.debug(`[setIsFinishedFirstOpen] finished`);
}
