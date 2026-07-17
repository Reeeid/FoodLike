"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import Select from "react-select";
import { api } from "@/lib/api";
import { useSession } from "@/lib/session";
import { InviteModal } from "@/components/invite-modal";
import { ChatCard } from "@/components/chat-card";
import type { Candidate, Group, SearchOptions } from "@/types";

type AreaOption = { value: string; label: string };

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
  const [options, setOptions] = useState<SearchOptions | null>(null);
  const [areaOption, setAreaOption] = useState<AreaOption | null>(null);
  const [budget, setBudget] = useState("");
  const [candidates, setCandidates] = useState<Candidate[] | null>(null);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  // エリア/予算のマスタを取得(HOTPEPPER_API_KEY未設定なら空=絞り込み無効)。
  useEffect(() => {
    api
      .getSearchOptions()
      .then(setOptions)
      .catch(() => {
        // マスタが無くても提案自体は動く(全エリア扱い)。黙って無視。
      });
  }, []);

  // 中エリア粒度(例: 新宿・代々木・大久保)で重複を除いて検索候補にする。
  const areaOptions = useMemo<AreaOption[]>(() => {
    const seen = new Map<string, AreaOption>();
    for (const a of options?.areas ?? []) {
      if (a.middle_code && !seen.has(a.middle_code)) {
        seen.set(a.middle_code, {
          value: a.middle_code,
          label: `${a.middle_name}（${a.large_name}）`,
        });
      }
    }
    return [...seen.values()];
  }, [options]);

  const suggest = async () => {
    setBusy(true);
    setError("");
    try {
      setCandidates(
        await api.suggest(groupId, {
          middle_area: areaOption?.value,
          budget: budget || undefined,
        }),
      );
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

      <div className="stack" style={{ gap: "0.5rem" }}>
        <Select<AreaOption>
          instanceId="area-select"
          classNamePrefix="rs"
          options={areaOptions}
          value={areaOption}
          onChange={setAreaOption}
          isClearable
          isDisabled={areaOptions.length === 0}
          placeholder={
            areaOptions.length === 0
              ? "エリア指定なし(全エリア)"
              : "エリアを検索(任意)"
          }
          noOptionsMessage={() => "該当なし"}
        />
        <select value={budget} onChange={(e) => setBudget(e.target.value)}>
          <option value="">予算指定なし</option>
          {(options?.budgets ?? []).map((b) => (
            <option key={b.code} value={b.code}>
              {b.name}
            </option>
          ))}
        </select>
      </div>

      {/* 最重要CTA: フィッツの法則に基づき、幅いっぱい・高さ大で
          カード下部(親指ゾーン)に配置し、迷わず押せるようにする */}
      <button className="btn btn-cta" onClick={suggest} disabled={busy}>
        {busy ? "検索中…" : "🔍 お店を提案してもらう"}
      </button>

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
