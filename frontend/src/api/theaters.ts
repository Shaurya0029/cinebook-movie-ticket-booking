import { apiRequest } from "./client";
import type { Theater } from "../types";

export function nearestTheater(lat: number, lng: number) {
  return apiRequest<Theater>(`/theaters/nearest?lat=${lat}&lng=${lng}`);
}

export function getTheater(theaterId: number) {
  return apiRequest<Theater>(`/theaters/${theaterId}`);
}
