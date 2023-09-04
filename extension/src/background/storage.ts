const golinkUrlKey = "golinkUrl";
const managedInstanceUrlKey = "golinkInstanceUrl";

export async function getGolinkUrl(): Promise<string | undefined> {
  console.debug(`[getGolinkUrl] started`);
  const result = await chrome.storage.managed.get(managedInstanceUrlKey);
  if (result && result[managedInstanceUrlKey]) {
    console.debug(`[getGolinkUrl] got url from managed storage`, result);
    return result[managedInstanceUrlKey];
  }

  const url = (await chrome.storage.sync.get(golinkUrlKey))[golinkUrlKey];
  console.debug(`[getGolinkUrl] got url from sync storage`, url);
  console.debug(`[getGolinkUrl] finished`);
  return url;
}

export async function setGolinkUrl(url: string) {
  console.debug(`[setGolinkUrl] started setting URL ${url}`);
  await chrome.storage.sync.set({ [golinkUrlKey]: url });
  console.debug(`[setGolinkUrl] finished setting URL ${url}`);
}

export async function getIsManaged(): Promise<boolean> {
  console.debug(`[isManaged] started`);
  const result = await chrome.storage.managed.get(managedInstanceUrlKey);
  console.debug(`[isManaged] finished`);
  return result && result[managedInstanceUrlKey] ? true : false;
}
