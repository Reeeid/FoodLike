"use client";

import { useCallback, useEffect, useState } from "react";
import { api } from "@/lib/api";
import type { Preference, PreferenceCategory } from "@/types";

// 自分の苦手なもの編集カード。データはユーザー単位(/api/me/preferences)。
export function PreferenceCard() {
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
