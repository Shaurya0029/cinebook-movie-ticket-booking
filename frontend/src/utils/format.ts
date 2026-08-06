export function formatMoney(cents: number): string {
  return `₹${(cents / 100).toFixed(2)}`;
}

export function formatShowDateTime(iso: string): string {
  const date = new Date(iso);
  return date.toLocaleString("en-IN", {
    weekday: "short",
    day: "2-digit",
    month: "short",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function formatShowDate(iso: string): string {
  const date = new Date(iso);
  return date.toLocaleDateString("en-IN", { weekday: "short", day: "2-digit", month: "short" });
}

export function formatShowTime(iso: string): string {
  const date = new Date(iso);
  return date.toLocaleTimeString("en-IN", { hour: "2-digit", minute: "2-digit" });
}
