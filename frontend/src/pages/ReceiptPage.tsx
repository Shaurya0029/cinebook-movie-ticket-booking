import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import * as bookingsApi from "../api/bookings";
import { formatMoney, formatShowDateTime } from "../utils/format";
import type { Receipt } from "../types";
import styles from "./ReceiptPage.module.css";

export default function ReceiptPage() {
  const { bookingId } = useParams();
  const [receipt, setReceipt] = useState<Receipt | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!bookingId) return;
    bookingsApi
      .getBookingReceipt(Number(bookingId))
      .then(setReceipt)
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load receipt"));
  }, [bookingId]);

  if (error) {
    return (
      <div className="page container">
        <p className="errorText">{error}</p>
      </div>
    );
  }

  if (!receipt) {
    return (
      <div className="page container">
        <div className="skeleton" style={{ height: 480, maxWidth: 480, margin: "0 auto", borderRadius: 16 }} />
      </div>
    );
  }

  const succeeded = receipt.payment?.status === "success";

  return (
    <div className="page container fadeIn">
      <div className={`card ${styles.wrap} fadeInUp`}>
        <div className={`${styles.statusBand} ${succeeded ? styles.success : styles.failed}`}>
          <div className={`${styles.iconCircle} ${succeeded ? styles.success : styles.failed}`}>
            {succeeded ? "✓" : "✕"}
          </div>
          <h1 className={styles.statusTitle}>{succeeded ? "Booking confirmed!" : "Payment declined"}</h1>
          <p className={styles.statusDesc}>
            {succeeded
              ? `Thanks, ${receipt.user_first_name}. Your tickets are booked.`
              : receipt.payment?.failure_reason ?? "Your payment could not be processed. Your seats have been released."}
          </p>
        </div>

        <div className={styles.ticket}>
          <h2 className={styles.movieTitle}>{receipt.movie_title}</h2>

          <div className={styles.detailGrid}>
            <div className={styles.detailItem}>
              <span className={styles.detailLabel}>Booked by</span>
              <span className={styles.detailValue}>
                {receipt.user_first_name} {receipt.user_last_name}
              </span>
            </div>
            <div className={styles.detailItem}>
              <span className={styles.detailLabel}>Theater</span>
              <span className={styles.detailValue}>{receipt.theater_name}</span>
            </div>
            <div className={styles.detailItem}>
              <span className={styles.detailLabel}>City</span>
              <span className={styles.detailValue}>{receipt.city_name}</span>
            </div>
            <div className={styles.detailItem}>
              <span className={styles.detailLabel}>Showtime</span>
              <span className={styles.detailValue}>{formatShowDateTime(receipt.show_starts_at)}</span>
            </div>
            <div className={styles.detailItem}>
              <span className={styles.detailLabel}>Seats</span>
              <span className={styles.detailValue}>{(receipt.seat_labels ?? []).join(", ")}</span>
            </div>
            <div className={styles.detailItem}>
              <span className={styles.detailLabel}>Price / ticket</span>
              <span className={styles.detailValue}>{formatMoney(receipt.price_cents)}</span>
            </div>
          </div>

          <div className={styles.divider} />

          <div className={styles.amountRow}>
            <span>Subtotal</span>
            <span>{formatMoney(receipt.booking.subtotal_cents)}</span>
          </div>
          <div className={styles.amountRow}>
            <span>GST ({receipt.booking.gst_rate_percent}%)</span>
            <span>{formatMoney(receipt.booking.gst_amount_cents)}</span>
          </div>
          <div className={styles.totalRow}>
            <span>{succeeded ? "Amount paid" : "Amount"}</span>
            <span>{formatMoney(receipt.booking.total_amount_cents)}</span>
          </div>

          {receipt.payment && (
            <div className={styles.txnId}>
              Transaction #{receipt.payment.id} · Card ending {receipt.payment.card_last4}
            </div>
          )}
        </div>

        <div className={styles.actions}>
          {succeeded ? (
            <>
              <Link to="/my-bookings" className="btn btnSecondary">
                My bookings
              </Link>
              <Link to="/" className="btn btnPrimary">
                Book more
              </Link>
            </>
          ) : (
            <>
              <Link to="/my-bookings" className="btn btnSecondary">
                My bookings
              </Link>
              <Link to="/" className="btn btnPrimary">
                Try again
              </Link>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
