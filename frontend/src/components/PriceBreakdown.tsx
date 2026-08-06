import { formatMoney } from "../utils/format";
import styles from "./PriceBreakdown.module.css";

interface Props {
  seatCount: number;
  priceCentsEach: number;
  gstRatePercent: number;
}

export default function PriceBreakdown({ seatCount, priceCentsEach, gstRatePercent }: Props) {
  const subtotal = seatCount * priceCentsEach;
  const gst = Math.round((subtotal * gstRatePercent) / 100);
  const total = subtotal + gst;

  return (
    <div className={`card ${styles.wrap}`}>
      <div className={styles.row}>
        <span>
          Tickets ({seatCount} x {formatMoney(priceCentsEach)})
        </span>
        <span key={subtotal} className="fadeIn">
          {formatMoney(subtotal)}
        </span>
      </div>
      <div className={styles.row}>
        <span>GST ({gstRatePercent}%)</span>
        <span key={gst} className="fadeIn">
          {formatMoney(gst)}
        </span>
      </div>
      <div className={styles.divider} />
      <div className={styles.total}>
        <span>Total</span>
        <span key={total} className="fadeIn">
          {formatMoney(total)}
        </span>
      </div>
    </div>
  );
}
