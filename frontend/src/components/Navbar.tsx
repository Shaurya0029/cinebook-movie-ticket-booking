import { Link, useLocation as useRouterLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useLocation as useCinemaLocation } from "../context/LocationContext";
import styles from "./Navbar.module.css";

export default function Navbar() {
  const { user, logout } = useAuth();
  const { selectedCity, selectedTheater } = useCinemaLocation();
  const navigate = useNavigate();
  const routerLocation = useRouterLocation();

  function handleLogout() {
    logout();
    navigate("/");
  }

  function handleBack() {
    if (window.history.length > 1) {
      navigate(-1);
    } else {
      navigate("/");
    }
  }

  const showBack = routerLocation.pathname !== "/";

  return (
    <header className={styles.header}>
      <div className={`container ${styles.bar}`}>
        <div className={styles.leftGroup}>
          {showBack && (
            <button className={styles.backButton} onClick={handleBack} aria-label="Go back" title="Go back">
              ←
            </button>
          )}
          <Link to="/" className={styles.brand}>
            Cine<span>Book</span>
          </Link>
        </div>

        <nav className={styles.nav}>
          <button className={styles.locationButton} onClick={() => navigate("/select-location")}>
            <span className={styles.locationLabel}>Location</span>
            <span className={styles.locationValue}>
              {selectedTheater ? selectedTheater.name : selectedCity ? selectedCity.name : "Select"}
            </span>
          </button>

          {user ? (
            <div className={styles.userMenu}>
              <Link to="/my-bookings" className={styles.navLink}>
                My Bookings
              </Link>
              <div className={styles.avatar}>{user.first_name[0]?.toUpperCase()}</div>
              <button className="btn btnSecondary" onClick={handleLogout}>
                Log out
              </button>
            </div>
          ) : (
            <div className={styles.userMenu}>
              <Link to="/login" className={styles.navLink}>
                Log in
              </Link>
              <Link to="/register" className="btn btnPrimary">
                Sign up
              </Link>
            </div>
          )}
        </nav>
      </div>
    </header>
  );
}
