const golinkUrlKey = "golinkUrl";
const managedInstanceUrlKey = "golinkInstanceUrl";

export async function getGolinkUrl(): Promise<string | undefined> {
  console.debug(`[storage.getGolinkUrl] started`);
  const result = await chrome.storage.managed.get(managedInstanceUrlKey);
  if (result && result[managedInstanceUrlKey]) {
    console.debug(
      `[storage.getGolinkUrl] got url from managed storage`,
      result
    );
    return result[managedInstanceUrlKey];
  }

  const url = (await chrome.storage.sync.get(golinkUrlKey))[golinkUrlKey];
  console.debug(`[storage.getGolinkUrl] got url from sync storage`, url);
  console.debug(`[storage.getGolinkUrl] finished`);
  return url;
}

export async function setGolinkUrl(url: string) {
  console.debug(`[storage.setGolinkUrl] started setting URL ${url}`);
  await chrome.storage.sync.set({ [golinkUrlKey]: url });
  console.debug(`[storage.setGolinkUrl] finished setting URL ${url}`);
}

export async function getIsManaged(): Promise<boolean> {
  console.debug(`[storage.isManaged] started`);
  const result = await chrome.storage.managed.get(managedInstanceUrlKey);
  console.debug(`[storage.isManaged] finished`);
  return result && result[managedInstanceUrlKey] ? true : false;
}
