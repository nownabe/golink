import { HandlerFn } from "./router";
import { getGolinkUrl, setGolinkUrl, getIsManaged } from "./storage";
import { updateRedirectRule } from "./updateRedirectRule";

export type SaveGolinkUrlRequestData = { url: string };
export type SaveGolinkUrlResponseData = {};
export const saveGolinkUrlName = "saveGolinkUrl";
export const saveGolinkUrl: HandlerFn<
  SaveGolinkUrlRequestData,
  SaveGolinkUrlResponseData
> = async ({ data }) => {
  const { url } = data;
  await setGolinkUrl(url);
  await updateRedirectRule();
  return {};
};

export type FetchGolinkUrlRequestData = {};
export type FetchGolinkUrlResponseData = { url: string | undefined };
export const fetchGolinkUrlName = "fetchGolinkUrl";
export const fetchGolinkUrl: HandlerFn<
  FetchGolinkUrlRequestData,
  FetchGolinkUrlResponseData
> = async () => {
  const url = await getGolinkUrl();
  return { url };
};

export type IsManagedRequestData = {};
export type IsManagedResponseData = { isManaged: boolean };
export const isManagedName = "isManaged";
export const isManaged: HandlerFn<
  IsManagedRequestData,
  IsManagedResponseData
> = async () => {
  const isManaged = await getIsManaged();
  return { isManaged };
};
