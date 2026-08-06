import { apiRequest } from "./client";
import type { AuthResponse, User } from "../types";

export function register(firstName: string, lastName: string, email: string, password: string) {
  return apiRequest<AuthResponse>("/auth/register", {
    method: "POST",
    body: { first_name: firstName, last_name: lastName, email, password },
  });
}

export function login(email: string, password: string) {
  return apiRequest<AuthResponse>("/auth/login", {
    method: "POST",
    body: { email, password },
  });
}

export function me() {
  return apiRequest<User>("/auth/me");
}
