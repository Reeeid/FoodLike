"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import { getLastRead } from "@/lib/unread";
import { useSession } from "@/lib/session";
import type { Group } from "@/types";

type UnreadGroup = { group: Group; count: number };

// ヘッダー右上の通知ベル。全グループ横断で未読を集計し、
// クリックで未読のあるグループへリンク誘導する。
// 全ページに常駐するため、ポーリングは可視時のみ・15秒間隔に抑える。
export function NotificationBell() {
  const session = useSession();
  const [open, setOpen] = useState(false);
  const [items, setItems] = useState<UnreadGroup[]>([]);

  const refresh = useCallback(async () => {
    if (typeof document !== "undefined" && document.hidden) return;
    try {
      const groups = await api.listGroups();
      const withCounts = await Promise.all(
        groups.map(async (g) => ({
          group: g,
          count: (await api.listMessages(g.id, getLastRead(g.id))).length,
        })),
      );
      setItems(withCounts.filter((x) => x.count > 0));
    } catch {
      // 一時的な失敗は次のポーリングに任せる
    }
  }, []);

  useEffect(() => {
    // 未認証ならポーリングを張らない(このときコンポーネントはnullを返す)。
    // setItemsは全てコールバック(タイマー/イベント)経由で呼び、
    // effect本体からの同期setStateを避ける。初回もsetTimeout(0)で1tick遅延。
    if (session.status !== "authed") return;
    const kickoff = setTimeout(refresh, 0);
    const timer = setInterval(refresh, 15000);
    const onVisible = () => {
      if (!document.hidden) refresh();
    };
    document.addEventListener("visibilitychange", onVisible);
    return () => {
      clearTimeout(kickoff);
      clearInterval(timer);
      document.removeEventListener("visibilitychange", onVisible);
    };
  }, [session.status, refresh]);

  if (session.status !== "authed") return null;

  const total = items.reduce((sum, x) => sum + x.count, 0);

  return (
    <div className="bell">
      <button
        className="bell-btn"
        onClick={() => setOpen((v) => !v)}
        aria-label={`通知${total > 0 ? ` (未読${total}件)` : ""}`}
      >
        {/* Font Awesome solid bell(SVGパスのみインライン、ライブラリ非依存) */}
        <svg viewBox="0 0 448 512" width="20" height="20" fill="currentColor" aria-hidden="true">
          <path d="M224 0c-17.7 0-32 14.3-32 32V49.9C119.5 61.4 64 124.2 64 200v33.4c0 45.4-15.5 89.5-43.8 124.9L5.3 377c-5.8 7.2-6.9 17.1-2.9 25.4S14.8 416 24 416H424c9.2 0 17.6-5.3 21.6-13.6s2.9-18.2-2.9-25.4l-14.9-18.6C399.5 322.9 384 278.8 384 233.4V200c0-75.8-55.5-138.6-128-150.1V32c0-17.7-14.3-32-32-32zm45.3 493.3c12-12 18.7-28.3 18.7-45.3H160c0 17 6.7 33.3 18.7 45.3s28.3 18.7 45.3 18.7s33.3-6.7 45.3-18.7z" />
        </svg>
        {total > 0 && (
          <span className="bell-badge">{total > 9 ? "9+" : total}</span>
        )}
      </button>
      {open && (
        <div className="bell-panel">
          {items.length === 0 ? (
            <p className="bell-empty">未読はありません</p>
          ) : (
            items.map(({ group, count }) => (
              <Link
                key={group.id}
                href={`/groups/${group.id}`}
                className="bell-item"
                onClick={() => setOpen(false)}
              >
                <span>{group.name}</span>
                <span className="badge badge-warn">💬 {count}</span>
              </Link>
            ))
          )}
        </div>
      )}
    </div>
  );
}
