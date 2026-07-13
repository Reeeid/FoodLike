import type { Member } from "@/types";

// MVPの仮認証: 登録したメンバーをlocalStorageに保持する外部ストア。
// useSyncExternalStoreで購読できる形にしてある。
// Firebase Auth導入(issue #8)時にこのモジュールごと差し替える。

const KEY = "foodlike:member";

let cache: { raw: string | null; member: Member | null } = {
  raw: null,
  member: null,
};
const listeners = new Set<() => void>();

function emit() {
  listeners.forEach((l) => l());
}

export function subscribeMember(cb: () => void): () => void {
  listeners.add(cb);
  window.addEventListener("storage", cb);
  return () => {
    listeners.delete(cb);
    window.removeEventListener("storage", cb);
  };
}

export function getMemberSnapshot(): Member | null {
  const raw = localStorage.getItem(KEY);
  if (raw !== cache.raw) {
    cache = { raw, member: raw ? (JSON.parse(raw) as Member) : null };
  }
  return cache.member;
}

// SSR/プリレンダリング時はログイン状態を持たない。
export function getMemberServerSnapshot(): Member | null {
  return null;
}

export function loadMember(): Member | null {
  if (typeof window === "undefined") return null;
  return getMemberSnapshot();
}

export function saveMember(member: Member) {
  localStorage.setItem(KEY, JSON.stringify(member));
  emit();
}

export function clearMember() {
  localStorage.removeItem(KEY);
  emit();
}
