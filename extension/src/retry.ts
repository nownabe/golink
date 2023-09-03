function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export async function retry(fn: Function, maxAttempts = 5) {
  let attempts = 0;
  let backoff = 100;

  while (true) {
    try {
      return await fn();
    } catch (e) {
      console.debug(`[retry] attempt ${attempts} failed`, e);
      attempts++;
      if (attempts >= maxAttempts) {
        throw e;
      }
      await sleep(backoff);
      backoff = backoff * 2;
    }
  }
}
