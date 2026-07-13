import type { Candidate, Group, Member, Preference } from "@/types";
import { currentIdToken } from "@/lib/firebase";
import { loadMember } from "@/lib/member-store";

const BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

// 認証ヘッダー: FirebaseログインならBearerトークン、
// モックモード(Firebase未設定)なら従来のX-Member-ID。
async function authHeaders(): Promise<Record<string, string>> {
  const token = await currentIdToken();
  if (token) return { Authorization: `Bearer ${token}` };
  const member = loadMember();
  return member ? { "X-Member-ID": String(member.id) } : {};
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(await authHeaders()),
      ...init?.headers,
    },
  });
  if (!res.ok) {
    const body = await res.json().catch(() => null);
    throw new Error(body?.error ?? `API error: ${res.status}`);
  }
  if (res.status === 204) return undefined as T;
  return (await res.json()) as T;
}

export const api = {
  me: () => request<Member>("/api/me"),

  registerMember: (name: string) =>
    request<Member>("/api/members", {
      method: "POST",
      body: JSON.stringify({ name }),
    }),

  listGroups: () => request<Group[]>("/api/groups"),

  getGroup: (id: number) => request<Group>(`/api/groups/${id}`),

  createGroup: (name: string) =>
    request<Group>("/api/groups", {
      method: "POST",
      body: JSON.stringify({ name }),
    }),

  joinGroup: (inviteCode: string) =>
    request<Group>("/api/groups/join", {
      method: "POST",
      body: JSON.stringify({ invite_code: inviteCode }),
    }),

  listMyPreferences: () => request<Preference[]>("/api/me/preferences"),

  saveMyPreferences: (preferences: Preference[]) =>
    request<void>("/api/me/preferences", {
      method: "PUT",
      body: JSON.stringify({ preferences }),
    }),

  suggest: (groupId: number, area: string) =>
    request<Candidate[]>(
      `/api/groups/${groupId}/suggestions?area=${encodeURIComponent(area)}`,
    ),
};
