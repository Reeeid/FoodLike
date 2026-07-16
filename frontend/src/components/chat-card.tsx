"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { api } from "@/lib/api";
import { setLastRead } from "@/lib/unread";
import type { Message } from "@/types";

// グループチャット。メンバー同士のやり取り+AI検索(回答はストリーミング表示)。
// 新着は3秒間隔のポーリングで取得する。
export function ChatCard({
  groupId,
  myMemberId,
}: {
  groupId: number;
  myMemberId: number;
}) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [text, setText] = useState("");
  // null=非表示 / 文字列=ストリーミング中のAI回答(空文字は考え中)
  const [aiText, setAiText] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const logRef = useRef<HTMLDivElement>(null);
  const lastIdRef = useRef(0);

  const append = useCallback((incoming: Message[]) => {
    if (incoming.length === 0) return;
    setMessages((prev) => {
      const seen = new Set(prev.map((m) => m.id));
      const added = incoming.filter((m) => !seen.has(m.id));
      if (added.length === 0) return prev;
      return [...prev, ...added].sort((a, b) => a.id - b.id);
    });
  }, []);

  // 新着反映のたびに: 取得位置と既読位置を進め、最下部へスクロール
  useEffect(() => {
    const last = messages[messages.length - 1];
    if (last) {
      lastIdRef.current = Math.max(lastIdRef.current, last.id);
      setLastRead(groupId, last.id);
    }
    logRef.current?.scrollTo({ top: logRef.current.scrollHeight });
  }, [messages, aiText, groupId]);

  // ポーリング
  useEffect(() => {
    let cancelled = false;
    const tick = async () => {
      // 非表示タブでは無駄打ちしない(バッテリー/通信/無料枠の節約)
      if (document.hidden) return;
      try {
        const news = await api.listMessages(groupId, lastIdRef.current);
        if (!cancelled) append(news);
      } catch {
        // 一時的な失敗は次のポーリングに任せる
      }
    };
    tick();
    const timer = setInterval(tick, 3000);
    return () => {
      cancelled = true;
      clearInterval(timer);
    };
  }, [groupId, append]);

  const send = async () => {
    const v = text.trim();
    if (!v) return;
    setBusy(true);
    setError("");
    try {
      const m = await api.postMessage(groupId, v);
      append([m]);
      setText("");
    } catch (e) {
      setError(e instanceof Error ? e.message : "送信に失敗しました");
    } finally {
      setBusy(false);
    }
  };

  const aiSearch = async () => {
    const v = text.trim();
    if (!v) return;
    setBusy(true);
    setError("");
    setAiText("");
    try {
      const { question, answer } = await api.aiSearch(groupId, v, (chunk) =>
        setAiText((prev) => (prev ?? "") + chunk),
      );
      append([question, answer]);
      setText("");
    } catch (e) {
      setError(e instanceof Error ? e.message : "AI検索に失敗しました");
    } finally {
      setAiText(null);
      setBusy(false);
    }
  };

  return (
    <div className="card">
      <h2 className="card-title">💬 グループチャット</h2>
      <p className="section-note">
        お店の相談はここで。「AI検索」を押すと、みんなの苦手なものを踏まえたおすすめがチャットに流れます。
      </p>

      <div className="chat-log" ref={logRef}>
        {messages.length === 0 && aiText === null && (
          <p className="empty">まだメッセージがありません</p>
        )}
        {messages.map((m) => (
          <ChatBubble
            key={m.id}
            message={m}
            mine={m.role === "member" && m.member_id === myMemberId}
          />
        ))}
        {aiText !== null && (
          <div className="bubble-row">
            <div className="bubble bubble-ai">
              <span className="bubble-name">🤖 AI</span>
              <span className="bubble-text">
                {aiText === "" ? "考え中…" : aiText}
              </span>
            </div>
          </div>
        )}
      </div>

      <div className="field-row" style={{ marginTop: "0.75rem" }}>
        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && !busy && send()}
          placeholder="メッセージ / AIに聞きたいこと"
          maxLength={200}
        />
        <button className="btn" onClick={send} disabled={busy || !text.trim()}>
          送信
        </button>
        <button
          className="btn btn-ghost"
          onClick={aiSearch}
          disabled={busy || !text.trim()}
        >
          AI検索
        </button>
      </div>
      {error && <p className="error">{error}</p>}
    </div>
  );
}

function ChatBubble({ message, mine }: { message: Message; mine: boolean }) {
  const cls =
    message.role === "ai" ? "bubble-ai" : mine ? "bubble-mine" : "bubble-other";
  return (
    <div className={`bubble-row${mine ? " bubble-row-mine" : ""}`}>
      <div className={`bubble ${cls}`}>
        {!mine && (
          <span className="bubble-name">
            {message.role === "ai" ? "🤖 AI" : `👤 ${message.member_name}`}
          </span>
        )}
        <span className="bubble-text">{message.text}</span>
      </div>
    </div>
  );
}
