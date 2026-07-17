import type {
  Candidate,
  Group,
  Member,
  Message,
  Preference,
  SearchCriteria,
  SearchOptions,
} from "@/types";
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

  getSearchOptions: () => request<SearchOptions>("/api/search-options"),

  suggest: (groupId: number, criteria: SearchCriteria) => {
    const q = new URLSearchParams();
    if (criteria.large_area) q.set("large_area", criteria.large_area);
    if (criteria.middle_area) q.set("middle_area", criteria.middle_area);
    if (criteria.small_area) q.set("small_area", criteria.small_area);
    if (criteria.budget) q.set("budget", criteria.budget);
    return request<Candidate[]>(
      `/api/groups/${groupId}/suggestions?${q.toString()}`,
    );
  },

  listMessages: (groupId: number, afterId = 0) =>
    request<Message[]>(`/api/groups/${groupId}/messages?after_id=${afterId}`),

  postMessage: (groupId: number, text: string) =>
    request<Message>(`/api/groups/${groupId}/messages`, {
      method: "POST",
      body: JSON.stringify({ text }),
    }),

  // AI検索(SSE)。EventSourceは認証ヘッダーを積めないのでfetchのストリームで読む。
  // chunkイベントをonChunkに流し、doneで保存済みの(質問, 回答)を返す。
  aiSearch: async (
    groupId: number,
    q: string,
    onChunk: (text: string) => void,
  ): Promise<{ question: Message; answer: Message }> => {
    const res = await fetch(
      `${BASE_URL}/api/groups/${groupId}/ai-search?q=${encodeURIComponent(q)}`,
      { headers: await authHeaders() },
    );
    if (!res.ok || !res.body) {
      throw new Error(`AI検索に失敗しました (${res.status})`);
    }
    const reader = res.body.getReader();
    const decoder = new TextDecoder();
    let buf = "";
    let result: { question: Message; answer: Message } | null = null;

    const handleEvent = (raw: string) => {
      let event = "";
      let data = "";
      for (const line of raw.split("\n")) {
        if (line.startsWith("event: ")) event = line.slice(7).trim();
        else if (line.startsWith("data: ")) data += line.slice(6);
      }
      if (!event || !data) return;
      if (event === "chunk") onChunk(JSON.parse(data) as string);
      else if (event === "done")
        result = JSON.parse(data) as { question: Message; answer: Message };
      else if (event === "error")
        throw new Error((JSON.parse(data) as { error: string }).error);
    };

    for (;;) {
      const { done, value } = await reader.read();
      if (done) break;
      buf += decoder.decode(value, { stream: true });
      let idx;
      while ((idx = buf.indexOf("\n\n")) >= 0) {
        handleEvent(buf.slice(0, idx));
        buf = buf.slice(idx + 2);
      }
    }
    if (!result) throw new Error("AI検索が中断されました");
    return result;
  },
};
