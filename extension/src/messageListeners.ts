import { HandlerFn } from "./router";
import { golinkUrlKey, updateRedirectRule } from "./background/service-worker";

export type SaveGolinkUrlRequestData = { url: string };
export type SaveGolinkUrlResponseData = {};
export const saveGolinkUrlName = "saveGolinkUrl";
export const saveGolinkUrl: HandlerFn<
  SaveGolinkUrlRequestData,
  SaveGolinkUrlResponseData
> = async ({ data }) => {
  const { url } = data;
  await chrome.storage.sync.set({ [golinkUrlKey]: url });
  await updateRedirectRule(url);
  return { success: true };
};
