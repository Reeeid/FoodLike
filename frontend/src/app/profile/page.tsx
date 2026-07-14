"use client";

import { useEffect } from "react";
import Link from "next/link";
import { PreferenceCard } from "@/components/preference-card";
import { useSession } from "@/lib/session";

export default function ProfilePage() {
  const session = useSession();

  useEffect(() => {
    if (session.status === "guest") {
      window.location.href = "/";
    }
  }, [session.status]);

  if (session.status !== "authed") {
    return <p className="empty">読み込み中…</p>;
  }

  const member = session.member!;
  return (
    <>
      <Link href="/" className="back-link">
        ← ホームへ
      </Link>

      <div className="card">
        <h2 className="card-title">プロフィール</h2>
        <div className="list-item">
          <span>
            👤 {member.name} <span className="muted">(ID: {member.id})</span>
          </span>
        </div>
      </div>

      <PreferenceCard />
    </>
  );
}
