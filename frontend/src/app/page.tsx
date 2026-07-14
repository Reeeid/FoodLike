"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import { saveMember } from "@/lib/member-store";
import { getLastRead } from "@/lib/unread";
import { useSession } from "@/lib/session";
import { LoginScreen } from "./login-screen";
import type { Group, Member } from "@/types";

export default function HomePage() {
  const session = useSession();

  if (session.status === "loading") {
    return <p className="empty">読み込み中…</p>;
  }

  if (session.status === "guest") {
    // Firebase設定済みならGoogleログイン、未設定(ローカル開発)は従来の仮登録
    return session.firebaseEnabled ? (
      <LoginScreen
        onGoogleSignIn={session.signInWithGoogle}
        onEmailSignIn={session.signInWithEmail}
        onEmailSignUp={session.signUpWithEmail}
        onGuestSignIn={session.signInAsGuest}
      />
    ) : (
      <RegisterCard onRegistered={saveMember} />
    );
  }

  const member = session.member!;
  return (
    <>
      <div className="card">
        <div className="list-item">
          <span>
            👤 {member.name} <span className="muted">(ID: {member.id})</span>
          </span>
          <div className="field-row" style={{ width: "auto" }}>
            <Link href="/profile" className="btn btn-ghost btn-sm">
              プロフィール
            </Link>
            <button className="btn btn-ghost btn-sm" onClick={session.signOut}>
              {session.firebaseEnabled ? "ログアウト" : "切り替え"}
            </button>
          </div>
        </div>
      </div>
      <GroupSection />
    </>
  );
}

function RegisterCard({
  onRegistered,
}: {
  onRegistered: (m: Member) => void;
}) {
  const [name, setName] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  const register = async () => {
    if (!name.trim()) return;
    setBusy(true);
    setError("");
    try {
      onRegistered(await api.registerMember(name.trim()));
    } catch (e) {
      setError(e instanceof Error ? e.message : "登録に失敗しました");
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="card">
      <h2 className="card-title">はじめまして 👋</h2>
      <p className="section-note">
        ニックネームを登録して始めましょう。好き嫌いはあなた以外には見えません。
      </p>
      <div className="field-row">
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && register()}
          placeholder="ニックネーム"
          maxLength={64}
        />
        <button className="btn" onClick={register} disabled={busy || !name.trim()}>
          登録
        </button>
      </div>
      {error && <p className="error">{error}</p>}
    </div>
  );
}

function GroupSection() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [groupName, setGroupName] = useState("");
  const [inviteCode, setInviteCode] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  // アプリ内通知: グループごとの未読メッセージ数(既読位置より新しい件数)
  const [unread, setUnread] = useState<Record<number, number>>({});

  useEffect(() => {
    let cancelled = false;
    (async () => {
      const entries = await Promise.all(
        groups.map(async (g) => {
          try {
            const news = await api.listMessages(g.id, getLastRead(g.id));
            return [g.id, news.length] as const;
          } catch {
            return [g.id, 0] as const;
          }
        }),
      );
      if (!cancelled) setUnread(Object.fromEntries(entries));
    })();
    return () => {
      cancelled = true;
    };
  }, [groups]);

  const refresh = useCallback(async () => {
    try {
      setGroups(await api.listGroups());
    } catch (e) {
      setError(e instanceof Error ? e.message : "グループの取得に失敗しました");
    }
  }, []);

  useEffect(() => {
    api
      .listGroups()
      .then(setGroups)
      .catch((e: unknown) =>
        setError(
          e instanceof Error ? e.message : "グループの取得に失敗しました",
        ),
      );
  }, []);

  const createGroup = async () => {
    if (!groupName.trim()) return;
    setBusy(true);
    setError("");
    try {
      await api.createGroup(groupName.trim());
      setGroupName("");
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : "作成に失敗しました");
    } finally {
      setBusy(false);
    }
  };

  const joinGroup = async () => {
    if (!inviteCode.trim()) return;
    setBusy(true);
    setError("");
    try {
      await api.joinGroup(inviteCode.trim());
      setInviteCode("");
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : "参加に失敗しました");
    } finally {
      setBusy(false);
    }
  };

  return (
    <>
      <div className="card">
        <h2 className="card-title">グループ</h2>
        {groups.length === 0 ? (
          <p className="empty">まだグループがありません</p>
        ) : (
          <div className="list">
            {groups.map((g) => (
              <Link key={g.id} href={`/groups/${g.id}`} className="list-item-link">
                <div className="stack">
                  <strong>
                    {g.name}
                    {(unread[g.id] ?? 0) > 0 && (
                      <span
                        className="badge badge-warn"
                        style={{ marginLeft: "0.5rem" }}
                      >
                        💬 未読{unread[g.id]}件
                      </span>
                    )}
                  </strong>
                  <span className="muted">{g.members.length}人のメンバー</span>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>

      <div className="card">
        <h2 className="card-title">グループを作る</h2>
        <div className="field-row">
          <input
            value={groupName}
            onChange={(e) => setGroupName(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && createGroup()}
            placeholder="グループ名(家族、サークル…)"
            maxLength={64}
          />
          <button className="btn" onClick={createGroup} disabled={busy || !groupName.trim()}>
            作成
          </button>
        </div>
      </div>

      <div className="card">
        <h2 className="card-title">招待コードで参加</h2>
        <div className="field-row">
          <input
            value={inviteCode}
            onChange={(e) => setInviteCode(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && joinGroup()}
            placeholder="招待コード"
            maxLength={32}
          />
          <button className="btn" onClick={joinGroup} disabled={busy || !inviteCode.trim()}>
            参加
          </button>
        </div>
      </div>

      {error && <p className="error">{error}</p>}
    </>
  );
}
