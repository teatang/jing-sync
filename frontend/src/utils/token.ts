export function setToken(token: string): void {
  localStorage.setItem("token", token);
}

export function getToken(): string {
  return localStorage.getItem("token") == null ? "" : localStorage.getItem("token") as string;
}

export function removeToken(): void {
  localStorage.removeItem("token");
}
