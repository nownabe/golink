export function validateEmail(email: string): boolean {
  const re = /\S+@\S+\.\S+/;
  return re.test(email);
}

export function validateGolinkName(name: string): boolean {
  if (name === "") {
    return false;
  }

  if (name.startsWith("-") || name.endsWith("-")) {
    return false;
  }

  if (name.match(/^_+$/) || name.startsWith("__") || name.endsWith("__")) {
    return false;
  }

  return true;
}

export function validateUrl(url: string): boolean {
  try {
    new URL(url);
    return true;
  } catch (e) {
    return false;
  }
}
