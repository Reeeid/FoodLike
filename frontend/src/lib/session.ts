"use client";

import { useCallback, useEffect, useState, useSyncExternalStore } from "react";
import {
  GoogleAuthProvider,
  createUserWithEmailAndPassword,
  onAuthStateChanged,
  signInAnonymously,
  signInWithEmailAndPassword,
  signInWithPopup,
  signOut as firebaseSignOut,
} from "firebase/auth";

import { api } from "@/lib/api";
import { firebaseEnabled, getFirebaseAuth } from "@/lib/firebase";
import {
  clearMember,
  getMemberServerSnapshot,
  getMemberSnapshot,
  subscribeMember,
} from "@/lib/member-store";
import type { Member } from "@/types";

// ログイン状態を一本化するフック。
// - Firebaseモード: onAuthStateChangedを購読し、ログイン後に GET /api/me で
//   内部メンバーを解決する(初回はバックエンドがJITで作成)。
// - モックモード(Firebase未設定): 従来どおりmember-store(localStorage)を使う。

export type SessionStatus = "loading" | "guest" | "authed";

export type Session = {
  status: SessionStatus;
  member: Member | null;
  firebaseEnabled: boolean;
  signInWithGoogle: () => Promise<void>;
  signInWithEmail: (email: string, password: string) => Promise<void>;
  signUpWithEmail: (email: string, password: string) => Promise<void>;
  signInAsGuest: () => Promise<void>;
  signOut: () => Promise<void>;
};

export function useSession(): Session {
  // モックモード用(Firebaseモードでは常にnullのまま)
  const mockMember = useSyncExternalStore(
    subscribeMember,
    getMemberSnapshot,
    getMemberServerSnapshot,
  );

  // Firebaseモード用
  const [fb, setFb] = useState<{ status: SessionStatus; member: Member | null }>({
    status: "loading",
    member: null,
  });

  useEffect(() => {
    if (!firebaseEnabled) return;
    return onAuthStateChanged(getFirebaseAuth(), async (user) => {
      if (!user) {
        setFb({ status: "guest", member: null });
        return;
      }
      try {
        setFb({ status: "authed", member: await api.me() });
      } catch {
        // トークンは有効だがAPIに届かない(バックエンド停止等)。ゲスト扱いに戻す。
        setFb({ status: "guest", member: null });
      }
    });
  }, []);

  const signInWithGoogle = useCallback(async () => {
    await signInWithPopup(getFirebaseAuth(), new GoogleAuthProvider());
    // 後続はonAuthStateChangedが拾う
  }, []);

  const signInWithEmail = useCallback(async (email: string, password: string) => {
    await signInWithEmailAndPassword(getFirebaseAuth(), email, password);
  }, []);

  const signUpWithEmail = useCallback(async (email: string, password: string) => {
    await createUserWithEmailAndPassword(getFirebaseAuth(), email, password);
  }, []);

  // 匿名ログイン。UIDは端末ごとに発行され、ログアウトすると二度と戻れない点に注意。
  const signInAsGuest = useCallback(async () => {
    await signInAnonymously(getFirebaseAuth());
  }, []);

  const signOut = useCallback(async () => {
    if (firebaseEnabled) {
      await firebaseSignOut(getFirebaseAuth());
    } else {
      clearMember();
    }
  }, []);

  if (firebaseEnabled) {
    return {
      ...fb,
      firebaseEnabled,
      signInWithGoogle,
      signInWithEmail,
      signUpWithEmail,
      signInAsGuest,
      signOut,
    };
  }
  return {
    status: mockMember ? "authed" : "guest",
    member: mockMember,
    firebaseEnabled,
    signInWithGoogle,
    signInWithEmail,
    signUpWithEmail,
    signInAsGuest,
    signOut,
  };
}
