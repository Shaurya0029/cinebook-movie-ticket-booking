import { BrowserRouter, Route, Routes } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import { LocationProvider } from "./context/LocationContext";
import Navbar from "./components/Navbar";
import RequireAuthRoute from "./components/RequireAuthRoute";
import HomePage from "./pages/HomePage";
import CityTheaterPickerPage from "./pages/CityTheaterPickerPage";
import MovieDetailsPage from "./pages/MovieDetailsPage";
import SeatSelectionPage from "./pages/SeatSelectionPage";
import PaymentPage from "./pages/PaymentPage";
import ReceiptPage from "./pages/ReceiptPage";
import MyBookingsPage from "./pages/MyBookingsPage";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <LocationProvider>
          <Navbar />
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/select-location" element={<CityTheaterPickerPage />} />
            <Route path="/movies/:movieId" element={<MovieDetailsPage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route
              path="/shows/:showId/seats"
              element={
                <RequireAuthRoute>
                  <SeatSelectionPage />
                </RequireAuthRoute>
              }
            />
            <Route
              path="/bookings/:bookingId/pay"
              element={
                <RequireAuthRoute>
                  <PaymentPage />
                </RequireAuthRoute>
              }
            />
            <Route
              path="/bookings/:bookingId/receipt"
              element={
                <RequireAuthRoute>
                  <ReceiptPage />
                </RequireAuthRoute>
              }
            />
            <Route
              path="/my-bookings"
              element={
                <RequireAuthRoute>
                  <MyBookingsPage />
                </RequireAuthRoute>
              }
            />
          </Routes>
        </LocationProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}
