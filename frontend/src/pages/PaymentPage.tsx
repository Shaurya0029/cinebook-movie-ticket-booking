import { useEffect, useState, type FormEvent } from "react";
import { useNavigate, useParams } from "react-router-dom";
import * as bookingsApi from "../api/bookings";
import { formatMoney } from "../utils/format";
import type { Booking } from "../types";
import styles from "./PaymentPage.module.css";

export default function PaymentPage() {
  const { bookingId } = useParams();
  const navigate = useNavigate();

  const [booking, setBooking] = useState<Booking | null>(null);
  const [cardNumber, setCardNumber] = useState("");
  const [cardName, setCardName] = useState("");
  const [expiry, setExpiry] = useState("");
  const [cvv, setCvv] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!bookingId) return;
    bookingsApi
      .getBooking(Number(bookingId))
      .then((b) => {
        if (b.status !== "pending_payment") {
          navigate(`/bookings/${b.id}/receipt`, { replace: true });
          return;
        }
        setBooking(b);
      })
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load booking"));
  }, [bookingId, navigate]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!booking) return;
    setSubmitting(true);
    setError(null);
    try {
      await bookingsApi.payForBooking(booking.id, {
        card_number: cardNumber,
        card_name: cardName,
        expiry,
        cvv,
      });
      navigate(`/bookings/${booking.id}/receipt`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Payment failed");
      setSubmitting(false);
    }
  }

  if (error && !booking) {
    return (
      <div className="page container">
        <p className="errorText">{error}</p>
      </div>
    );
  }

  if (!booking) {
    return (
      <div className="page container">
        <div className="skeleton" style={{ height: 360, borderRadius: 16 }} />
      </div>
    );
  }

  return (
    <div className="page container fadeIn">
      <div className={styles.layout}>
        <div className={`card ${styles.formCard}`}>
          <h1 className={styles.title}>Payment details</h1>
          <p className={styles.subtitle}>This is a simulated payment — no real card is charged.</p>

          <div className={styles.hint}>
            Any card number succeeds, except the test decline card <code>4000 0000 0000 0002</code>, which always
            fails so you can see what a declined payment looks like.
          </div>

          <form onSubmit={handleSubmit}>
            <div className="formField">
              <label htmlFor="cardName">Name on card</label>
              <input id="cardName" value={cardName} onChange={(e) => setCardName(e.target.value)} required autoFocus />
            </div>
            <div className="formField">
              <label htmlFor="cardNumber">Card number</label>
              <input
                id="cardNumber"
                inputMode="numeric"
                placeholder="1234 5678 9012 3456"
                value={cardNumber}
                onChange={(e) => setCardNumber(formatCardNumber(e.target.value))}
                maxLength={19}
                required
              />
            </div>
            <div className={styles.summaryRow} style={{ display: "flex", gap: 12 }}>
              <div className="formField" style={{ flex: 1 }}>
                <label htmlFor="expiry">Expiry (MM/YY)</label>
                <input
                  id="expiry"
                  placeholder="12/28"
                  value={expiry}
                  onChange={(e) => setExpiry(formatExpiry(e.target.value))}
                  maxLength={5}
                  required
                />
              </div>
              <div className="formField" style={{ flex: 1 }}>
                <label htmlFor="cvv">CVV</label>
                <input
                  id="cvv"
                  inputMode="numeric"
                  placeholder="123"
                  value={cvv}
                  onChange={(e) => setCvv(e.target.value.replace(/\D/g, "").slice(0, 3))}
                  maxLength={3}
                  required
                />
              </div>
            </div>

            {error && <p className="errorText">{error}</p>}

            <button type="submit" className={`btn btnPrimary ${styles.payBtn}`} disabled={submitting}>
              {submitting ? <span className="spinner" /> : `Pay ${formatMoney(booking.total_amount_cents)}`}
            </button>
            <div className={styles.securePill}>
              <span className={styles.dot} />
              Simulated secure checkout
            </div>
          </form>
        </div>

        <div className={`card ${styles.summary}`}>
          <div className={styles.summaryRow}>
            <span>Seats</span>
            <span>{booking.seat_count}</span>
          </div>
          <div className={styles.summaryRow}>
            <span>Subtotal</span>
            <span>{formatMoney(booking.subtotal_cents)}</span>
          </div>
          <div className={styles.summaryRow}>
            <span>GST ({booking.gst_rate_percent}%)</span>
            <span>{formatMoney(booking.gst_amount_cents)}</span>
          </div>
          <div className={styles.divider} />
          <div className={styles.total}>
            <span>Total</span>
            <span>{formatMoney(booking.total_amount_cents)}</span>
          </div>
        </div>
      </div>
    </div>
  );
}

function formatCardNumber(raw: string): string {
  const digits = raw.replace(/\D/g, "").slice(0, 16);
  return digits.replace(/(.{4})/g, "$1 ").trim();
}

function formatExpiry(raw: string): string {
  const digits = raw.replace(/\D/g, "").slice(0, 4);
  if (digits.length <= 2) return digits;
  return `${digits.slice(0, 2)}/${digits.slice(2)}`;
}
