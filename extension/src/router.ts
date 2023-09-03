export interface Request<T> {
  name: string;
  data: T;
}

export type RequestMessage<T> = {
  request: Request<T>;
  sender: chrome.runtime.MessageSender;
  sendResponse: (response?: any) => void;
};

export interface Response<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface SimpleResponse extends Response<void> {
  success: boolean;
}

export type HandlerFn<T, S> = (request: Request<T>) => Promise<Response<S>>;

type listenerFn = (
  request: Request<any>,
  sender: chrome.runtime.MessageSender,
  sendResponse: (response?: any) => void
) => boolean;

export class Router {
  routes: Map<string, HandlerFn<any, any>>;

  constructor() {
    this.routes = new Map();
  }

  public on<T, S>(name: string, fn: HandlerFn<T, S>) {
    this.routes.set(name, fn);
  }

  public listener(): listenerFn {
    return (request, sender, sendResponse) => {
      console.log("[Router] received message", request, sender);

      const { name, data } = request;

      const handler = this.routes.get(name);

      if (!handler) {
        console.error(`[Router] no handler found for ${name}`);
        sendResponse({
          success: false,
          error: `[Router] no handler found for ${name}`,
        });
        return false;
      }

      (async () => {
        console.log(`[Router][${name}] started handling message with`, data);
        const response = await handler({ name, data });
        sendResponse(response);
        console.log(
          `[Router][${name}] finished handling message with`,
          response
        );
      })();

      // Returning true indicates that the response will be sent asynchronously.
      return true;
    };
  }
}

export async function send<T, S>(name: string, req: T): Promise<S> {
  console.log(`[send] sending request`, req);
  const msg: Request<T> = {
    name,
    data: req,
  };
  const response = await chrome.runtime.sendMessage(msg);
  console.log(`[send] received response`, response);
  return response;
}
