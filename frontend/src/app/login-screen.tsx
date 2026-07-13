"use client";

import { useState } from "react";

// Firebaseのエラーコードを利用者向けの日本語に変換する。
function errorMessage(code: string): string {
  if (code.includes("email-already-in-use"))
    return "このメールアドレスは登録済みです。ログインをお試しください。";
  if (code.includes("invalid-credential") || code.includes("wrong-password") || code.includes("user-not-found"))
    return "メールアドレスまたはパスワードが違います。";
  if (code.includes("weak-password"))
    return "パスワードは6文字以上にしてください。";
  if (code.includes("invalid-email")) return "メールアドレスの形式が正しくありません。";
  if (code.includes("too-many-requests"))
    return "試行回数が多すぎます。しばらく待ってからお試しください。";
  return "ログインに失敗しました。もう一度お試しください。";
}

// 認証画面(デザインモック準拠)。Firebaseモードでのみ表示される。
export function LoginScreen({
  onGoogleSignIn,
  onEmailSignIn,
  onEmailSignUp,
  onGuestSignIn,
}: {
  onGoogleSignIn: () => Promise<void>;
  onEmailSignIn: (email: string, password: string) => Promise<void>;
  onEmailSignUp: (email: string, password: string) => Promise<void>;
  onGuestSignIn: () => Promise<void>;
}) {
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const [showEmail, setShowEmail] = useState(false);
  const [mode, setMode] = useState<"signin" | "signup">("signin");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const signIn = async () => {
    setBusy(true);
    setError("");
    try {
      await onGoogleSignIn();
    } catch (e) {
      // ユーザーがポップアップを閉じただけならエラー表示しない
      const code = (e as { code?: string })?.code ?? "";
      if (!code.includes("popup-closed") && !code.includes("cancelled")) {
        setError(errorMessage(code));
      }
    } finally {
      setBusy(false);
    }
  };

  const signInGuest = async () => {
    setBusy(true);
    setError("");
    try {
      await onGuestSignIn();
    } catch (e) {
      setError(errorMessage((e as { code?: string })?.code ?? ""));
    } finally {
      setBusy(false);
    }
  };

  const submitEmail = async () => {
    if (!email.trim() || !password) return;
    setBusy(true);
    setError("");
    try {
      if (mode === "signup") {
        await onEmailSignUp(email.trim(), password);
      } else {
        await onEmailSignIn(email.trim(), password);
      }
    } catch (e) {
      setError(errorMessage((e as { code?: string })?.code ?? ""));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="login">
      <div className="login-hero">
        <div className="login-mark" aria-hidden="true">
          <svg width="44" height="44" viewBox="0 0 48 48" fill="none">
            <path
              d="M8 24 H40 C40 33 33 39.5 24 39.5 C15 39.5 8 33 8 24 Z"
              fill="#fff"
            />
            <path d="M18 40 H30 L28.5 43.5 H19.5 Z" fill="#fff" />
            <path
              d="M19 16.5 C17 13.5 21 11.5 19 8.5"
              stroke="#fff"
              strokeWidth="2.6"
              strokeLinecap="round"
            />
            <path
              d="M29 16.5 C27 13.5 31 11.5 29 8.5"
              stroke="#fff"
              strokeWidth="2.6"
              strokeLinecap="round"
            />
          </svg>
        </div>
        <h1 className="login-title">
          Food<span>Like</span>
        </h1>
        <p className="login-tagline">みんなの「食べたい」が、ひとつのお店になる。</p>
      </div>

      <div className="login-beats" aria-hidden="true">
        <span className="beat">好みを登録</span>
        <span className="beat">グループで共有</span>
        <span className="beat">妥協点のお店を提案</span>
      </div>

      <div className="login-actions">
        <button className="btn-google" onClick={signIn} disabled={busy}>
          <span className="g-badge" aria-hidden="true">
            <svg width="16" height="16" viewBox="0 0 18 18">
              <path
                fill="#4285F4"
                d="M17.64 9.2c0-.63-.06-1.25-.16-1.84H9v3.49h4.84a4.14 4.14 0 0 1-1.8 2.71v2.26h2.92a8.78 8.78 0 0 0 2.68-6.62z"
              />
              <path
                fill="#34A853"
                d="M9 18c2.43 0 4.47-.8 5.96-2.18l-2.92-2.26c-.8.54-1.84.86-3.04.86-2.34 0-4.32-1.58-5.03-3.71H.96v2.33A9 9 0 0 0 9 18z"
              />
              <path
                fill="#FBBC05"
                d="M3.97 10.71A5.41 5.41 0 0 1 3.68 9c0-.59.1-1.17.28-1.71V4.96H.96a9 9 0 0 0 0 8.08l3.01-2.33z"
              />
              <path
                fill="#EA4335"
                d="M9 3.58c1.32 0 2.5.45 3.44 1.35l2.58-2.59A9 9 0 0 0 .96 4.96l3.01 2.33C4.68 5.16 6.66 3.58 9 3.58z"
              />
            </svg>
          </span>
          {busy ? "ログイン中…" : "Googleでログイン"}
        </button>

        <div className="login-divider">または</div>

        {!showEmail ? (
          <button
            className="btn-email"
            onClick={() => setShowEmail(true)}
            disabled={busy}
          >
            メールアドレスで続ける
          </button>
        ) : (
          <div className="email-form">
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="メールアドレス"
              autoComplete="email"
            />
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && submitEmail()}
              placeholder="パスワード(6文字以上)"
              autoComplete={mode === "signup" ? "new-password" : "current-password"}
            />
            <button
              className="btn-email"
              onClick={submitEmail}
              disabled={busy || !email.trim() || !password}
            >
              {mode === "signup" ? "アカウントを作成" : "メールでログイン"}
            </button>
            <button
              className="link-btn"
              onClick={() => {
                setMode(mode === "signup" ? "signin" : "signup");
                setError("");
              }}
            >
              {mode === "signup"
                ? "アカウントをお持ちの方はログイン"
                : "はじめての方はアカウント作成"}
            </button>
          </div>
        )}
        <button className="link-btn" onClick={signInGuest} disabled={busy}>
          登録せずゲストとして試す
        </button>
        {error && <p className="error">{error}</p>}
      </div>

      <p className="login-legal">
        続行すると、利用規約・プライバシーポリシー・Cookieポリシーに同意したものとみなされます。
      </p>
    </div>
  );
}
