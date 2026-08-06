import { apiRequest, apiRequestAllowStatus } from "./client";
import type { Booking, PayResponse, Receipt } from "../types";

export function createBooking(showId: number, seatIds: number[]) {
  return apiRequest<Booking>("/bookings", {
    method: "POST",
    body: { show_id: showId, seat_ids: seatIds },
  });
}

export function listMyBookings() {
  return apiRequest<Booking[]>("/bookings");
}

export function getBooking(bookingId: number) {
  return apiRequest<Booking>(`/bookings/${bookingId}`);
}

export function getBookingReceipt(bookingId: number) {
  return apiRequest<Receipt>(`/bookings/${bookingId}/receipt`);
}

export interface PayCardDetails {
  card_number: string;
  card_name: string;
  expiry: string;
  cvv: string;
}

// 402 is a normal, expected outcome here (the simulated card was declined),
// not a transport error, so it's allowed through alongside 200.
export function payForBooking(bookingId: number, card: PayCardDetails) {
  return apiRequestAllowStatus<PayResponse>(`/bookings/${bookingId}/pay`, [200, 402], {
    method: "POST",
    body: card,
  });
}
