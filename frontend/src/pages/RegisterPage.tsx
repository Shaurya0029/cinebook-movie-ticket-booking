import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import styles from "./AuthForm.module.css";

export default function RegisterPage() {
  const { register } = useAuth();
  const navigate = useNavigate();

  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      await register(firstName, lastName, email, password);
      navigate("/", { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="page container">
      <div className={`card ${styles.wrap} fadeInUp`}>
        <h1 className={styles.title}>Create your account</h1>
        <p className={styles.subtitle}>Sign up to start booking tickets.</p>

        <form onSubmit={handleSubmit}>
          <div className={styles.row}>
            <div className="formField">
              <label htmlFor="firstName">First name</label>
              <input id="firstName" value={firstName} onChange={(e) => setFirstName(e.target.value)} required autoFocus />
            </div>
            <div className="formField">
              <label htmlFor="lastName">Last name</label>
              <input id="lastName" value={lastName} onChange={(e) => setLastName(e.target.value)} required />
            </div>
          </div>
          <div className="formField">
            <label htmlFor="email">Email</label>
            <input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
          </div>
          <div className="formField">
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              minLength={8}
              required
            />
          </div>

          {error && <p className="errorText">{error}</p>}

          <button type="submit" className={`btn btnPrimary ${styles.submit}`} disabled={submitting}>
            {submitting ? <span className="spinner" /> : "Sign up"}
          </button>
        </form>

        <p className={styles.footer}>
          Already have an account? <Link to="/login">Log in</Link>
        </p>
      </div>
    </div>
  );
}
