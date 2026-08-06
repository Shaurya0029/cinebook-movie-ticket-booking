import { createContext, useContext, useEffect, useState, type ReactNode } from "react";
import * as citiesApi from "../api/cities";
import * as theatersApi from "../api/theaters";
import type { City, Theater } from "../types";

const CITY_STORAGE_KEY = "selected_city";
const THEATER_STORAGE_KEY = "selected_theater";

interface LocationContextValue {
  selectedCity: City | null;
  selectedTheater: Theater | null;
  setCity: (city: City) => void;
  setTheater: (theater: Theater | null) => void;
  useCurrentLocation: () => Promise<void>;
  locating: boolean;
  locationError: string | null;
}

const LocationContext = createContext<LocationContextValue | undefined>(undefined);

export function LocationProvider({ children }: { children: ReactNode }) {
  const [selectedCity, setSelectedCity] = useState<City | null>(() => readStored(CITY_STORAGE_KEY));
  const [selectedTheater, setSelectedTheater] = useState<Theater | null>(() => readStored(THEATER_STORAGE_KEY));
  const [locating, setLocating] = useState(false);
  const [locationError, setLocationError] = useState<string | null>(null);

  useEffect(() => {
    if (selectedCity) {
      localStorage.setItem(CITY_STORAGE_KEY, JSON.stringify(selectedCity));
    }
  }, [selectedCity]);

  useEffect(() => {
    if (selectedTheater) {
      localStorage.setItem(THEATER_STORAGE_KEY, JSON.stringify(selectedTheater));
    } else {
      localStorage.removeItem(THEATER_STORAGE_KEY);
    }
  }, [selectedTheater]);

  function setCity(city: City) {
    setSelectedCity(city);
    setSelectedTheater(null);
  }

  function setTheater(theater: Theater | null) {
    setSelectedTheater(theater);
  }

  async function useCurrentLocation() {
    setLocating(true);
    setLocationError(null);
    try {
      const position = await getCurrentPosition();
      const { latitude, longitude } = position.coords;
      const [city, theater] = await Promise.all([
        citiesApi.nearestCity(latitude, longitude),
        theatersApi.nearestTheater(latitude, longitude),
      ]);
      setSelectedCity(city);
      setSelectedTheater(theater);
    } catch (err) {
      setLocationError(err instanceof Error ? err.message : "Could not determine your location");
    } finally {
      setLocating(false);
    }
  }

  return (
    <LocationContext.Provider
      value={{ selectedCity, selectedTheater, setCity, setTheater, useCurrentLocation, locating, locationError }}
    >
      {children}
    </LocationContext.Provider>
  );
}

export function useLocation() {
  const ctx = useContext(LocationContext);
  if (!ctx) throw new Error("useLocation must be used within a LocationProvider");
  return ctx;
}

function readStored<T>(key: string): T | null {
  const raw = localStorage.getItem(key);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as T;
  } catch {
    return null;
  }
}

function getCurrentPosition(): Promise<GeolocationPosition> {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      reject(new Error("Geolocation is not supported by this browser"));
      return;
    }
    navigator.geolocation.getCurrentPosition(resolve, () => reject(new Error("Location permission denied")), {
      timeout: 10000,
    });
  });
}
