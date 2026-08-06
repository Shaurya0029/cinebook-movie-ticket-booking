import { apiRequest } from "./client";
import type { City, Theater } from "../types";

export function listCities() {
  return apiRequest<City[]>("/cities");
}

export function nearestCity(lat: number, lng: number) {
  return apiRequest<City>(`/cities/nearest?lat=${lat}&lng=${lng}`);
}

export function listTheatersByCity(cityId: number) {
  return apiRequest<Theater[]>(`/cities/${cityId}/theaters`);
}
