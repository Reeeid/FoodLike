import { getApps, initializeApp } from "firebase/app";
import { getAuth, type Auth } from "firebase/auth";

// Firebase Authの初期化。NEXT_PUBLIC_FIREBASE_API_KEYが未設定なら
// Firebase無効(モック認証モード)として扱い、初期化しない。
// 設定値はコンソールの「プロジェクトの設定 > マイアプリ」からコピーする。

const apiKey = process.env.NEXT_PUBLIC_FIREBASE_API_KEY;

export const firebaseEnabled = Boolean(apiKey);

export function getFirebaseAuth(): Auth {
  if (!apiKey) {
    throw new Error("Firebase is not configured (NEXT_PUBLIC_FIREBASE_API_KEY)");
  }
  const app =
    getApps()[0] ??
    initializeApp({
      apiKey,
      authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
      projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
      appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
    });
  return getAuth(app);
}

// APIリクエストに載せるIDトークン。未ログイン・Firebase無効時はnull。
// getIdToken()は期限切れが近ければ自動でリフレッシュしてくれる。
export async function currentIdToken(): Promise<string | null> {
  if (!firebaseEnabled || typeof window === "undefined") return null;
  const user = getFirebaseAuth().currentUser;
  return user ? user.getIdToken() : null;
}
