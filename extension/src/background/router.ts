export type Request<T> = {
  name: string;
  data: T;
};

export type RequestMessage<T> = {
  request: Request<T>;
  sender: chrome.runtime.MessageSender;
  sendResponse: (response?: any) => void;
};

export type Response<T> =
  | {
      success: true;
      data: T;
    }
  | {
      success: false;
      error: string;
    };

export type HandlerFn<T, S> = (request: Request<T>) => Promise<S>;

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
        try {
          const responseData = await handler({ name, data });
          const response = { success: true, data: responseData };
          sendResponse(response);
          console.log(
            `[Router][${name}] finished handling message with`,
            response
          );
        } catch (e) {
          console.error(`[Router][${name}] failed to handle message`, e);

          var msg;
          if (e instanceof Error) {
            msg = e.message;
          } else {
            msg = "unknown message";
          }
          sendResponse({
            success: false,
            error: msg,
          });
        }
      })();

      // Returning true indicates that the response will be sent asynchronously.
      return true;
    };
  }
}
