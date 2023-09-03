import {
  FetchGolinkUrlRequestData,
  FetchGolinkUrlResponseData,
  IsManagedRequestData,
  IsManagedResponseData,
  SaveGolinkUrlRequestData,
  SaveGolinkUrlResponseData,
  fetchGolinkUrlName,
  isManagedName,
  saveGolinkUrlName,
} from "./background/messageListeners";
import { Request, Response } from "./background/router";
import { retry } from "./retry";

async function send<T, S>(name: string, req: T): Promise<Response<S>> {
  console.log(`[send] sending request to ${name}`, req);
  const msg: Request<T> = {
    name,
    data: req,
  };

  const response = await retry(
    async () => await chrome.runtime.sendMessage(msg)
  );
  console.log(`[send] received response`, response);
  return response;
}

export async function saveGolinkUrl(url: string) {
  const response = await send<
    SaveGolinkUrlRequestData,
    SaveGolinkUrlResponseData
  >(saveGolinkUrlName, { url });
  if (!response.success) {
    throw new Error(`${saveGolinkUrlName} failed: ${response.error}`);
  }
}

export async function fetchGolinkUrl(): Promise<string | undefined> {
  const response = await send<
    FetchGolinkUrlRequestData,
    FetchGolinkUrlResponseData
  >(fetchGolinkUrlName, {});
  if (!response.success) {
    throw new Error(`${fetchGolinkUrl} failed: ${response.error}`);
  }
  return response.data.url;
}

export async function fetchIsManaged(): Promise<boolean> {
  const response = await send<IsManagedRequestData, IsManagedResponseData>(
    isManagedName,
    {}
  );
  if (!response.success) {
    throw new Error(`${isManagedName} failed: ${response.error}`);
  }
  return response.data.isManaged;
}
