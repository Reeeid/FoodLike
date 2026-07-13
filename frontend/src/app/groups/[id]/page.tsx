"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { api } from "@/lib/api";
import { useSession } from "@/lib/session";
import { InviteModal } from "@/components/invite-modal";
import type {
  Candidate,
  Group,
  Preference,
  PreferenceCategory,
} from "@/types";

export default function GroupPage() {
  const params = useParams<{ id: string }>();
  const groupId = Number(params.id);
  const [group, setGroup] = useState<Group | null>(null);
  const [error, setError] = useState("");
  const [inviteOpen, setInviteOpen] = useState(false);

  const session = useSession();

  useEffect(() => {
    if (session.status === "guest") {
      window.location.href = "/";
      return;
    }
    if (session.status !== "authed") return;
    api
      .getGroup(groupId)
      .then(setGroup)
      .catch((e) =>
        setError(e instanceof Error ? e.message : "グループの取得に失敗しました"),
      );
  }, [groupId, session.status]);

  if (error) {
    return (
      <>
        <Link href="/" className="back-link">
          ← ホームへ
        </Link>
        <p className="error">{error}</p>
      </>
    );
  }
  if (!group) return null;

  return (
    <>
      <Link href="/" className="back-link">
        ← ホームへ
      </Link>

      <div className="card">
        <h2 className="card-title">{group.name}</h2>
        <div className="chips">
          {group.members.map((m) => (
            <span key={m.id} className="chip">
              👤 {m.name}
            </span>
          ))}
        </div>
      </div>

      <div className="card">
        <h2 className="card-title">友だちを招待</h2>
        <p className="section-note">
          QRコードや招待コードを共有すると参加できます。
        </p>
        <button className="btn" onClick={() => setInviteOpen(true)}>
          招待コード・QRを表示
        </button>
      </div>
      {inviteOpen && (
        <InviteModal group={group} onClose={() => setInviteOpen(false)} />
      )}

      <PreferenceCard />
      <SuggestionCard groupId={groupId} />
    </>
  );
}

function PreferenceCard() {
  const [prefs, setPrefs] = useState<Preference[]>([]);
  const [value, setValue] = useState("");
  const [category, setCategory] = useState<PreferenceCategory>("genre");
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    api
      .listMyPreferences()
      .then(setPrefs)
      .catch((e) =>
        setError(e instanceof Error ? e.message : "取得に失敗しました"),
      );
  }, []);

  const save = useCallback(async (next: Preference[]) => {
    setSaving(true);
    setError("");
    try {
      await api.saveMyPreferences(next);
      setPrefs(next);
    } catch (e) {
      setError(e instanceof Error ? e.message : "保存に失敗しました");
    } finally {
      setSaving(false);
    }
  }, []);

  const add = () => {
    const v = value.trim();
    if (!v) return;
    if (prefs.some((p) => p.value === v && p.category === category)) {
      setValue("");
      return;
    }
    save([...prefs, { kind: "dislike", category, value: v }]);
    setValue("");
  };

  const remove = (target: Preference) => {
    save(
      prefs.filter(
        (p) => !(p.value === target.value && p.category === target.category),
      ),
    );
  };

  return (
    <div className="card">
      <h2 className="card-title">🤫 あなたの苦手なもの</h2>
      <p className="section-note">
        ここに登録した内容は他のメンバーには一切見えません。提案の絞り込みにだけ使われます。
      </p>
      <div className="field-row">
        <select
          value={category}
          onChange={(e) => setCategory(e.target.value as PreferenceCategory)}
          style={{ maxWidth: "8rem" }}
        >
          <option value="genre">ジャンル</option>
          <option value="ingredient">食材</option>
        </select>
        <input
          value={value}
          onChange={(e) => setValue(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && add()}
          placeholder={category === "genre" ? "例: 辛い物" : "例: エビ"}
          maxLength={64}
        />
        <button className="btn" onClick={add} disabled={saving || !value.trim()}>
          追加
        </button>
      </div>
      {prefs.length > 0 && (
        <div className="chips" style={{ marginTop: "0.75rem" }}>
          {prefs.map((p) => (
            <span key={`${p.category}:${p.value}`} className="chip">
              {p.category === "genre" ? "🍽" : "🥕"} {p.value}
              <button onClick={() => remove(p)} aria-label="削除">
                ×
              </button>
            </span>
          ))}
        </div>
      )}
      {error && <p className="error">{error}</p>}
    </div>
  );
}

function SuggestionCard({ groupId }: { groupId: number }) {
  const [area, setArea] = useState("");
  const [candidates, setCandidates] = useState<Candidate[] | null>(null);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  const suggest = async () => {
    setBusy(true);
    setError("");
    try {
      setCandidates(await api.suggest(groupId, area.trim()));
    } catch (e) {
      setError(e instanceof Error ? e.message : "提案の取得に失敗しました");
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="card">
      <h2 className="card-title">🍜 お店を提案してもらう</h2>
      <p className="section-note">
        メンバー全員の苦手なものにこっそり配慮したお店を探します。
      </p>
      <div className="field-row">
        <input
          value={area}
          onChange={(e) => setArea(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && suggest()}
          placeholder="エリア(例: 新宿)※空欄で全エリア"
        />
        <button className="btn" onClick={suggest} disabled={busy}>
          {busy ? "検索中…" : "提案"}
        </button>
      </div>

      {candidates && (
        <div className="list" style={{ marginTop: "1rem" }}>
          {candidates.length === 0 && (
            <p className="empty">条件に合うお店が見つかりませんでした</p>
          )}
          {candidates.map((c) => (
            <div key={c.id} className="list-item">
              <div className="stack">
                <strong>{c.name}</strong>
                <div className="candidate-meta">
                  <span>📍 {c.area}</span>
                  <span>💰 ~¥{c.budget.toLocaleString()}</span>
                  <span>{c.genres.join(" / ")}</span>
                </div>
              </div>
              {c.matched_all ? (
                <span className="badge badge-ok">全員OK</span>
              ) : (
                <span className="badge badge-warn">妥協案</span>
              )}
            </div>
          ))}
        </div>
      )}
      {error && <p className="error">{error}</p>}
    </div>
  );
}
