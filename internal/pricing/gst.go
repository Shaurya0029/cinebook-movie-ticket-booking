package pricing

import "math"

// Breakdown holds the computed price breakdown for a booking, all amounts in cents (paise).
type Breakdown struct {
	SubtotalCents    int
	GSTRatePercent   float64
	GSTAmountCents   int
	TotalAmountCents int
}

// Calculate computes subtotal/GST/total for seatCount seats at priceCents each.
func Calculate(priceCents, seatCount int, gstRatePercent float64) Breakdown {
	subtotal := priceCents * seatCount
	gst := int(math.Round(float64(subtotal) * gstRatePercent / 100))

	return Breakdown{
		SubtotalCents:    subtotal,
		GSTRatePercent:   gstRatePercent,
		GSTAmountCents:   gst,
		TotalAmountCents: subtotal + gst,
	}
}
