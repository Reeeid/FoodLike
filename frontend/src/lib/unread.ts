// グループごとの既読位置(最後に読んだメッセージID)をlocalStorageに保持する。
// アプリ内通知(未読バッジ)は「既読位置より新しいメッセージ数」で出す。

const key = (groupId: number) => `foodlike:lastRead:${groupId}`;

export function getLastRead(groupId: number): number {
  if (typeof window === "undefined") return 0;
  return Number(localStorage.getItem(key(groupId))) || 0;
}

export function setLastRead(groupId: number, messageId: number) {
  if (typeof window === "undefined") return;
  if (messageId > getLastRead(groupId)) {
    localStorage.setItem(key(groupId), String(messageId));
  }
}
