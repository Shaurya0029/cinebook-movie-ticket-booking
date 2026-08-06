import { useState, type FormEvent } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import styles from "./AuthForm.module.css";

export default function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      await login(email, password);
      const from = (location.state as { from?: string } | null)?.from ?? "/";
      navigate(from, { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="page container">
      <div className={`card ${styles.wrap} fadeInUp`}>
        <h1 className={styles.title}>Welcome back</h1>
        <p className={styles.subtitle}>Log in to book tickets and see your bookings.</p>

        <form onSubmit={handleSubmit}>
          <div className="formField">
            <label htmlFor="email">Email</label>
            <input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required autoFocus />
          </div>
          <div className="formField">
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          {error && <p className="errorText">{error}</p>}

          <button type="submit" className={`btn btnPrimary ${styles.submit}`} disabled={submitting}>
            {submitting ? <span className="spinner" /> : "Log in"}
          </button>
        </form>

        <p className={styles.footer}>
          New here? <Link to="/register">Create an account</Link>
        </p>
      </div>
    </div>
  );
}
