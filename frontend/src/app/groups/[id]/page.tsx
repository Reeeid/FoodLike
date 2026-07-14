"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { api } from "@/lib/api";
import { useSession } from "@/lib/session";
import { InviteModal } from "@/components/invite-modal";
import { ChatCard } from "@/components/chat-card";
import type { Candidate, Group } from "@/types";

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

      <ChatCard groupId={groupId} myMemberId={session.member?.id ?? 0} />
      <SuggestionCard groupId={groupId} />
    </>
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
        メンバー全員の苦手なものにこっそり配慮したお店を探します。あなたの苦手なものは
        <Link href="/profile" className="inline-link">
          プロフィール
        </Link>
        で登録できます。
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
