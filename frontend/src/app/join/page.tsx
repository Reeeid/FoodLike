"use client";

import { Suspense, useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { api } from "@/lib/api";
import { useSession } from "@/lib/session";

// QRコード(/join?code=xxx)の着地ページ。
// ログイン済みならそのまま参加してグループページへ、未ログインならホームへ誘導する。
export default function JoinPage() {
  return (
    <Suspense fallback={<p className="empty">読み込み中…</p>}>
      <JoinContent />
    </Suspense>
  );
}

function JoinContent() {
  const router = useRouter();
  const code = useSearchParams().get("code") ?? "";
  const session = useSession();
  const [error, setError] = useState("");
  const joining = useRef(false);

  useEffect(() => {
    if (session.status !== "authed" || !code || joining.current) return;
    joining.current = true;
    api
      .joinGroup(code)
      .then((g) => router.replace(`/groups/${g.id}`))
      .catch((e) => {
        joining.current = false;
        setError(e instanceof Error ? e.message : "参加に失敗しました");
      });
  }, [session.status, code, router]);

  if (!code) {
    return (
      <div className="card">
        <h2 className="card-title">招待コードがありません</h2>
        <p className="section-note">招待リンクが正しいか確認してください。</p>
        <Link href="/" className="back-link">
          ← ホームへ
        </Link>
      </div>
    );
  }

  if (session.status === "loading") {
    return <p className="empty">読み込み中…</p>;
  }

  if (session.status === "guest") {
    return (
      <div className="card">
        <h2 className="card-title">グループに招待されています 🎉</h2>
        <p className="section-note">
          参加するには先にログイン(または登録)して、ホームの「招待コードで参加」に下のコードを入力してください。
        </p>
        <div className="invite-code" style={{ textAlign: "center" }}>
          {code}
        </div>
        <Link href="/" className="back-link" style={{ marginTop: "0.75rem" }}>
          ← ホームへ
        </Link>
      </div>
    );
  }

  return (
    <div className="card">
      <h2 className="card-title">グループに参加中…</h2>
      {error ? (
        <>
          <p className="error">{error}</p>
          <Link href="/" className="back-link">
            ← ホームへ
          </Link>
        </>
      ) : (
        <p className="empty">しばらくお待ちください</p>
      )}
    </div>
  );
}
