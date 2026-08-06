import type { Seat } from "../types";
import styles from "./SeatMap.module.css";

interface Props {
  seats: Seat[];
  selectedIds: Set<number>;
  onToggle: (seat: Seat) => void;
}

export default function SeatMap({ seats, selectedIds, onToggle }: Props) {
  const rows = groupByRow(seats);

  return (
    <div className={styles.wrap}>
      <div className={styles.screen}>
        <div className={styles.screenBar} />
        <div className={styles.screenLabel}>SCREEN</div>
      </div>

      <div className={styles.grid}>
        {rows.map(([rowLabel, rowSeats]) => (
          <div key={rowLabel} className={styles.row}>
            <span className={styles.rowLabel}>{rowLabel}</span>
            <div className={styles.seats}>
              {rowSeats.map((seat) => {
                const isSelected = selectedIds.has(seat.id);
                const classes = [styles.seat, isSelected ? styles.selected : "", seat.is_booked ? styles.booked : ""]
                  .filter(Boolean)
                  .join(" ");
                return (
                  <button
                    key={seat.id}
                    type="button"
                    className={classes}
                    disabled={seat.is_booked}
                    onClick={() => onToggle(seat)}
                    title={seat.seat_label}
                  >
                    {seat.seat_label.slice(1)}
                  </button>
                );
              })}
            </div>
          </div>
        ))}
      </div>

      <div className={styles.legend}>
        <div className={styles.legendItem}>
          <span className={styles.legendSwatch} style={{ background: "var(--color-seat-available)" }} />
          Available
        </div>
        <div className={styles.legendItem}>
          <span className={styles.legendSwatch} style={{ background: "var(--color-seat-selected)" }} />
          Selected
        </div>
        <div className={styles.legendItem}>
          <span className={styles.legendSwatch} style={{ background: "var(--color-seat-booked-bg)" }} />
          Booked
        </div>
      </div>
    </div>
  );
}

function groupByRow(seats: Seat[]): [string, Seat[]][] {
  const map = new Map<string, Seat[]>();
  for (const seat of seats) {
    const row = seat.seat_label[0];
    if (!map.has(row)) map.set(row, []);
    map.get(row)!.push(seat);
  }
  for (const rowSeats of map.values()) {
    rowSeats.sort((a, b) => parseInt(a.seat_label.slice(1), 10) - parseInt(b.seat_label.slice(1), 10));
  }
  return Array.from(map.entries()).sort((a, b) => a[0].localeCompare(b[0]));
}
