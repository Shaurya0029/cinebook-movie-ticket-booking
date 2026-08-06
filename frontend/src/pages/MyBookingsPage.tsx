import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import * as bookingsApi from "../api/bookings";
import { formatMoney, formatShowDateTime } from "../utils/format";
import type { Receipt } from "../types";
import styles from "./MyBookingsPage.module.css";

export default function MyBookingsPage() {
  const [receipts, setReceipts] = useState<Receipt[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    bookingsApi
      .listMyBookings()
      .then(async (bookings) => {
        const sorted = [...bookings].sort((a, b) => b.id - a.id);
        const details = await Promise.all(sorted.map((b) => bookingsApi.getBookingReceipt(b.id)));
        setReceipts(details);
      })
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load bookings"));
  }, []);

  return (
    <div className="page container fadeIn">
      <h1 className="sectionTitle" style={{ marginTop: 8 }}>
        My Bookings
      </h1>

      {error && <p className="errorText">{error}</p>}

      {!receipts ? (
        <div className={styles.list}>
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="skeleton" style={{ height: 96, borderRadius: 16 }} />
          ))}
        </div>
      ) : receipts.length === 0 ? (
        <div className={`card ${styles.empty}`}>
          <p>You haven't booked any tickets yet.</p>
          <Link to="/" className="btn btnPrimary" style={{ marginTop: 12 }}>
            Browse movies
          </Link>
        </div>
      ) : (
        <div className={styles.list}>
          {receipts.map((r, i) => (
            <Link
              key={r.booking.id}
              to={`/bookings/${r.booking.id}/receipt`}
              className={`card ${styles.row} fadeInUp`}
              style={{ animationDelay: `${Math.min(i, 10) * 30}ms` }}
            >
              <div className={styles.info}>
                <div className={styles.movieTitle}>{r.movie_title}</div>
                <div className={styles.meta}>
                  {r.theater_name} · {r.city_name} · {formatShowDateTime(r.show_starts_at)}
                </div>
                <div className={styles.meta}>{(r.seat_labels ?? []).join(", ")}</div>
              </div>
              <div className={styles.right}>
                <span className={styles.amount}>{formatMoney(r.booking.total_amount_cents)}</span>
                <span className={`${styles.statusBadge} ${styles[r.booking.status]}`}>
                  {r.booking.status.replace("_", " ")}
                </span>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
