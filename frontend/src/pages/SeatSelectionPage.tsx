import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import * as showsApi from "../api/shows";
import * as bookingsApi from "../api/bookings";
import { ApiRequestError } from "../api/client";
import SeatMap from "../components/SeatMap";
import PriceBreakdown from "../components/PriceBreakdown";
import { formatShowDateTime } from "../utils/format";
import type { Seat, Show } from "../types";
import styles from "./SeatSelectionPage.module.css";

const GST_RATE_PERCENT = 18; // mirrors the backend's default; the booking response is authoritative.
const MAX_SEATS = 8;

export default function SeatSelectionPage() {
  const { showId } = useParams();
  const navigate = useNavigate();

  const [show, setShow] = useState<Show | null>(null);
  const [seats, setSeats] = useState<Seat[] | null>(null);
  const [selected, setSelected] = useState<Map<number, Seat>>(new Map());
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!showId) return;
    showsApi
      .getShow(Number(showId))
      .then(setShow)
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load show"));
    showsApi
      .listSeatsForShow(Number(showId))
      .then((data) => setSeats(data ?? []))
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load seats"));
  }, [showId]);

  function toggleSeat(seat: Seat) {
    setSelected((prev) => {
      const next = new Map(prev);
      if (next.has(seat.id)) {
        next.delete(seat.id);
      } else {
        if (next.size >= MAX_SEATS) return prev;
        next.set(seat.id, seat);
      }
      return next;
    });
  }

  async function handleProceed() {
    if (!show || selected.size === 0) return;
    setSubmitting(true);
    setError(null);
    try {
      const booking = await bookingsApi.createBooking(show.id, Array.from(selected.keys()));
      navigate(`/bookings/${booking.id}/pay`);
    } catch (err) {
      if (err instanceof ApiRequestError && err.status === 409) {
        setError("One or more selected seats were just taken. Please pick different seats.");
        const fresh = await showsApi.listSeatsForShow(show.id);
        setSeats(fresh ?? []);
        setSelected(new Map());
      } else {
        setError(err instanceof Error ? err.message : "Failed to create booking");
      }
    } finally {
      setSubmitting(false);
    }
  }

  if (error && !show) {
    return (
      <div className="page container">
        <p className="errorText">{error}</p>
      </div>
    );
  }

  if (!show || !seats) {
    return (
      <div className="page container">
        <div className="skeleton" style={{ height: 420, borderRadius: 16 }} />
      </div>
    );
  }

  const selectedIds = new Set(selected.keys());

  return (
    <div className="page container fadeIn">
      <div className={styles.header}>
        <h1 className={styles.movieTitle}>{show.movie_title}</h1>
        <div className={styles.showMeta}>
          {show.theater_name} · {show.city_name} · {formatShowDateTime(show.starts_at)}
        </div>
      </div>

      <div className={styles.layout}>
        <div className={`card ${styles.seatCard}`}>
          <SeatMap seats={seats} selectedIds={selectedIds} onToggle={toggleSeat} />
        </div>

        <div className={styles.summary}>
          <div className={`card ${styles.selectedSeats}`}>
            <div className={styles.selectedLabel}>
              {selected.size === 0 ? `Select up to ${MAX_SEATS} seats` : `${selected.size} seat(s) selected`}
            </div>
            <div className={styles.seatChips}>
              {Array.from(selected.values()).map((s) => (
                <span key={s.id} className={styles.seatChip}>
                  {s.seat_label}
                </span>
              ))}
            </div>
          </div>

          {selected.size > 0 ? (
            <PriceBreakdown seatCount={selected.size} priceCentsEach={show.price_cents} gstRatePercent={GST_RATE_PERCENT} />
          ) : (
            <div className="card" style={{ padding: 20, color: "var(--color-text-muted)", fontSize: 14 }}>
              Select a seat to see the price breakdown.
            </div>
          )}

          {error && <p className="errorText">{error}</p>}

          <button
            className={`btn btnPrimary ${styles.proceedBtn}`}
            disabled={selected.size === 0 || submitting}
            onClick={handleProceed}
          >
            {submitting ? <span className="spinner" /> : "Proceed to payment"}
          </button>
        </div>
      </div>
    </div>
  );
}
