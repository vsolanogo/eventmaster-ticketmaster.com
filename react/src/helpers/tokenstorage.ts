const tokenKey = "access_token";

export function isLocalStorageSupported(): boolean {
  try {
    localStorage.setItem("test", "test");
    localStorage.removeItem("test");
    return true;
  } catch (e) {
    return false;
  }
}

export function storeToken(token: string): void {
  if (isLocalStorageSupported()) {
    localStorage.setItem(tokenKey, token);
  } else {
    sessionStorage.setItem(tokenKey, token);
  }
}

export function getToken(): string | null {
  if (isLocalStorageSupported()) {
    return localStorage.getItem(tokenKey);
  } else {
    return sessionStorage.getItem(tokenKey);
  }
}

export function removeToken(): void {
  if (isLocalStorageSupported()) {
    localStorage.removeItem(tokenKey);
  } else {
    sessionStorage.removeItem(tokenKey);
  }
}
