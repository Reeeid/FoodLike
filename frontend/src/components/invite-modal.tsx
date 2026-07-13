"use client";

import { useEffect, useMemo, useState } from "react";
import QRCode from "react-qr-code";
import type { Group } from "@/types";

// 招待コードとQRコードを表示するモーダル。
// QRには /join?code=xxx のURLを埋め込み、スマホで読み取るとそのまま参加できる。
export function InviteModal({
  group,
  onClose,
}: {
  group: Group;
  onClose: () => void;
}) {
  const [copied, setCopied] = useState(false);

  const joinUrl = useMemo(
    () =>
      `${window.location.origin}/join?code=${encodeURIComponent(group.invite_code)}`,
    [group.invite_code],
  );

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => e.key === "Escape" && onClose();
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [onClose]);

  const copy = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch {
      // クリップボード非対応環境ではuser-select: allのコード表示から手動コピーしてもらう
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div
        className="modal"
        role="dialog"
        aria-modal="true"
        aria-label={`${group.name} への招待`}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="modal-header">
          <h2 className="card-title" style={{ marginBottom: 0 }}>
            友だちを招待
          </h2>
          <button className="modal-close" onClick={onClose} aria-label="閉じる">
            ×
          </button>
        </div>

        <p className="section-note">
          QRコードを読み取るか、招待コードを共有すると「{group.name}
          」に参加できます。
        </p>

        <div className="qr-box">
          <QRCode value={joinUrl} size={192} />
        </div>

        <div className="field-row" style={{ marginTop: "1rem" }}>
          <div className="invite-code" style={{ flex: 1, textAlign: "center" }}>
            {group.invite_code}
          </div>
          <button className="btn" onClick={() => copy(group.invite_code)}>
            {copied ? "コピーしました" : "コピー"}
          </button>
        </div>

        <button
          className="btn btn-ghost"
          style={{ width: "100%", marginTop: "0.75rem" }}
          onClick={() => copy(joinUrl)}
        >
          招待リンクをコピー
        </button>
      </div>
    </div>
  );
}
