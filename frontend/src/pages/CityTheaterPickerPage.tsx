import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import * as citiesApi from "../api/cities";
import { useLocation as useCinemaLocation } from "../context/LocationContext";
import type { City, Theater } from "../types";
import styles from "./CityTheaterPickerPage.module.css";

export default function CityTheaterPickerPage() {
  const { selectedCity, selectedTheater, setCity, setTheater, useCurrentLocation, locating, locationError } =
    useCinemaLocation();
  const navigate = useNavigate();

  const [cities, setCities] = useState<City[] | null>(null);
  const [theaters, setTheaters] = useState<Theater[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    citiesApi
      .listCities()
      .then(setCities)
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load cities"));
  }, []);

  useEffect(() => {
    if (!selectedCity) {
      setTheaters(null);
      return;
    }
    citiesApi
      .listTheatersByCity(selectedCity.id)
      .then(setTheaters)
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load theaters"));
  }, [selectedCity]);

  async function handleUseLocation() {
    await useCurrentLocation();
  }

  function handleTheaterPick(theater: Theater) {
    setTheater(theater);
    navigate(-1);
  }

  return (
    <div className="page container fadeIn">
      <h1 className="sectionTitle" style={{ marginTop: 8 }}>
        Choose your location
      </h1>

      <div className={`card ${styles.locateCard} fadeInUp`}>
        <div>
          <div className={styles.locateTitle}>Use my current location</div>
          <div className={styles.locateDesc}>We'll match you to the nearest city and theater automatically.</div>
          {locationError && <div className="errorText" style={{ marginTop: 8 }}>{locationError}</div>}
        </div>
        <button className="btn btnPrimary" onClick={handleUseLocation} disabled={locating}>
          {locating ? <span className="spinner" /> : "Use current location"}
        </button>
      </div>

      {error && <p className="errorText">{error}</p>}

      <h2 className="sectionTitle">Or pick a city</h2>
      <div className={styles.cityGrid}>
        {(cities ?? []).map((city, i) => (
          <button
            key={city.id}
            className={`card ${styles.cityButton} ${selectedCity?.id === city.id ? styles.active : ""} fadeInUp`}
            style={{ animationDelay: `${i * 30}ms` }}
            onClick={() => setCity(city)}
          >
            <div className={styles.cityName}>{city.name}</div>
            {city.state && <div className={styles.cityState}>{city.state}</div>}
          </button>
        ))}
      </div>

      {selectedCity && (
        <>
          <h2 className="sectionTitle">Theaters in {selectedCity.name}</h2>
          <div className={styles.theaterList}>
            {(theaters ?? []).map((theater, i) => (
              <button
                key={theater.id}
                className={`card ${styles.theaterRow} fadeInUp`}
                style={{ animationDelay: `${i * 30}ms` }}
                onClick={() => handleTheaterPick(theater)}
              >
                <div>
                  <div className={styles.theaterName}>{theater.name}</div>
                  {theater.address && <div className={styles.theaterAddress}>{theater.address}</div>}
                </div>
                {selectedTheater?.id === theater.id && <span style={{ color: "var(--color-accent)" }}>Selected</span>}
              </button>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
