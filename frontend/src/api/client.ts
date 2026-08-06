import type { ApiError } from "../types";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL as string;

export class ApiRequestError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

let authToken: string | null = localStorage.getItem("token");

export function setAuthToken(token: string | null) {
  authToken = token;
  if (token) {
    localStorage.setItem("token", token);
  } else {
    localStorage.removeItem("token");
  }
}

export function getAuthToken() {
  return authToken;
}

interface RequestOptions {
  method?: "GET" | "POST" | "PUT" | "DELETE";
  body?: unknown;
}

async function rawRequest(path: string, options: RequestOptions): Promise<Response> {
  const headers: Record<string, string> = {};
  if (options.body !== undefined) {
    headers["Content-Type"] = "application/json";
  }
  if (authToken) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }

  return fetch(`${API_BASE_URL}${path}`, {
    method: options.method ?? "GET",
    headers,
    body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
  });
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const res = await rawRequest(path, options);

  if (!res.ok) {
    let message = `Request failed with status ${res.status}`;
    try {
      const data = (await res.json()) as ApiError;
      if (data.error) message = data.error;
    } catch {
      // response wasn't JSON; keep the default message
    }
    throw new ApiRequestError(res.status, message);
  }

  if (res.status === 204) {
    return undefined as T;
  }
  return (await res.json()) as T;
}

// apiRequestAllowStatus is for endpoints where a non-2xx status is still a
// meaningful, expected business outcome (e.g. a simulated card decline),
// not a transport/client error that should be thrown.
export async function apiRequestAllowStatus<T>(
  path: string,
  allowedStatuses: number[],
  options: RequestOptions = {},
): Promise<T> {
  const res = await rawRequest(path, options);

  if (!res.ok && !allowedStatuses.includes(res.status)) {
    let message = `Request failed with status ${res.status}`;
    try {
      const data = (await res.json()) as ApiError;
      if (data.error) message = data.error;
    } catch {
      // response wasn't JSON; keep the default message
    }
    throw new ApiRequestError(res.status, message);
  }

  return (await res.json()) as T;
}
